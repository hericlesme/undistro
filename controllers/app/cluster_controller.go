/*
Copyright 2020 The UnDistro authors

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

package controllers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/predicate"
	"github.com/getupio-undistro/undistro/pkg/retry"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/template"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/uuid"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cl := appv1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, &cl); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log := r.Log.WithValues("cluster", req.NamespacedName)

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&cl, meta.Finalizer) {
		controllerutil.AddFinalizer(&cl, meta.Finalizer)
		_, err := util.CreateOrUpdate(ctx, r.Client, &cl)
		if err != nil {
			log.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}
	if cl.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}
	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(&cl, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	patchOpts := []patch.Option{}
	capiCluster := capi.Cluster{}
	err = r.Get(ctx, client.ObjectKeyFromObject(&cl), &capiCluster)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if !cl.DeletionTimestamp.IsZero() {
		log.Info("Deleting")
		return r.reconcileDelete(ctx, log, cl, capiCluster)
	}
	cl, result, err := r.reconcile(ctx, log, cl, capiCluster)
	if err == nil {
		patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
	}
	patchErr := patchHelper.Patch(ctx, &cl, patchOpts...)
	if err != nil {
		err = kerrors.NewAggregate([]error{patchErr, err})
	}
	return result, err
}

func (r *ClusterReconciler) templateVariables(ctx context.Context, capiCluster *capi.Cluster, cl *appv1alpha1.Cluster, hasDiff bool) (map[string]interface{}, error) {
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
	if hasDiff {
		cl.Status.LastUsedUID = string(uuid.NewUUID())
	}
	return vars, nil
}

func (r *ClusterReconciler) reconcileMachines(ctx context.Context, capiCluster *capi.Cluster, cl *appv1alpha1.Cluster) (bool, error) {
	cpMachines, wMachines, err := util.GetMachinesForCluster(ctx, r.Client, capiCluster)
	if err != nil {
		return false, err
	}
	if len(cpMachines.Items) == 0 && len(wMachines.Items) == 0 {
		return true, nil
	}
	for _, wm := range wMachines.Items {
		if wm.Spec.Version != nil {
			if *wm.Spec.Version != cl.Spec.KubernetesVersion {
				return true, nil
			}
		}
	}
	if cl.Spec.ControlPlane != nil {
		cpProvider, err := util.GetProviderMachinesUnstructured(ctx, r.Client, cpMachines)
		if err != nil {
			return false, err
		}
		for _, o := range cpProvider.Items {
			mt, ok, err := unstructured.NestedString(o.Object, "spec", "instanceType")
			if err != nil {
				return false, err
			}
			if !ok {
				continue
			}
			if mt != cl.Spec.ControlPlane.MachineType {
				return true, nil
			}
		}
	}
	acceptableTypes := make([]string, len(cl.Spec.Workers))
	for i := range cl.Spec.Workers {
		acceptableTypes[i] = cl.Spec.Workers[i].MachineType
	}
	typesSet := sets.NewString(acceptableTypes...)
	wProvider, err := util.GetProviderMachinesUnstructured(ctx, r.Client, wMachines)
	if err != nil {
		return false, err
	}
	for _, o := range wProvider.Items {
		mt, ok, err := unstructured.NestedString(o.Object, "spec", "instanceType")
		if err != nil {
			return false, err
		}
		if !ok {
			continue
		}
		if !util.ContainsStringInSlice(typesSet.List(), mt) {
			return true, nil
		}
	}
	return false, nil
}

func (r *ClusterReconciler) getBastionIP(ctx context.Context, log logr.Logger, cl appv1alpha1.Cluster, capiCluster capi.Cluster) (string, error) {
	ref := capiCluster.Spec.InfrastructureRef
	if cl.Spec.InfrastructureProvider.Managed {
		ref = capiCluster.Spec.ControlPlaneRef
	}
	if ref != nil {
		key := client.ObjectKey{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		}
		o := unstructured.Unstructured{}
		o.SetGroupVersionKind(ref.GroupVersionKind())
		err := r.Get(ctx, key, &o)
		if err != nil {
			return "", err
		}
		ip, _, err := unstructured.NestedString(o.Object, "status", "bastion", "publicIp")
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return "", nil
}

func (r *ClusterReconciler) reconcile(ctx context.Context, log logr.Logger, cl appv1alpha1.Cluster, capiCluster capi.Cluster) (appv1alpha1.Cluster, ctrl.Result, error) {
	if cl.Status.ObservedGeneration != cl.Generation {
		cl.Status.ObservedGeneration = cl.Generation
		cl = appv1alpha1.ClusterProgressing(cl)
		cl.Status.TotalWorkerPools = int32(len(cl.Spec.Workers))
		cl.Status.TotalWorkerReplicas = 0
		for _, w := range cl.Spec.Workers {
			cl.Status.TotalWorkerReplicas += *w.Replicas
		}
		if _, updateStatusErr := util.CreateOrUpdate(ctx, r.Client, &cl); updateStatusErr != nil {
			log.Error(updateStatusErr, "unable to update status after generation update")
			return cl, ctrl.Result{Requeue: true}, updateStatusErr
		}
	}
	for _, cond := range cl.Status.Conditions {
		meta.SetResourceCondition(&cl, cond.Type, cond.Status, cond.Reason, cond.Message)
	}
	if capiCluster.Status.ControlPlaneInitialized && !capiCluster.Status.ControlPlaneReady && !cl.Spec.InfrastructureProvider.Managed {
		log.Info("installing calico")
		err := r.installCNI(ctx, cl)
		if err != nil {
			meta.SetResourceCondition(&cl, meta.CNIInstalledCondition, metav1.ConditionFalse, meta.CNIInstalledFailedReason, err.Error())
			return cl, ctrl.Result{}, err
		}
		meta.SetResourceCondition(&cl, meta.CNIInstalledCondition, metav1.ConditionTrue, meta.CNIInstalledSuccessReason, "calico installed")
	}
	if cl.Spec.Bastion != nil {
		if cl.Spec.Bastion.Enabled && cl.Status.BastionPublicIP == "" {
			var err error
			cl.Status.BastionPublicIP, err = r.getBastionIP(ctx, log, cl, capiCluster)
			if err != nil {
				return appv1alpha1.ClusterNotReady(cl, meta.WaitProvisionReason, err.Error()), ctrl.Result{Requeue: true}, nil
			}
		}
	}
	if capiCluster.Spec.ClusterNetwork != nil {
		if !reflect.DeepEqual(*capiCluster.Spec.ClusterNetwork, capi.ClusterNetwork{}) && reflect.DeepEqual(cl.Spec.Network.ClusterNetwork, capi.ClusterNetwork{}) {
			cl.Spec.Network.ClusterNetwork = *capiCluster.Spec.ClusterNetwork
			_, err := util.CreateOrUpdate(ctx, r.Client, &cl)
			if err != nil {
				return cl, ctrl.Result{Requeue: true}, err
			}
		}
	}
	if !reflect.DeepEqual(capiCluster.Spec.ControlPlaneEndpoint, capi.APIEndpoint{}) && reflect.DeepEqual(cl.Spec.ControlPlane.Endpoint, capi.APIEndpoint{}) {
		cl.Spec.ControlPlane.Endpoint = capiCluster.Spec.ControlPlaneEndpoint
		_, err := util.CreateOrUpdate(ctx, r.Client, &cl)
		if err != nil {
			return cl, ctrl.Result{Requeue: true}, err
		}
	}
	hasDiff, err := r.reconcileMachines(ctx, &capiCluster, &cl)
	if err != nil {
		return appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error()), ctrl.Result{Requeue: true}, err
	}

	if capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseProvisioning {
		return appv1alpha1.ClusterNotReady(cl, meta.WaitProvisionReason, "wait cluster to be provisioned"), ctrl.Result{Requeue: true}, nil
	}
	if !hasDiff && capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseProvisioned && capiCluster.Status.ControlPlaneReady {
		return appv1alpha1.ClusterReady(cl), ctrl.Result{}, nil
	}
	if hasDiff && capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseProvisioned && capiCluster.Status.ControlPlaneReady && !cl.Spec.InfrastructureProvider.Managed {
		meta.SetResourceCondition(&cl, meta.ReadyCondition, metav1.ConditionTrue, meta.ReconciliationSucceededReason, "updating")
	}
	if !hasDiff && capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseProvisioned && !capiCluster.Status.ControlPlaneReady {
		return appv1alpha1.ClusterNotReady(cl, meta.ArtifactFailedReason, "wait control plane to be ready"), ctrl.Result{Requeue: true}, nil
	}
	vars, err := r.templateVariables(ctx, &capiCluster, &cl, hasDiff)
	if err != nil {
		return appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error()), ctrl.Result{Requeue: true}, err
	}
	tpl := template.New(template.Options{
		Directory: "clustertemplates",
	})
	buff := &bytes.Buffer{}
	err = tpl.YAML(buff, cl.Spec.InfrastructureProvider.Name, vars)
	if err != nil {
		return appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error()), ctrl.Result{Requeue: true}, err
	}
	objs, err := util.ToUnstructured(buff.Bytes())
	if err != nil {
		return appv1alpha1.ClusterNotReady(cl, meta.TemplateAppliedFailed, err.Error()), ctrl.Result{Requeue: true}, err
	}
	for _, o := range objs {
		err = retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
			_, err = util.CreateOrUpdate(ctx, r.Client, &o)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return cl, ctrl.Result{Requeue: true}, err
		}
	}
	if capiCluster.Status.InfrastructureReady && capiCluster.Status.ControlPlaneReady {
		return appv1alpha1.ClusterReady(cl), ctrl.Result{}, nil
	}
	return cl, ctrl.Result{Requeue: true}, nil
}

func (r *ClusterReconciler) installCNI(ctx context.Context, cl appv1alpha1.Cluster) error {
	const addr = "https://docs.projectcalico.org/manifests/calico.yaml"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
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
	buff := bytes.Buffer{}
	_, err = io.Copy(&buff, resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unable to get calico: %s", buff.String())
	}
	objs, err := util.ToUnstructured(buff.Bytes())
	if err != nil {
		return err
	}
	byt, err := kubeconfig.FromSecret(ctx, r.Client, client.ObjectKeyFromObject(&cl))
	if err != nil {
		return err
	}
	restGetter := kube.NewMemoryRESTClientGetter(byt, cl.GetNamespace())
	workloadClientConfig, err := restGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	workloadClient, err := client.New(workloadClientConfig, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	for _, o := range objs {
		_, err = util.CreateOrUpdate(ctx, workloadClient, &o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ClusterReconciler) reconcileDelete(ctx context.Context, logger logr.Logger, cl appv1alpha1.Cluster, capiCluster capi.Cluster) (ctrl.Result, error) {
	if capiCluster.Status.GetTypedPhase() != capi.ClusterPhaseUnknown && capiCluster.Status.GetTypedPhase() != capi.ClusterPhaseDeleting {
		return ctrl.Result{Requeue: true}, r.Delete(ctx, &capiCluster)
	}
	if capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseDeleting {
		return ctrl.Result{Requeue: true}, nil
	}
	controllerutil.RemoveFinalizer(&cl, meta.Finalizer)
	_, err := util.CreateOrUpdate(ctx, r.Client, &cl)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	err = r.Delete(ctx, &cl)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Cluster{}, builder.WithPredicates(predicate.Changed{})).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(r)
}
