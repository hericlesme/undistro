/*
Copyright 2020-2021 The UnDistro authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/getupio-undistro/controllerlib"
	"github.com/getupio-undistro/meta"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/cloud"
	"github.com/getupio-undistro/undistro/pkg/fs"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/retry"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/template"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	capi "sigs.k8s.io/cluster-api/api/v1alpha4"
	capicp "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now().UTC()

	// Retrieve UnDistro cluster instance
	uCluster := &appv1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, uCluster); err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Setup log with default values
	keysAndValues := []interface{}{
		"cluster", req.NamespacedName,
		"infra", uCluster.Spec.InfrastructureProvider.Name,
		"flavor", uCluster.Spec.InfrastructureProvider.Flavor,
		"uid", uCluster.UID,
	}
	log := setDefaultLog(ctx)
	log.WithValues(keysAndValues...)

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(uCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer controllerlib.PatchInstance(ctx, controllerlib.InstanceOpts{
		Controller: "ClusterController",
		Request:    req.String(),
		Object:     uCluster,
		Error:      err,
		Helper:     patchHelper,
	})

	log.Info("Checking object generation")
	if uCluster.Generation < uCluster.Status.ObservedGeneration {
		log.Info("Skipping this old version of reconciled object")
		return ctrl.Result{}, nil
	}

	log.Info("Checking if has finalizer")
	if !controllerutil.ContainsFinalizer(uCluster, meta.Finalizer) {
		log.Info("Adding Finalizer")
		controllerutil.AddFinalizer(uCluster, meta.Finalizer)
		return ctrl.Result{}, nil
	}

	log.Info("Checking if object is paused")
	if uCluster.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		uCluster = appv1alpha1.ClusterPaused(uCluster)
		return ctrl.Result{}, nil
	}

	log.Info("Checking if under deletion")
	if !uCluster.DeletionTimestamp.IsZero() {
		uCluster = appv1alpha1.ClusterDeleting(uCluster)
		return r.reconcileDelete(ctx, uCluster)
	}

	result, err := r.reconcile(ctx, uCluster)

	durationMsg := fmt.Sprintf("Reconcilation finished in %s", time.Since(start).String())
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	log.Info(durationMsg)
	return result, err
}

func (r *ClusterReconciler) reconcile(ctx context.Context, cl *appv1alpha1.Cluster) (ctrl.Result, error) {
	log := setDefaultLog(ctx)

	// Retrieve Cluster API Cluster object
	capiCluster := capi.Cluster{}
	err := r.Get(ctx, client.ObjectKeyFromObject(cl), &capiCluster)
	// We ignore because in the first iteration of the loop the CAPI Cluster and related objects will be created bellow
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	// Update the actual worker replicas state
	cl.Status.TotalWorkerPools = int32(len(cl.Spec.Workers))
	cl.Status.TotalWorkerReplicas = 0
	for _, w := range cl.Spec.Workers {
		cl.Status.TotalWorkerReplicas += *w.Replicas
	}
	log.Info("Cluster capabilities", "totalWorkerPools", cl.Status.TotalWorkerPools, "totalWorkerReplicas", cl.Status.TotalWorkerReplicas)

	// In managed flavors, we need to install calico too for network policy support. So we don't verify
	err = r.reconcileCNI(ctx, cl)
	if err != nil {
		meta.SetResourceCondition(cl, meta.CNIInstalledCondition, metav1.ConditionFalse, meta.CNIInstalledFailedReason, err.Error())
		return ctrl.Result{}, err
	}

	// Install all the AutoScaler related
	if cl.Annotations != nil {
		_, autoscalerEnabled := cl.Annotations[meta.EnableClusterAutoscaler]
		if autoscalerEnabled {
			err := r.reconcileClusterAutoscaler(ctx, cl)
			if err != nil {
				meta.SetResourceCondition(cl, meta.ClusterAutoscalerInstalledCondition, metav1.ConditionFalse, meta.ClusterAutoscalerInstalledFailedReason, err.Error())
				return ctrl.Result{}, err
			}
		}
	}

	// Install and config cloud provider specific things
	err = cloud.ReconcileProviderIntegration(ctx, r.Client, log, cl)
	if err != nil {
		meta.SetResourceCondition(cl, meta.CloudProviderInstalledCondition, metav1.ConditionFalse, meta.CloudProvideInstalledFailedReason, err.Error())
		return ctrl.Result{}, err
	}
	log.Info("Cloud provider integration reconciled")

	// Check if a bastion was specified
	if cl.Spec.Bastion != nil {
		log.Info("Bastion specification exists", "enabled", *cl.Spec.Bastion.Enabled)
		if *cl.Spec.Bastion.Enabled && cl.Status.BastionPublicIP == "" {
			cl.Status.BastionPublicIP, err = r.tryRetrieveBastionIP(ctx, capiCluster)
			// In the first iteration, bastion is nil because the CAPI Cluster has not been created yet
			if err != nil {
				cl = appv1alpha1.ClusterNotReady(cl, meta.WaitProvisionReason, err.Error())
				return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
			}
		}
	}

	log.Info("Reconciling cluster config for each provider")
	err = cloud.ReconcileClusterConf(ctx, r.Client, cl, &capiCluster)
	if err != nil {
		cl = appv1alpha1.ClusterNotReady(cl, meta.ReconcileLaunchTemplateFailed, err.Error())
		return ctrl.Result{}, err
	}

	log.Info("Reconciling cluster network")
	err = cloud.ReconcileNetwork(ctx, r.Client, cl, &capiCluster)
	if err != nil {
		cl = appv1alpha1.ClusterNotReady(cl, meta.ReconcileNetworkFailed, err.Error())
		return ctrl.Result{}, err
	}

	return r.reconcileCAPIClusterTemplate(ctx, cl, capiCluster)
}

func (r *ClusterReconciler) reconcileCAPIClusterTemplate(ctx context.Context, cl *appv1alpha1.Cluster, capiCl capi.Cluster) (ctrl.Result, error) {
	log := setDefaultLog(ctx)
	log.Info("Checking if has diff between templates", "spec", cl.Spec, "status", cl.Status)
	if r.hasDiff(ctx, cl) {
		log.Info("Has diff!")
		vars, err := r.fillTemplateVariables(ctx, r.Client, cl)
		if err != nil {
			cl = appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error())
			return ctrl.Result{}, err
		}
		log.Info("Template variables", "vars", vars)

		objs, err := template.GetObjs(fs.FS, "clustertemplates", cl.GetTemplate(), vars)
		if err != nil {
			cl = appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error())
			return ctrl.Result{}, err
		}
		// OpenStackCluster

		switch cl.Spec.InfrastructureProvider.Name {
		case appv1alpha1.OpenStack.String():
			// apply objects in specific way
		case appv1alpha1.Amazon.String():
			err := r.applyObjects(ctx, cl, objs)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	cl.Status.KubernetesVersion = cl.Spec.KubernetesVersion
	cl.Status.ControlPlane = *cl.Spec.ControlPlane
	cl.Status.Workers = cl.Spec.Workers
	cl.Status.BastionConfig = cl.Spec.Bastion
	log.Info("Cluster status updated", "status", cl.Status)
	if capiCl.Status.ControlPlaneReady && capiCl.Status.InfrastructureReady {
		cl = appv1alpha1.ClusterReady(cl)
		err := kube.EnsureComponentsConfig(ctx, r.Client, cl)
		if err != nil {
			return ctrl.Result{}, err
		}
		if cl.Status.ConciergeInfo == nil {
			err = r.reconcileAuthEndpoint(ctx, cl)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}
	cl = appv1alpha1.ClusterNotReady(cl, meta.WaitProvisionReason, "Waiting cluster to be provisioned")
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *ClusterReconciler) applyObjects(ctx context.Context, cl *appv1alpha1.Cluster, objs []unstructured.Unstructured) error {
	for _, o := range objs {
		if o.GetAPIVersion() == capi.GroupVersion.String() && o.GetKind() == "Cluster" {
			err := ctrl.SetControllerReference(cl, &o, scheme.Scheme)
			if err != nil {
				return err
			}
		}
		err := retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
			labels := o.GetLabels()
			if labels == nil {
				labels = make(map[string]string)
			}
			labels[meta.LabelUndistro] = ""
			labels[meta.LabelUndistroClusterName] = cl.Name
			labels[capi.ClusterLabelName] = cl.Name
			o.SetLabels(labels)
			_, err := util.CreateOrUpdate(ctx, r.Client, &o)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ClusterReconciler) fillTemplateVariables(ctx context.Context, c client.Client, cl *appv1alpha1.Cluster) (map[string]interface{}, error) {
	log := setDefaultLog(ctx)

	vars := make(map[string]interface{})
	v := make(map[string]interface{})
	err := template.SetVariablesFromEnvVar(ctx, template.VariablesInput{
		ClientSet:      r.Client,
		NamespacedName: client.ObjectKeyFromObject(cl),
		Variables:      v,
		EnvVars:        cl.Spec.InfrastructureProvider.Env,
	})
	if err != nil {
		return nil, err
	}

	vars["Cluster"] = cl
	vars["ENV"] = v
	acc, err := cloud.GetAccount(ctx, c, cl)
	if err != nil {
		return nil, err
	}
	vars["Account"] = acc

	labels := cl.GetLabels()
	_, moved := labels[meta.LabelUndistroMoved]
	validDiff := true
	isNotValidDiff := moved && !cl.Spec.InfrastructureProvider.IsManaged() && cl.Status.LastUsedUID == ""

	log.Info("Checking if is a valid diff", "lastUserUUID", cl.Status.LastUsedUID)
	if isNotValidDiff {
		validDiff = false
		cp := capicp.KubeadmControlPlane{}
		err := r.Get(ctx, client.ObjectKeyFromObject(cl), &cp)
		if client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		if !apierrors.IsNotFound(err) {
			split := strings.Split(cp.Spec.MachineTemplate.InfrastructureRef.Name, "-")
			cl.Status.LastUsedUID = split[len(split)-1]
		}
		log.Info("Value of the last used UUID", "lastUserUUID", cl.Status.LastUsedUID)
	}

	if r.hasDiff(ctx, cl) && validDiff {
		cpChanged, workersChanged := r.machineTypeChanged(ctx, cl)
		log.Info("Changes in cluster specification", "controlPlaneChanged", cpChanged, "workersChanged", workersChanged)

		newUUID := string(uuid.NewUUID())
		log.Info("New UUID", "newUUID", newUUID)
		if cpChanged || len(workersChanged) > 0 {
			vars["CPID"] = newUUID
			cl.Status.LastUsedUID = newUUID
		} else {
			vars["CPID"] = cl.Status.LastUsedUID
		}

		log.Info("Control Plane UUID", "controlPlaneUUID", vars["CPID"])
		vars["WorkersChanged"] = workersChanged
		vars["OldID"] = cl.Status.LastUsedUID
		log.Info("Template variables", "vars", vars)
	}
	return vars, nil
}

func (r *ClusterReconciler) tryRetrieveBastionIP(ctx context.Context, capiCluster capi.Cluster) (string, error) {
	ref := capiCluster.Spec.InfrastructureRef
	if ref != nil {
		key := client.ObjectKey{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		}
		o := unstructured.Unstructured{}
		o.SetGroupVersionKind(ref.GroupVersionKind())
		err := r.Get(ctx, key, &o)
		if err != nil {
			return "", client.IgnoreNotFound(err)
		}
		ip, _, err := unstructured.NestedString(o.Object, "status", "bastion", "publicIp")
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return "", nil
}

func (r *ClusterReconciler) reconcileAuthEndpoint(ctx context.Context, cl *appv1alpha1.Cluster) error {
	clObj := client.ObjectKey{
		Name:      cl.Name,
		Namespace: cl.GetNamespace(),
	}
	cfg, err := kube.GetInternalKubeconfig(ctx, r.Client, clObj)
	if err != nil {
		return err
	}
	info, err := kube.ConciergeInfoFromConfig(ctx, cfg)
	if err != nil {
		return kube.IgnoreConciergeNotInstalled(err)
	}
	cl.Status.ConciergeInfo = info
	return nil
}

func (r *ClusterReconciler) machineTypeChanged(ctx context.Context, cl *appv1alpha1.Cluster) (bool, []int) {
	log := setDefaultLog(ctx)

	cpChanged := false
	workersChanged := make([]int, 0)
	if !cl.Spec.InfrastructureProvider.IsManaged() && cl.Spec.ControlPlane != nil {
		cpChanged = cl.Spec.ControlPlane.MachineType != cl.Status.ControlPlane.MachineType
	}
	log.Info("Old vs new control plane template comparison", "controlPlaneSpec", cl.Spec.ControlPlane, "controlPlaneStatus", cl.Status.ControlPlane)
	for i, w := range cl.Spec.Workers {
		if len(cl.Status.Workers)-1 >= i {
			if w.MachineType != cl.Status.Workers[i].MachineType {
				workersChanged = append(workersChanged, i)
			}
		} else {
			workersChanged = append(workersChanged, i)
		}
	}
	return cpChanged, workersChanged
}

func (r *ClusterReconciler) hasDiff(ctx context.Context, cl *appv1alpha1.Cluster) bool {
	log := setDefaultLog(ctx)

	if cl.Spec.KubernetesVersion != cl.Status.KubernetesVersion {
		log.Info(KubernetesVersionChanged, "old", cl.Status.KubernetesVersion, "new", cl.Spec.KubernetesVersion)
		return true
	}

	if r.hasDiffInControlPlane(log, cl) {
		return true
	}

	if r.hasDiffInWorkers(log, cl) {
		return true
	}

	return !reflect.DeepEqual(cl.Spec.Bastion, cl.Status.BastionConfig)
}

func (r *ClusterReconciler) hasDiffInControlPlane(log logr.Logger, cl *appv1alpha1.Cluster) bool {
	if !cl.Spec.InfrastructureProvider.IsManaged() && cl.Spec.ControlPlane != nil {
		log.Info("Checking changes in control plane replicas")
		if *cl.Spec.ControlPlane.Replicas != *cl.Status.ControlPlane.Replicas {
			log.Info(ControlPlaneReplicasChanged, "old", cl.Status.ControlPlane.Replicas, "new", cl.Spec.ControlPlane.Replicas)
			return true
		}
		if cl.Spec.ControlPlane.MachineType != cl.Status.ControlPlane.MachineType {
			log.Info(ControlPlaneMachineTypeChanged, "old", cl.Status.ControlPlane.MachineType, "new", cl.Spec.ControlPlane.MachineType)
			return true
		}
		if !reflect.DeepEqual(cl.Spec.ControlPlane.Labels, cl.Status.ControlPlane.Labels) {
			log.Info(ControlPlaneLabelsChanged, "old", cl.Status.ControlPlane.Labels, "new", cl.Spec.ControlPlane.Labels)
			return true
		}
		if !reflect.DeepEqual(cl.Spec.ControlPlane.Taints, cl.Status.ControlPlane.Taints) {
			log.Info(ControlPlaneTaintsChanged, "old", cl.Status.ControlPlane.Taints, "new", cl.Spec.ControlPlane.Taints)
			return true
		}
		if !reflect.DeepEqual(cl.Spec.ControlPlane.ProviderTags, cl.Status.ControlPlane.ProviderTags) {
			log.Info(ControlPlaneProviderTagsChanged, "old", cl.Status.ControlPlane.ProviderTags, "new", cl.Spec.ControlPlane.ProviderTags)
			return true
		}
	}
	return false
}

func (r *ClusterReconciler) hasDiffInWorkers(log logr.Logger, cl *appv1alpha1.Cluster) bool {
	if len(cl.Spec.Workers) != len(cl.Status.Workers) {
		log.Info(WorkersChanged, "old", cl.Status.Workers, "new", cl.Spec.Workers)
		return true
	}

	for i, w := range cl.Spec.Workers {
		if *w.Replicas != *cl.Status.Workers[i].Replicas {
			log.Info(WorkersReplicasChanged, "old", cl.Status.Workers[i].Replicas, "new", w.Replicas)
			return true
		}
		if w.MachineType != cl.Status.Workers[i].MachineType {
			log.Info(WorkerMachineTypeChanged, "old", cl.Status.Workers[i].MachineType, "new", w.MachineType)
			return true
		}
		if !reflect.DeepEqual(w.Labels, cl.Status.Workers[i].Labels) {
			log.Info(WorkerLabelsChanged, "old", cl.Status.Workers[i].Labels, "new", w.Labels)
			return true
		}
		if !reflect.DeepEqual(w.Taints, cl.Status.Workers[i].Taints) {
			log.Info(WorkerTaintsChanged, "old", cl.Status.Workers[i].Taints, "new", w.Taints)
			return true
		}
		if !reflect.DeepEqual(w.ProviderTags, cl.Status.Workers[i].ProviderTags) {
			log.Info(WorkerProviderTagsChanged, "old", cl.Status.Workers[i].ProviderTags, "new", w.ProviderTags)
			return true
		}
		if !reflect.DeepEqual(w.Autoscale, cl.Status.Workers[i].Autoscale) {
			log.Info(WorkerAutoscalingChanged, "old", cl.Status.Workers[i].Autoscale, "new", w.Autoscale)
			return true
		}
	}
	return false
}

func (r *ClusterReconciler) reconcileCNI(ctx context.Context, cl *appv1alpha1.Cluster) error {
	log := setDefaultLog(ctx)

	const (
		cniCalicoName = "calico"
		calicoVersion = "3.19.1"
	)
	log.Info("Reconciling CNI")

	key := client.ObjectKey{
		Name:      hr.GetObjectName(cniCalicoName, cl.Name),
		Namespace: cl.GetNamespace(),
	}
	release := appv1alpha1.HelmRelease{}
	err := r.Get(ctx, key, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	if meta.InReadyCondition(release.Status.Conditions) {
		meta.SetResourceCondition(cl, meta.CNIInstalledCondition, metav1.ConditionTrue, meta.CNIInstalledSuccessReason, "calico installed")
	}

	calicoValues := cloud.CalicoValues(cl)
	release, err = hr.Prepare(cniCalicoName, "kube-system", cl.GetNamespace(), calicoVersion, cl.Name, calicoValues)
	if err != nil {
		return err
	}

	if release.Labels == nil {
		release.Labels = make(map[string]string)
	}
	release.Labels[meta.LabelUndistroMove] = ""
	if release.Annotations == nil {
		release.Annotations = make(map[string]string)
	}
	release.Annotations[meta.SetupAnnotation] = cniCalicoName

	err = hr.Install(ctx, r.Client, log, release, cl)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClusterReconciler) reconcileClusterAutoscaler(ctx context.Context, cl *appv1alpha1.Cluster) error {
	log := setDefaultLog(ctx)

	const (
		chartName    = "cluster-autoscaler"
		chartVersion = "9.10.9"
	)
	log.Info("Reconciling Cluster Autoscaler")

	key := client.ObjectKey{
		Name:      hr.GetObjectName(chartName, cl.Name),
		Namespace: cl.GetNamespace(),
	}
	release := appv1alpha1.HelmRelease{}
	err := r.Get(ctx, key, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	if meta.InReadyCondition(release.Status.Conditions) {
		meta.SetResourceCondition(cl, meta.ClusterAutoscalerInstalledCondition, metav1.ConditionTrue, meta.ClusterAutoscalerInstalledSuccessReason, "cluster-autoscaler installed")
	}

	values := map[string]interface{}{
		"clusterName": cl.Name,
	}
	release, err = hr.Prepare(chartName, cl.GetNamespace(), cl.GetNamespace(), chartVersion, cl.Name, values)
	if err != nil {
		return err
	}

	if release.Labels == nil {
		release.Labels = make(map[string]string)
	}
	release.Labels[meta.LabelUndistroMove] = ""
	if release.Annotations == nil {
		release.Annotations = make(map[string]string)
	}
	release.Annotations[meta.HelmReleaseLocation] = ""

	err = hr.Install(ctx, r.Client, log, release, cl)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClusterReconciler) reconcileDelete(ctx context.Context, uCluster *appv1alpha1.Cluster) (ctrl.Result, error) {
	log := setDefaultLog(ctx)

	log.Info("Cluster is under deletion")
	capiCluster := capi.Cluster{}
	err := r.Get(ctx, client.ObjectKeyFromObject(uCluster), &capiCluster)
	if apierrors.IsNotFound(err) {
		if controllerutil.ContainsFinalizer(uCluster, meta.Finalizer) {
			controllerutil.RemoveFinalizer(uCluster, meta.Finalizer)
		}
		// update is done by patch helper in the main reconcile func
		return ctrl.Result{}, nil
	}

	log.Info("Deleting CAPI Cluster")
	err = r.Delete(ctx, &capiCluster)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.removeDeps(ctx, *uCluster)
	if err != nil {
		return ctrl.Result{}, err
	}
	if controllerutil.ContainsFinalizer(uCluster, meta.Finalizer) {
		controllerutil.RemoveFinalizer(uCluster, meta.Finalizer)
	}
	// update is done by patch helper in the main reconcile func
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) removeDeps(ctx context.Context, undistroCluster appv1alpha1.Cluster) error {
	log := setDefaultLog(ctx)
	log.Info("Removing UnDistro Cluster dependencies")
	releaseClusterName := fmt.Sprintf("%s/%s", undistroCluster.GetNamespace(), undistroCluster.Name)
	releaseList := appv1alpha1.HelmReleaseList{}
	err := r.List(ctx, &releaseList)
	if err != nil {
		return err
	}
	log.Info("Removing cluster releases")
	for _, item := range releaseList.Items {
		if item.Spec.ClusterName == releaseClusterName {
			err = r.Delete(ctx, &item)
			if err != nil {
				return err
			}
		}
	}
	policies := appv1alpha1.DefaultPoliciesList{}
	err = r.List(ctx, &policies, client.InNamespace(undistroCluster.GetNamespace()))
	if err != nil {
		return err
	}
	log.Info("Removing cluster policies")
	for _, item := range policies.Items {
		if item.Spec.ClusterName == undistroCluster.Name {
			err = r.Delete(ctx, &item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ClusterReconciler) capiToUndistro(o client.Object) []ctrl.Request {
	capiCluster := o.(*capi.Cluster).DeepCopy()
	if capiCluster.Status.Phase == "" {
		return nil
	}
	return []ctrl.Request{
		{
			NamespacedName: client.ObjectKeyFromObject(capiCluster),
		},
	}
}

func setDefaultLog(ctx context.Context) logr.Logger {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}
	return log
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Cluster{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Watches(
			&source.Kind{
				Type: &capi.Cluster{},
			},
			handler.EnqueueRequestsFromMapFunc(r.capiToUndistro),
		).
		Complete(r)
}
