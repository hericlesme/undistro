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
	"strings"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	observerRequeueAfter   = time.Minute * 5
	monitoringNs           = "monitoring"
	kubeStackVersion       = "18.0.4"
	kubeStackReleaseName   = "kube-prometheus-stack"
	eckOperatorVersion     = "1.8.0"
	eckOperatorReleaseName = "eck-operator"
	fluentBitVersion       = "0.18.0"
	fluentBitReleaseName   = "fluent-bit"
	fluentdVersion         = "0.2.12"
	fluentdReleaseName     = "fluentd"
)

// ObserverReconciler reconciles a Observer object
type ObserverReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=app.undistro.io,resources=observers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.undistro.io,resources=observers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.undistro.io,resources=observers/finalizers,verbs=update

func (r *ObserverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	r.Log.Info("Reconciling Observer state", "request", req.String())
	instance := &appv1alpha1.Observer{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		} else {
			return ctrl.Result{}, err
		}
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(instance, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer patchInstance(Instance{
		Ctx:        ctx,
		Log:        r.Log,
		Controller: "DefaultPoliciesController",
		Request:    req.String(),
		Object:     instance,
		Error:      err,
		Helper:     patchHelper,
	})

	r.Log.Info("Checking paused")
	if instance.Spec.Paused {
		r.Log.Info("Reconciliation is paused for this object")
		// nolint
		instance = appv1alpha1.ObserverPaused(*instance)
		return ctrl.Result{}, nil
	}
	r.Log.Info("Checking object age")
	if instance.Generation < instance.Status.ObservedGeneration {
		r.Log.Info("Skipping this old version of reconciled object")
		return ctrl.Result{}, nil
	}
	r.Log.Info("Checking if object is under deletion")
	if instance.DeletionTimestamp.IsZero() {
		// Add our finalizer if it does not exist
		if !controllerutil.ContainsFinalizer(instance, meta.Finalizer) {
			r.Log.Info("Adding finalizer", "finalizer", meta.Finalizer)
			controllerutil.AddFinalizer(instance, meta.Finalizer)
			return ctrl.Result{}, nil
		}
	} else {
		return r.reconcileDelete(ctx, instance)
	}
	result, err := r.reconcile(ctx, *instance)
	durationMsg := fmt.Sprintf("Reconcilation finished in %s", time.Since(start).String())
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	r.Log.Info(durationMsg)
	return result, err
}

func (r *ObserverReconciler) reconcile(ctx context.Context, observer appv1alpha1.Observer) (ctrl.Result, error) {
	res, err := r.reconcileMetrics(ctx, observer)
	if err != nil {
		return ctrl.Result{}, err
	}
	res, err = r.reconcileLog(ctx, observer)
	if err != nil {
		return ctrl.Result{}, err
	}
	res.RequeueAfter = observerRequeueAfter
	return res, err
}

func (r *ObserverReconciler) reconcileLog(ctx context.Context, observer appv1alpha1.Observer) (ctrl.Result, error) {
	r.Log.Info("Reconciling log stack", "observerName", observer.Name)
	res, err := r.reconcileElasticStack(ctx, observer)
	if err != nil {
		r.Log.Info("error reconciling elastic eck operator", "error", err.Error())
		return res, err
	}
	return r.reconcileFluentStack(ctx, observer)
}

func (r *ObserverReconciler) reconcileFluentStack(ctx context.Context, observer appv1alpha1.Observer) (ctrl.Result, error) {
	r.Log.Info("Reconciling fluent stack", "observerName", observer.Name)
	values := map[string]interface{}{}
	cl := &appv1alpha1.Cluster{}
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		cl.Name = "management"
		cl.Namespace = undistro.Namespace
	} else {
		key := client.ObjectKey{
			Name:      observer.Spec.ClusterName,
			Namespace: observer.GetNamespace(),
		}
		err := r.Get(ctx, key, cl)
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		if cl.HasInfraNodes() {
			values = map[string]interface{}{
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
		}
	}
	if err := r.installRelease(ctx, fluentBitReleaseName, fluentBitVersion, values, &observer, cl); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.installRelease(ctx, fluentdReleaseName, fluentdVersion, values, &observer, cl); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ObserverReconciler) reconcileElasticStack(ctx context.Context, observer appv1alpha1.Observer) (ctrl.Result, error) {
	r.Log.Info("Reconciling elastic stack", "observerName", observer.Name)
	values := map[string]interface{}{}
	cl := &appv1alpha1.Cluster{}
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		cl.Name = "management"
		cl.Namespace = undistro.Namespace
	} else {
		key := client.ObjectKey{
			Name:      observer.Spec.ClusterName,
			Namespace: observer.GetNamespace(),
		}
		err := r.Get(ctx, key, cl)
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{Requeue: true}, err
		}
		if cl.HasInfraNodes() {
			values = map[string]interface{}{
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
		}
	}
	if err := r.installRelease(ctx, eckOperatorReleaseName, eckOperatorVersion, values, &observer, cl); err != nil {
		return ctrl.Result{}, err
	}
	return r.createElasticsearchCluster(ctx, &observer, cl)
}

