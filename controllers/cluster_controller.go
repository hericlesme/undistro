/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	undistrov1 "github.com/getupcloud/undistro/api/v1alpha1"
	uclient "github.com/getupcloud/undistro/client"
	"github.com/getupcloud/undistro/client/config"
	"github.com/getupcloud/undistro/internal/patch"
	"github.com/getupcloud/undistro/internal/util"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	clusterApi "sigs.k8s.io/cluster-api/api/v1alpha3"
	kubeadmApi "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha3"
	utilresource "sigs.k8s.io/cluster-api/util/resource"
	"sigs.k8s.io/cluster-api/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	jobOwnerKey = ".metadata.controller"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	RestConfig *rest.Config
}

// +kubebuilder:rbac:urls=/metrics,verbs=get;
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes/custom-host,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=acme.cert-manager.io,resources=*,verbs=deletecollection;get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=*,verbs=deletecollection;get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=*,verbs=deletecollection;get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=auditregistration.k8s.io,resources=auditsinks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiregistration.k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=getupcloud.com,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=getupcloud.com,resources=clusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=getupcloud.com,resources=providers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=delete;get;list;watch;create;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io;bootstrap.cluster.x-k8s.io;controlplane.cluster.x-k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (res ctrl.Result, err error) {
	ctx := context.Background()
	log := r.Log.WithValues("cluster", req.NamespacedName)
	var cluster undistrov1.Cluster
	if err := r.Get(ctx, req.NamespacedName, &cluster); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "couldn't get object", "name", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(&cluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer func() {
		err = patch.ControllerObject(ctx, patchHelper, &cluster, err)
	}()
	if !controllerutil.ContainsFinalizer(&cluster, undistrov1.ClusterFinalizer) {
		controllerutil.AddFinalizer(&cluster, undistrov1.ClusterFinalizer)
		return ctrl.Result{}, nil
	}
	if !cluster.DeletionTimestamp.IsZero() {
		log.Info("removing cluster", "name", req.NamespacedName)
		if err := r.delete(ctx, &cluster); err != nil {
			return ctrl.Result{}, err
		}
		if cluster.Status.Phase == undistrov1.DeletingPhase {
			return ctrl.Result{Requeue: true}, nil
		}
		controllerutil.RemoveFinalizer(&cluster, undistrov1.ClusterFinalizer)
		return ctrl.Result{}, nil
	}
	undistroClient, err := uclient.New("")
	if err != nil {
		return ctrl.Result{}, err
	}
	err = util.SetVariablesFromEnvVar(ctx, util.VariablesInput{
		VariablesClient: undistroClient.GetVariables(),
		ClientSet:       r.Client,
		NamespacedName:  req.NamespacedName,
		EnvVars:         cluster.Spec.InfrastructureProvider.Env,
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	err = config.TrySetCustomTemplates(&cluster, undistroClient.GetVariables())
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster.Status.Phase == undistrov1.NewPhase {
		log.Info("ensure mangement cluster is initialized and updated", "name", req.NamespacedName)
		if err = r.init(ctx, &cluster, undistroClient); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "couldn't initialize or update the mangement cluster", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if cluster.Status.Phase == undistrov1.InitializedPhase {
		log.Info("generanting cluster-api configuration", "name", req.NamespacedName)
		if err = r.config(ctx, &cluster, undistroClient); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "couldn't install cluster config", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if cluster.Status.Phase == undistrov1.ProvisioningPhase {
		res, err := r.provisioning(ctx, &cluster, undistroClient)
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "couldn't provision the cluster", "name", req.NamespacedName)
			return res, err
		}
		return res, nil
	}
	if r.hasDiff(ctx, &cluster) {
		return r.upgrade(ctx, &cluster, undistroClient)
	}
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) hasDiff(ctx context.Context, cl *undistrov1.Cluster) bool {
	diffCP := cmp.Diff(cl.Spec.ControlPlaneNode, cl.Status.ControlPlaneNode)
	diffW := cmp.Diff(cl.Spec.WorkerNode, cl.Status.WorkerNode)
	switch {
	case cl.Spec.KubernetesVersion != cl.Status.KubernetesVersion:
		return true
	case diffCP != "":
		return true
	case diffW != "":
		return true
	default:
		return false
	}
}

func (r *ClusterReconciler) delete(ctx context.Context, cl *undistrov1.Cluster) error {
	log := r.Log
	cl.Status.Ready = false
	if cl.Status.ClusterAPIRef == nil {
		cl.Status.Phase = undistrov1.DeletedPhase
		return nil
	}
	capiNM := types.NamespacedName{
		Name:      cl.Status.ClusterAPIRef.Name,
		Namespace: cl.Status.ClusterAPIRef.Namespace,
	}
	var capi clusterApi.Cluster
	err := r.Get(ctx, capiNM, &capi)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "couldn't get capi", "name", capiNM)
			return err
		}
		cl.Status.Phase = undistrov1.DeletedPhase
		return nil
	}
	if capi.Status.Phase != string(clusterApi.ClusterPhaseDeleting) {
		if err := r.Delete(ctx, &capi); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "couldn't delete capi", "name", capiNM)
				return err
			}
			cl.Status.Phase = undistrov1.DeletedPhase
			return nil
		}
	}
	cl.Status.Phase = undistrov1.DeletingPhase
	return nil
}

