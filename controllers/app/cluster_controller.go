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
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/predicate"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"
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
		if err := r.Update(ctx, &cl); err != nil {
			log.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}
	if cl.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}
	capiCluster := capi.Cluster{}
	err := r.Get(ctx, client.ObjectKeyFromObject(&cl), &capiCluster)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if !cl.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, log, cl, capiCluster)
	}
	cl, result, err := r.reconcile(ctx, log, cl, capiCluster)
	if updateStatusErr := r.patchStatus(ctx, &cl); updateStatusErr != nil {
		log.Error(updateStatusErr, "unable to update status after reconciliation")
		return ctrl.Result{Requeue: true}, updateStatusErr
	}
	return result, err
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
		if updateStatusErr := r.patchStatus(ctx, &cl); updateStatusErr != nil {
			log.Error(updateStatusErr, "unable to update status after generation update")
			return cl, ctrl.Result{Requeue: true}, updateStatusErr
		}
	}
	for _, cond := range cl.Status.Conditions {
		meta.SetResourceCondition(&cl, cond.Type, cond.Status, cond.Reason, cond.Message)
	}
	if capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseProvisioning {
		if capiCluster.Status.ControlPlaneInitialized && !capiCluster.Status.ControlPlaneReady && !cl.Spec.InfrastructureCluster.Managed {
			err := r.installCNI(ctx, cl)
			if err != nil {
				meta.SetResourceCondition(&cl, meta.CNIInstalledCondition, metav1.ConditionFalse, meta.CNIInstalledFailedReason, err.Error())
				return cl, ctrl.Result{}, err
			}
			meta.SetResourceCondition(&cl, meta.CNIInstalledCondition, metav1.ConditionTrue, meta.CNIInstalledSuccessReason, "calico installed")
		}
		return appv1alpha1.ClusterNotReady(cl, meta.WaitProvisionReason, "wait cluster to be provisioned"), ctrl.Result{Requeue: true}, nil
	}
	if !reflect.DeepEqual(*capiCluster.Spec.ClusterNetwork, capi.ClusterNetwork{}) && reflect.DeepEqual(cl.Spec.Network.ClusterNetwork, capi.ClusterNetwork{}) {
		cl.Spec.Network.ClusterNetwork = *capiCluster.Spec.ClusterNetwork
		if err := r.Update(ctx, &cl); err != nil {
			return cl, ctrl.Result{Requeue: true}, err
		}
	}
	if !reflect.DeepEqual(capiCluster.Spec.ControlPlaneEndpoint, capi.APIEndpoint{}) && reflect.DeepEqual(cl.Spec.ControlPlane.Endpoint, capi.APIEndpoint{}) {
		cl.Spec.ControlPlane.Endpoint = capiCluster.Spec.ControlPlaneEndpoint
		if err := r.Update(ctx, &cl); err != nil {
			return cl, ctrl.Result{Requeue: true}, err
		}
	}
	return cl, ctrl.Result{}, nil
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
	for _, o := range objs {
		_, err = util.CreateOrUpdate(ctx, r.Client, &o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ClusterReconciler) patchStatus(ctx context.Context, cl *appv1alpha1.Cluster) error {
	latest := &appv1alpha1.HelmRelease{}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(cl), latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, cl, client.MergeFrom(latest))
}

func (r *ClusterReconciler) reconcileDelete(ctx context.Context, logger logr.Logger, cl appv1alpha1.Cluster, capiCluster capi.Cluster) (ctrl.Result, error) {
	if capiCluster.Status.GetTypedPhase() != capi.ClusterPhaseUnknown && capiCluster.Status.GetTypedPhase() != capi.ClusterPhaseDeleting {
		return ctrl.Result{Requeue: true}, r.Delete(ctx, &capiCluster)
	}
	if capiCluster.Status.GetTypedPhase() == capi.ClusterPhaseDeleting {
		return ctrl.Result{Requeue: true}, nil
	}
	controllerutil.RemoveFinalizer(&cl, meta.Finalizer)
	if err := r.Update(ctx, &cl); err != nil {
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
