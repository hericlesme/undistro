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
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/predicate"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	jobOwnerKey = ".metadata.controller"
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

	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &capi.Cluster{}, jobOwnerKey, func(rawObj client.Object) []string {
		cluster := rawObj.(*capi.Cluster)
		owner := metav1.GetControllerOf(cluster)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != appv1alpha1.GroupVersion.String() || owner.Kind != "Cluster" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Cluster{}, builder.WithPredicates(predicate.Changed{})).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(r)
}