func (r *ClusterReconciler) init(ctx context.Context, cl *undistrov1.Cluster, c uclient.Client) error {
	opts := uclient.InitOptions{
		Kubeconfig: uclient.Kubeconfig{
			RestConfig: r.RestConfig,
		},
		InfrastructureProviders: []string{cl.Spec.InfrastructureProvider.NameVersion()},
		TargetNamespace:         "undistro-system",
		LogUsageInstructions:    false,
	}
	if cl.Spec.BootstrapProvider != nil {
		opts.BootstrapProviders = []string{cl.Spec.BootstrapProvider.NameVersion()}
	}
	if cl.Spec.ControlPlaneProvider != nil {
		opts.ControlPlaneProviders = []string{cl.Spec.ControlPlaneProvider.NameVersion()}
	}
	components, err := c.Init(opts)
	if err != nil {
		return err
	}
	if len(components) == 0 {
		var comp uclient.Components
		comp, err = c.GetProviderComponents(cl.Spec.InfrastructureProvider.Name, undistrov1.InfrastructureProviderType, uclient.ComponentsOptions{
			TargetNamespace: "undistro-system",
		})
		if err != nil {
			return err
		}
		components = append(components, comp)
	}
	cl.Status.InstalledComponents = make([]undistrov1.InstalledComponent, len(components))
	for i, component := range components {
		ic := undistrov1.InstalledComponent{
			Name:    component.Name(),
			Version: component.Version(),
			URL:     component.URL(),
			Type:    component.Type(),
		}
		cl.Status.InstalledComponents[i] = ic
	}
	cl.Status.Phase = undistrov1.InitializedPhase
	cl.Status.KubernetesVersion = cl.Spec.KubernetesVersion
	cl.Status.WorkerNode = cl.Spec.WorkerNode
	cl.Status.ControlPlaneNode = cl.Spec.ControlPlaneNode
	cl.Status.InfrastructureName = cl.Spec.InfrastructureProvider.Name
	return nil
}

func (r *ClusterReconciler) config(ctx context.Context, cl *undistrov1.Cluster, c uclient.Client) error {
	log := r.Log
	for _, ic := range cl.Status.InstalledComponents {
		comp, err := uclient.GetProvider(c, ic.Name, ic.Type)
		if err != nil {
			return err
		}
		preConfigFunc := comp.GetPreConfigFunc()
		if preConfigFunc != nil {
			log.Info("executing pre config func", "component", comp.Name())
			err = preConfigFunc(ctx, cl, c.GetVariables(), r.Client)
			if err != nil {
				return err
			}
		}
	}
	opts := uclient.GetClusterTemplateOptions{
		Kubeconfig: uclient.Kubeconfig{
			RestConfig: r.RestConfig,
		},
		ClusterName:              cl.Name,
		TargetNamespace:          cl.Namespace,
		ListVariablesOnly:        false,
		KubernetesVersion:        cl.Spec.KubernetesVersion,
		ControlPlaneMachineCount: cl.Spec.ControlPlaneNode.Replicas,
		WorkerMachineCount:       cl.Spec.WorkerNode.Replicas,
	}
	if cl.Spec.Template != nil {
		opts.URLSource = &uclient.URLSourceOptions{
			URL: *cl.Spec.Template,
		}
	}
	if cl.Spec.BootstrapProvider != nil && cl.Spec.Template == nil {
		opts.ProviderRepositorySource = &uclient.ProviderRepositorySourceOptions{
			InfrastructureProvider: cl.Spec.InfrastructureProvider.Name,
			Flavor:                 cl.Spec.BootstrapProvider.Name,
		}
	}
	tpl, err := c.GetClusterTemplate(opts)
	if err != nil {
		return err
	}
	objs := utilresource.SortForCreate(tpl.Objs())
	for _, o := range objs {
		isCluster := false
		if o.GetKind() == "Cluster" && o.GroupVersionKind().GroupVersion().String() == clusterApi.GroupVersion.String() {
			isCluster = true
			err = ctrl.SetControllerReference(cl, &o, r.Scheme)
			if err != nil {
				return errors.Errorf("couldn't set reference: %v", err)
			}
		}
		err = r.Patch(ctx, &o, client.Apply, client.FieldOwner("undistro"))
		if err != nil {
			return err
		}
		if isCluster {
			cl.Status.ClusterAPIRef = &corev1.ObjectReference{
				Kind:            o.GetKind(),
				Namespace:       o.GetNamespace(),
				Name:            o.GetName(),
				UID:             o.GetUID(),
				ResourceVersion: o.GetResourceVersion(),
				APIVersion:      o.GetAPIVersion(),
			}
		}
	}
	cl.Status.Phase = undistrov1.ProvisioningPhase
	return nil
}

