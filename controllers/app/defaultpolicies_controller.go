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
	"path/filepath"
	"time"

	"github.com/getupio-undistro/controllerlib"
	"github.com/getupio-undistro/meta"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/fs"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// DefaultPoliciesReconciler reconciles a DefaultPolicies object
type DefaultPoliciesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	kyvernoReleaseName    = "kyverno"
	kyvernoReleaseVersion = "1.4.2"
)

func (r *DefaultPoliciesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()

	p := appv1alpha1.DefaultPolicies{}
	if err := r.Get(ctx, req.NamespacedName, &p); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}
	log.WithValues("DefaultPolicies", req.NamespacedName)

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(&p, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer controllerlib.PatchInstance(ctx, controllerlib.InstanceOpts{
		Controller: "DefaultPoliciesController",
		Request:    req.String(),
		Object:     &p,
		Error:      err,
		Helper:     patchHelper,
	})

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&p, meta.Finalizer) {
		log.Info("Adding finalizer")
		controllerutil.AddFinalizer(&p, meta.Finalizer)
		return ctrl.Result{}, nil
	}

	if p.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		p = appv1alpha1.DefaultPoliciesPaused(p)
		return ctrl.Result{}, nil
	}

	if p.Generation < p.Status.ObservedGeneration {
		log.Info("skipping this old version of reconciled object")
		return ctrl.Result{}, nil
	}

	if !p.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, p)
	}

	// if cluster name is empty, we install the chart only for the undistro cluster
	cl := &appv1alpha1.Cluster{}
	if !util.IsMgmtCluster(p.Spec.ClusterName) {
		key := client.ObjectKey{
			Name:      p.Spec.ClusterName,
			Namespace: p.GetNamespace(),
		}
		err = r.Get(ctx, key, cl)
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
	}

	p, result, err := r.reconcile(ctx, p, cl)
	durationMsg := fmt.Sprintf("Reconcilation finished in %s", time.Since(start).String())
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	log.Info(durationMsg)
	return result, err
}

func (r *DefaultPoliciesReconciler) reconcile(ctx context.Context, p appv1alpha1.DefaultPolicies, cl *appv1alpha1.Cluster) (appv1alpha1.DefaultPolicies, ctrl.Result, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	values := map[string]interface{}{
		"fullnameOverride": kyvernoReleaseName,
		"namespace":        kyvernoReleaseName,
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2Gi",
			},
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "500Mi",
			},
		},
	}

	if cl.HasInfraNodes() {
		infraValues := map[string]interface{}{
			"nodeSelector": map[string]interface{}{
				meta.LabelUndistroInfra: "true",
			},
			"tolerations": []map[string]interface{}{
				{
					"effect": "NoSchedule",
					"key":    "dedicated",
					"value":  "infra",
				},
			},
		}
		values = util.MergeMaps(values, infraValues)
	}

	release := appv1alpha1.HelmRelease{}
	key := client.ObjectKey{
		Name:      hr.GetObjectName(kyvernoReleaseName, p.Spec.ClusterName),
		Namespace: p.GetNamespace(),
	}
	err = r.Get(ctx, key, &release)
	if err != nil {
		if apierrors.IsNotFound(err) {
			release, err = hr.Prepare(kyvernoReleaseName, kyvernoReleaseName, p.GetNamespace(), kyvernoReleaseVersion, p.Spec.ClusterName, values)
			if err != nil {
				return appv1alpha1.DefaultPoliciesNotReady(p, meta.ObjectsApliedFailedReason, err.Error()), ctrl.Result{}, err
			}

			if release.Labels == nil {
				release.Labels = make(map[string]string)
			}
			release.Labels[meta.LabelUndistroMove] = ""

			err = hr.Install(ctx, r.Client, log, release, cl)
			if err != nil {
				return appv1alpha1.DefaultPoliciesNotReady(p, meta.ObjectsApliedFailedReason, err.Error()), ctrl.Result{}, err
			}
			return p, ctrl.Result{RequeueAfter: time.Second * 2}, nil
		}
		return p, ctrl.Result{}, err
	}

	meta.SetResourceCondition(&p, meta.ObjectsAppliedCondition, metav1.ConditionTrue, meta.ObjectsAppliedSuccessReason, "HelmRelease applied")
	if !meta.InReadyCondition(release.Status.Conditions) {
		return appv1alpha1.DefaultPoliciesNotReady(p, meta.WaitProvisionReason, "Wait Kyverno to be installed"), ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	clusterClient := r.Client
	if !util.IsMgmtCluster(p.Spec.ClusterName) {
		clusterClient, err = kube.NewClusterClient(ctx, r.Client, p.Spec.ClusterName, p.GetNamespace())
		if err != nil {
			return appv1alpha1.DefaultPoliciesNotReady(p, meta.GetClusterFailed, err.Error()), ctrl.Result{}, err
		}
	}

	p, err = r.applyPolicies(ctx, clusterClient, p)
	if err != nil {
		appv1alpha1.DefaultPoliciesNotReady(p, meta.ArtifactFailedReason, err.Error())
	}
	return appv1alpha1.DefaultPoliciesReady(p), ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *DefaultPoliciesReconciler) reconcileDelete(ctx context.Context, instance appv1alpha1.DefaultPolicies) (ctrl.Result, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}
	return hr.Uninstall(ctx, r.Client, log, kyvernoReleaseName, instance.Spec.ClusterName, instance.GetNamespace())
}

func (r *DefaultPoliciesReconciler) applyPolicies(ctx context.Context, clusterClient client.Client, p appv1alpha1.DefaultPolicies) (appv1alpha1.DefaultPolicies, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}
	dir, err := fs.PoliciesFS.ReadDir("policies")
	if err != nil {
		return p, err
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		log.Info("Installing policy", "name", f.Name())
		byt, err := fs.PoliciesFS.ReadFile(filepath.Join("policies", f.Name()))
		if err != nil {
			return p, err
		}

		objs, err := util.ToUnstructured(byt)
		if err != nil {
			return p, err
		}
		for _, o := range objs {
			if util.ContainsStringInSlice(p.Spec.ExcludePolicies, o.GetName()) {
				// delete policy if exists
				u := unstructured.Unstructured{}
				u.SetGroupVersionKind(o.GroupVersionKind())
				key := client.ObjectKey{
					Name: o.GetName(),
				}
				err = clusterClient.Get(ctx, key, &u)
				if !apierrors.IsNotFound(err) {
					err = clusterClient.Delete(ctx, &u)
					if err != nil {
						log.V(2).Error(err, "can't exclude policy", "name", u.GetName())
					}
				}
				continue
			}
			_, err = util.CreateOrUpdate(ctx, clusterClient, &o)
			if err != nil {
				log.Info("failed to apply policy", "name", o.GetName(), "err", err)
				return p, err
			}
			p.Status.AppliedPolicies = append(p.Status.AppliedPolicies, o.GetName())
			set := sets.NewString(p.Status.AppliedPolicies...)
			p.Status.AppliedPolicies = set.UnsortedList()
		}
	}
	return p, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DefaultPoliciesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.DefaultPolicies{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(r)
}
