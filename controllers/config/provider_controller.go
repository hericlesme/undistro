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

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	configv1alpha1 "github.com/getupio-undistro/undistro/apis/config/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/predicate"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ProviderReconciler reconciles a Provider object
type ProviderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	p := configv1alpha1.Provider{}
	if err := r.Get(ctx, req.NamespacedName, &p); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log := r.Log.WithValues("provider", req.NamespacedName)
	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&p, meta.Finalizer) {
		controllerutil.AddFinalizer(&p, meta.Finalizer)
		if err := r.Update(ctx, &p); err != nil {
			log.Error(err, "unable to register finalizer")
			return ctrl.Result{}, err
		}
	}
	if p.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}
	if !p.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, log, p)
	}
	p, result, err := r.reconcile(ctx, log, p)
	// Update status after reconciliation.
	if updateStatusErr := r.patchStatus(ctx, &p); updateStatusErr != nil {
		log.Error(updateStatusErr, "unable to update status after reconciliation")
		return ctrl.Result{Requeue: true}, updateStatusErr
	}
	return result, err
}

func (r *ProviderReconciler) reconcileDelete(ctx context.Context, logger logr.Logger, p configv1alpha1.Provider) (ctrl.Result, error) {
	key := client.ObjectKey{
		Name:      p.Status.HelmReleaseName,
		Namespace: p.GetNamespace(),
	}
	hr := appv1alpha1.HelmRelease{}
	err := r.Get(ctx, key, &hr)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if apierrors.IsNotFound(err) {
		// Remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(&p, meta.Finalizer)
		if err := r.Update(ctx, &p); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{Requeue: true}, nil
}

func (r *ProviderReconciler) reconcile(ctx context.Context, logger logr.Logger, p configv1alpha1.Provider) (configv1alpha1.Provider, ctrl.Result, error) {
	return p, ctrl.Result{}, nil
}

func (r *ProviderReconciler) patchStatus(ctx context.Context, p *configv1alpha1.Provider) error {
	latest := &configv1alpha1.Provider{}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(p), latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, p, client.MergeFrom(latest))
}

func (r *ProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.Provider{}, builder.WithPredicates(predicate.Changed{})).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(r)
}