func (r *ClusterReconciler) upgrade(ctx context.Context, cl *undistrov1.Cluster, uc uclient.Client) (ctrl.Result, error) {
	if cl.Status.ClusterAPIRef == nil {
		return ctrl.Result{}, errors.New("cluster API reference is nil")
	}
	capi := clusterApi.Cluster{}
	nm := types.NamespacedName{
		Name:      cl.Status.ClusterAPIRef.Name,
		Namespace: cl.Status.ClusterAPIRef.Namespace,
	}
	if err := r.Get(ctx, nm, &capi); err != nil {
		return ctrl.Result{}, err
	}
	actual := cl.Status.DeepCopy()
	cl.Status.KubernetesVersion = cl.Spec.KubernetesVersion
	cl.Status.WorkerNode = cl.Spec.WorkerNode
	cl.Status.ControlPlaneNode = cl.Spec.ControlPlaneNode
	cl.Status.InfrastructureName = cl.Spec.InfrastructureProvider.Name
	if err := r.Status().Update(ctx, cl); err != nil {
		return ctrl.Result{}, err
	}
	if cl.Namespace == "" {
		cl.Namespace = "default"
	}
	switch {
	case actual.KubernetesVersion != cl.Spec.KubernetesVersion,
		*actual.ControlPlaneNode.Replicas != *cl.Spec.ControlPlaneNode.Replicas,
		*actual.WorkerNode.Replicas != *cl.Spec.WorkerNode.Replicas:
		return ctrl.Result{}, r.upgradeRefs(ctx, cl, &capi, uc)
	case actual.ControlPlaneNode.MachineType != cl.Spec.ControlPlaneNode.MachineType,
		actual.WorkerNode.MachineType != cl.Spec.WorkerNode.MachineType:
		return ctrl.Result{}, r.upgradeInstance(ctx, cl, &capi, actual)
	}
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) upgradeInstance(ctx context.Context, cl *undistrov1.Cluster, capi *clusterApi.Cluster, actual *undistrov1.ClusterStatus) error {
	if actual.ControlPlaneNode.MachineType != cl.Spec.ControlPlaneNode.MachineType {
		nm := types.NamespacedName{
			Name:      capi.Spec.ControlPlaneRef.Name,
			Namespace: capi.Spec.ControlPlaneRef.Namespace,
		}
		kubeadmCP := kubeadmApi.KubeadmControlPlane{}
		err := r.Get(ctx, nm, &kubeadmCP)
		if err != nil {
			return err
		}
		nm.Name = kubeadmCP.Spec.InfrastructureTemplate.Name
		nm.Namespace = kubeadmCP.Spec.InfrastructureTemplate.Namespace
		o := unstructured.Unstructured{}
		o.SetGroupVersionKind(kubeadmCP.Spec.InfrastructureTemplate.GroupVersionKind())
		err = r.Get(ctx, nm, &o)
		if err != nil {
			return err
		}
		newObj := o.DeepCopy()
		newObj.SetName(fmt.Sprintf("%s-control-plane-%s", cl.Name, uuid.New().String()))
		newObj.SetNamespace(cl.Namespace)
		err = unstructured.SetNestedField(newObj.Object, cl.Spec.ControlPlaneNode.MachineType, "spec", "template", "spec", "instanceType")
		if err != nil {
			return err
		}
		newObj.SetResourceVersion("")
		err = r.Create(ctx, newObj)
		if err != nil {
			return err
		}
		kubeadmCP.Spec.InfrastructureTemplate.Name = newObj.GetName()
		kubeadmCP.Spec.InfrastructureTemplate.Namespace = newObj.GetNamespace()
		err = r.Update(ctx, &kubeadmCP)
		if err != nil {
			return err
		}
	}
	if actual.WorkerNode.MachineType != cl.Spec.WorkerNode.MachineType {
		md := clusterApi.MachineDeployment{}
		nm := types.NamespacedName{
			Name:      fmt.Sprintf("%s-md-0", cl.Name),
			Namespace: cl.Namespace,
		}
		err := r.Get(ctx, nm, &md)
		if err != nil {
			return err
		}
		if md.Spec.Template.Spec.InfrastructureRef.Namespace == "" {
			md.Spec.Template.Spec.InfrastructureRef.Namespace = "default"
		}
		nm.Name = md.Spec.Template.Spec.InfrastructureRef.Name
		nm.Namespace = md.Spec.Template.Spec.InfrastructureRef.Namespace
		o := unstructured.Unstructured{}
		o.SetGroupVersionKind(md.Spec.Template.Spec.InfrastructureRef.GroupVersionKind())
		err = r.Get(ctx, nm, &o)
		if err != nil {
			return err
		}
		newObj := o.DeepCopy()
		newObj.SetName(fmt.Sprintf("%s-md-0-%s", cl.Name, uuid.New().String()))
		newObj.SetNamespace(cl.Namespace)
		err = unstructured.SetNestedField(newObj.Object, cl.Spec.WorkerNode.MachineType, "spec", "template", "spec", "instanceType")
		if err != nil {
			return err
		}
		newObj.SetResourceVersion("")
		err = r.Create(ctx, newObj)
		if err != nil {
			return err
		}
		md.Spec.Template.Spec.InfrastructureRef.Name = newObj.GetName()
		md.Spec.Template.Spec.InfrastructureRef.Namespace = newObj.GetNamespace()
		err = r.Update(ctx, &md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ClusterReconciler) upgradeRefs(ctx context.Context, cl *undistrov1.Cluster, capi *clusterApi.Cluster, uc uclient.Client) error {
	log := r.Log
	p, err := uclient.GetProvider(uc, cl.Spec.InfrastructureProvider.Name, undistrov1.InfrastructureProviderType)
	if err != nil {
		return err
	}
	upgradeFunc := p.GetUpgradeFunc()
	if upgradeFunc != nil {
		log.Info("executing upgrade func", "component", p.Name())
		err = upgradeFunc(ctx, cl, capi, r.Client)
		if err != nil {
			return err
		}
	} else {
		o := unstructured.Unstructured{}
		o.SetGroupVersionKind(capi.Spec.ControlPlaneRef.GroupVersionKind())
		nm := types.NamespacedName{
			Name:      capi.Spec.ControlPlaneRef.Name,
			Namespace: capi.Spec.ControlPlaneRef.Namespace,
		}
		err := r.Get(ctx, nm, &o)
		if err != nil {
			return err
		}
		err = unstructured.SetNestedField(o.Object, cl.Spec.KubernetesVersion, "spec", "version")
		if err != nil {
			return err
		}
		err = unstructured.SetNestedField(o.Object, cl.Spec.ControlPlaneNode.Replicas, "spec", "replicas")
		if err != nil {
			return err
		}
		o.SetResourceVersion(capi.Spec.ControlPlaneRef.ResourceVersion)
		err = r.Update(ctx, &o)
		if err != nil {
			return err
		}
	}
	// upgrade worker nodes
	md := clusterApi.MachineDeployment{}
	nm := types.NamespacedName{
		Name:      fmt.Sprintf("%s-md-0", cl.Name),
		Namespace: cl.Namespace,
	}
	err = r.Get(ctx, nm, &md)
	if err != nil {
		return err
	}
	md.Spec.Template.Spec.Version = &cl.Spec.KubernetesVersion
	workerReplicas := int32(*cl.Spec.WorkerNode.Replicas)
	md.Spec.Replicas = &workerReplicas
	return r.Update(ctx, &md)
}

func (r *ClusterReconciler) provisioning(ctx context.Context, cl *undistrov1.Cluster, c uclient.Client) (ctrl.Result, error) {
	log := r.Log
	log.Info("provisioning cluster", "name", cl.Name, "namespace", cl.Namespace)
	var clusterList clusterApi.ClusterList
	if err := r.List(ctx, &clusterList, client.InNamespace(cl.Namespace), client.MatchingFields{jobOwnerKey: cl.Name}); err != nil {
		return ctrl.Result{}, err
	}
	if len(clusterList.Items) != 1 {
		return ctrl.Result{}, errors.Errorf("has more than one Cluster API Cluster owned by %s/%s", cl.Namespace, cl.Name)
	}
	capi := &clusterList.Items[0]
	if undistrov1.ClusterPhase(capi.Status.Phase) == undistrov1.FailedPhase {
		cl.Status.Phase = undistrov1.FailedPhase
		return ctrl.Result{}, nil
	}
	if capi.Status.ControlPlaneInitialized && !capi.Status.ControlPlaneReady && cl.Spec.CniName != undistrov1.ProviderCNI {
		if err := r.installCNI(ctx, cl, c); err != nil {
			return ctrl.Result{}, err
		}
	}
	if capi.Status.ControlPlaneReady && capi.Status.InfrastructureReady && undistrov1.ClusterPhase(capi.Status.Phase) == undistrov1.ProvisionedPhase {
		cl.Status.Phase = undistrov1.ProvisionedPhase
		cl.Status.Ready = true
	}
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) installCNI(ctx context.Context, cl *undistrov1.Cluster, c uclient.Client) error {
	log := r.Log
	cniAddr := cl.GetCNITemplateURL()
	if cniAddr == "" {
		return errors.Errorf("CNI %s is not supported", cl.Spec.CniName)
	}
	log.Info("getting CNI", "name", cl.Spec.CniName, "URL", cniAddr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cniAddr, nil)
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	objs, err := yaml.ToUnstructured(byt)
	if err != nil {
		return err
	}
	wc, err := c.GetWorkloadCluster(uclient.Kubeconfig{
		RestConfig: r.RestConfig,
	})
	if err != nil {
		return err
	}
	workloadCfg, err := wc.GetRestConfig(cl.Name, cl.Namespace)
	if err != nil {
		return err
	}
	workloadClient, err := client.New(workloadCfg, client.Options{Scheme: r.Scheme})
	if err != nil {
		return err
	}
	objs = utilresource.SortForCreate(objs)
	for _, o := range objs {
		workloadClient.Patch(ctx, &o, client.Apply, client.FieldOwner("undistro"))
	}
	return nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager, opts controller.Options) error {
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &clusterApi.Cluster{}, jobOwnerKey, func(rawObj runtime.Object) []string {
		cluster := rawObj.(*clusterApi.Cluster)
		owner := metav1.GetControllerOf(cluster)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != undistrov1.GroupVersion.String() || owner.Kind != "Cluster" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(opts).
		For(&undistrov1.Cluster{}).
		Owns(&clusterApi.Cluster{}).
		Watches(
			&source.Kind{Type: &clusterApi.Cluster{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.capiToUndistro),
			},
		).
		Complete(r)
}

func (r *ClusterReconciler) capiToUndistro(o handler.MapObject) []ctrl.Request {
	ctx := context.TODO()
	c, ok := o.Object.(*clusterApi.Cluster)
	if !ok {
		r.Log.Error(nil, fmt.Sprintf("expected a Cluster but got a %T", o.Object))
		return nil
	}
	if c.Status.Phase == "" {
		return nil
	}
	nm := types.NamespacedName{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
	uc := undistrov1.Cluster{}
	if err := r.Get(ctx, nm, &uc); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "couldn't get undistro cluster")
		}
		return nil
	}
	if uc.Status.Phase == undistrov1.ClusterPhase(c.Status.Phase) {
		return nil
	}
	uc.Status.Phase = undistrov1.ClusterPhase(c.Status.Phase)
	if err := r.Status().Update(ctx, &uc); client.IgnoreNotFound(err) != nil {
		r.Log.Error(err, "couldn't update status", "name", uc.Name, "namespace", uc.Namespace)
		return nil
	}
	return []ctrl.Request{
		{
			NamespacedName: nm,
		},
	}
}
