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
	"context"
	"fmt"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/helm"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/predicate"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// HelmReleaseReconciler reconciles a HelmRelease object
type HelmReleaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	config *rest.Config
}

// +kubebuilder:rbac:groups=app.undistro.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.undistro.io,resources=helmreleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.undistro.io,resources=helmreleases/finalizers,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *HelmReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	hr := appv1alpha1.HelmRelease{}
	if err := r.Get(ctx, req.NamespacedName, &hr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log := r.Log.WithValues("helmrelease", req.NamespacedName)
	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&hr, meta.Finalizer) {
		controllerutil.AddFinalizer(&hr, meta.Finalizer)
		if err := r.Update(ctx, &hr); err != nil {
			log.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}
	if hr.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	if !hr.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, log, hr)
	}
	hr, result, err := r.reconcile(ctx, log, hr)
	// Update status after reconciliation.
	if updateStatusErr := r.patchStatus(ctx, &hr); updateStatusErr != nil {
		log.Error(updateStatusErr, "unable to update status after reconciliation")
		return ctrl.Result{Requeue: true}, updateStatusErr
	}
	return result, err
}

func (r *HelmReleaseReconciler) reconcile(ctx context.Context, log logr.Logger, hr appv1alpha1.HelmRelease) (appv1alpha1.HelmRelease, ctrl.Result, error) {
	if hr.Status.ObservedGeneration != hr.Generation {
		hr.Status.ObservedGeneration = hr.Generation
		hr = appv1alpha1.HelmReleaseProgressing(hr)
		if updateStatusErr := r.patchStatus(ctx, &hr); updateStatusErr != nil {
			log.Error(updateStatusErr, "unable to update status after generation update")
			return hr, ctrl.Result{Requeue: true}, updateStatusErr
		}
	}
	if hr.Spec.SecretRef != nil {
		name := types.NamespacedName{
			Name:      hr.Spec.SecretRef.Name,
			Namespace: hr.GetNamespace(),
		}
		var secret corev1.Secret
		err := r.Client.Get(ctx, name, &secret)
		if err != nil {
			return hr, ctrl.Result{}, err
		}
	}
	return hr, ctrl.Result{}, nil
}

func (r *HelmReleaseReconciler) patchStatus(ctx context.Context, hr *appv1alpha1.HelmRelease) error {
	latest := &appv1alpha1.HelmRelease{}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(hr), latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, hr, client.MergeFrom(latest))
}

func (r *HelmReleaseReconciler) getRESTClientGetter(ctx context.Context, hr appv1alpha1.HelmRelease) (genericclioptions.RESTClientGetter, error) {
	if hr.Spec.ClusterName == "" {
		return kube.NewInClusterRESTClientGetter(r.config, hr.GetNamespace()), nil
	}
	secretName := types.NamespacedName{
		Namespace: hr.GetNamespace(),
		Name:      fmt.Sprintf("%s-kubeconfig", hr.Spec.ClusterName),
	}
	var secret corev1.Secret
	if err := r.Get(ctx, secretName, &secret); err != nil {
		return nil, fmt.Errorf("could not find KubeConfig secret '%s': %w", secretName, err)
	}
	kubeConfig, ok := secret.Data["value"]
	if !ok {
		return nil, fmt.Errorf("KubeConfig secret '%s' does not contain a 'value' key", secretName)
	}
	return kube.NewMemoryRESTClientGetter(kubeConfig, hr.GetNamespace()), nil
}

// reconcileDelete deletes the v1beta1.HelmChart of the v2beta1.HelmRelease,
// and uninstalls the Helm release if the resource has not been suspended.
func (r *HelmReleaseReconciler) reconcileDelete(ctx context.Context, logger logr.Logger, hr appv1alpha1.HelmRelease) (ctrl.Result, error) {
	restClient, err := r.getRESTClientGetter(ctx, hr)
	if err != nil {
		return ctrl.Result{}, err
	}
	runner, err := helm.NewRunner(restClient, hr.GetNamespace(), logger)
	if err != nil {
		return ctrl.Result{}, err
	}
	rel, err := runner.ObserveLastRelease(hr)
	if err != nil {
		return ctrl.Result{}, err
	}
	if rel == nil {
		return ctrl.Result{}, nil
	}
	err = runner.Uninstall(hr)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(&hr, meta.Finalizer)
	if err := r.Update(ctx, &hr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.config = mgr.GetConfig()
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.HelmRelease{}, builder.WithPredicates(predicate.Changed{})).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(r)
}