func (r *ObserverReconciler) createElasticsearchCluster(ctx context.Context, obs *appv1alpha1.Observer, cl *appv1alpha1.Cluster) (ctrl.Result, error) {
	replicaCount := 0
	if cl.HasInfraNodes() {
		for _, w := range cl.Spec.Workers {
			if w.InfraNode {
				replicaCount = int(*w.Replicas)
				r.Log.Info("Replicas incremented", "replicaCount", replicaCount)
				break
			}
		}
	} else {
		replicaCount = int(cl.Status.TotalWorkerReplicas)
	}
	elasticsearch := fmt.Sprintf(`
--- 
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata: 
  name: elasticsearch
  namespace: monitoring
spec: 
  nodeSets: 
    - 
      config: 
        node.attr.attr_name: attr_value
        node.roles: 
          - master
          - data
        node.store.allow_mmap: false
      count: 1
      name: master
      podTemplate: 
        spec: 
          containers: 
            - 
              name: elasticsearch
    - 
      config: 
        node.data: true
      count: %d
      name: worker
  version: "7.15.1"
`, replicaCount)
	objs, err := util.ToUnstructured([]byte(elasticsearch))
	if err != nil {
		return ctrl.Result{}, err
	}
	u := unstructured.Unstructured{}
	k := client.ObjectKey{}
	clusterClient := r.Client
	if !util.IsMgmtCluster(obs.Spec.ClusterName) {
		clusterClient, err = kube.NewClusterClient(ctx, r.Client, obs.Spec.ClusterName, obs.GetNamespace())
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	for _, o := range objs {
		_, err = util.CreateOrUpdate(ctx, clusterClient, &o)
		if err != nil {
			return ctrl.Result{}, err
		}
		u.SetGroupVersionKind(o.GroupVersionKind())
		k.Namespace = o.GetNamespace()
		k.Name = o.GetName()
	}
	err = clusterClient.Get(ctx, k, &u)
	if err != nil {
		return ctrl.Result{}, err
	}
	health, ok, err := unstructured.NestedString(u.Object, "status", "health")
	if err != nil {
		return ctrl.Result{}, err
	}
	if !ok {
		return ctrl.Result{}, errors.New("health field not present")
	}
	desiredHealth := "green"
	// maybe this will need refactoring
	for strings.ToLower(health) != desiredHealth {
		r.Log.Info("Waiting elasticsearch", "health", health, "desiredHealth", desiredHealth)
		err = clusterClient.Get(ctx, k, &u)
		if err != nil {
			return ctrl.Result{}, err
		}
		health, ok, err = unstructured.NestedString(u.Object, "status", "health")
		if err != nil {
			return ctrl.Result{}, err
		}
		if !ok {
			return ctrl.Result{}, errors.New("health field not present")
		}
		<-time.After(time.Minute * 1)
	}
	return ctrl.Result{}, nil
}

func (r *ObserverReconciler) reconcileMetrics(ctx context.Context, observer appv1alpha1.Observer) (ctrl.Result, error) {
	r.Log.Info("Reconciling metrics stack", "observerName", observer.Name)
	values := map[string]interface{}{
		"namespaceOverride": monitoringNs,
		"grafana": map[string]interface{}{
			"enabled": false,
		},
	}
	cl := &appv1alpha1.Cluster{}
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		cl.Name = "management"
		cl.Namespace = undistro.Namespace
		r.Log.Info("Cluster is a management cluster", "name", cl.Name, "namespace", undistro.Namespace)
	} else {
		r.Log.Info("Cluster is a managed cluster", "name", cl.Name, "namespace", undistro.Namespace)
		key := client.ObjectKey{
			Name:      observer.Spec.ClusterName,
			Namespace: observer.GetNamespace(),
		}
		err := r.Get(ctx, key, cl)
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		if cl.HasInfraNodes() {
			infraValues := map[string]interface{}{
				"prometheusOperator": map[string]interface{}{
					"admissionWebhooks": map[string]interface{}{
						"patch": map[string]interface{}{
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
						},
					},
					"prometheusSpec": map[string]interface{}{
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
					},
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
				},
				"alertmanager": map[string]interface{}{
					"alertmanagerSpec": map[string]interface{}{
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
					},
				},
				"kube-state-metrics": map[string]interface{}{
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
				},
				"prometheus": map[string]interface{}{
					"prometheusSpec": map[string]interface{}{
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
					},
				},
			}
			values = util.MergeMaps(values, infraValues)
		}
	}
	if err := r.installRelease(ctx, kubeStackReleaseName, kubeStackVersion, values, &observer, cl); err != nil {
		return ctrl.Result{}, err
	}
	// after the kube-prometheus-stack crds install
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		err := r.enableUnDistroMetrics(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *ObserverReconciler) installRelease(
	ctx context.Context, name, version string, values map[string]interface{}, observer *appv1alpha1.Observer, cl *appv1alpha1.Cluster) (err error) {
	r.Log.Info("Installing release", "releaseName", name, "observerName", observer.Name)
	key := client.ObjectKey{
		Name:      hr.GetObjectName(name, observer.Spec.ClusterName),
		Namespace: observer.GetNamespace(),
	}
	release := appv1alpha1.HelmRelease{}
	err = r.Get(ctx, key, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}
	release, err = hr.Prepare(name, monitoringNs, cl.GetNamespace(), version, observer.Spec.ClusterName, values)
	if err != nil {
		return err
	}
	if release.Labels == nil {
		release.Labels = make(map[string]string)
	}
	release.Labels[meta.LabelUndistroMove] = ""
	if err := hr.Install(ctx, r.Client, r.Log, release, cl); err != nil {
		return err
	}
	for !meta.InReadyCondition(release.Status.Conditions) {
		err = r.Get(ctx, key, &release)
		if err != nil {
			return err
		}
		if meta.InReadyCondition(release.Status.Conditions) {
			continue
		} else {
			r.Log.Info("Waiting release is ready", "release", release.Name, "namespace", release.Namespace)
			<-time.After(1 * time.Minute)
		}
		r.Log.Info("HelmRelease status", "lastAppliedVersion", release.Status.LastAppliedRevision)
	}
	return
}

func (r *ObserverReconciler) enableUnDistroMetrics(ctx context.Context) error {
	undistroServiceMonitor := `
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    control-plane: controller-manager
    undistro.io: undistro
  name: undistro-controller-manager-metrics-monitor
  namespace: undistro-system
spec:
  endpoints:
  - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    path: /metrics
    port: https
    scheme: https
    tlsConfig:
      insecureSkipVerify: true
  selector:
    matchLabels:
      control-plane: controller-manager
`
	objs, err := util.ToUnstructured([]byte(undistroServiceMonitor))
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

func (r *ObserverReconciler) reconcileDelete(ctx context.Context, instance *appv1alpha1.Observer) (res ctrl.Result, err error) {
	r.Log.Info("Reconciling delete", "clusterName", instance.Spec.ClusterName)
	releases := []string{kubeStackReleaseName, eckOperatorReleaseName, fluentBitVersion, fluentdReleaseName}
	for _, release := range releases {
		r.Log.Info("Deleting charts", "release", release, "namespace", instance.GetNamespace())
		res, err = hr.Uninstall(ctx, r.Client, r.Log, release, instance.Spec.ClusterName, instance.GetNamespace())
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		For(&appv1alpha1.Observer{}).
		Complete(r)
}
