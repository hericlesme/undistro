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
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	observerRequeueAfter = time.Minute * 5
	monitoringNS         = "monitoring"
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
	instance := &appv1alpha1.Observer{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(instance, r.Client)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to init patch helper")
	}
	defer func() {
		var patchOpts []patch.Option
		if err == nil {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}
		patchErr := patchHelper.Patch(ctx, instance, patchOpts...)
		if patchErr != nil {
			err = kerrors.NewAggregate([]error{patchErr, err})
			r.Log.Info("failed to Patch identity")
		}
	}()

	r.Log.Info("Checking paused")
	if instance.Spec.Paused {
		r.Log.Info("Reconciliation is paused for this object")
		instance = appv1alpha1.ObserverPaused(*instance)
		return ctrl.Result{}, nil
	}
	r.Log.Info("Checking object age")
	if instance.Generation < instance.Status.ObservedGeneration {
		r.Log.Info("Skipping this old version of reconciled object")
		return ctrl.Result{}, nil
	}
	r.Log.Info("Reconciling Observer state...")
	if err := r.reconcile(ctx, *instance); err != nil {
		r.Log.Info(err.Error())
		return ctrl.Result{}, err
	}
	elapsed := time.Since(start)
	msg := fmt.Sprintf("Queueing after %s", elapsed.String())
	r.Log.Info(msg)
	return ctrl.Result{RequeueAfter: observerRequeueAfter}, nil
}

func (r *ObserverReconciler) reconcile(ctx context.Context, observer appv1alpha1.Observer) (err error) {
	err = r.reconcileMetrics(ctx, observer)
	if err != nil {
		return err
	}
	err = r.reconcileLog(ctx, observer)
	if err != nil {
		return err
	}
	return
}

func (r *ObserverReconciler) reconcileLog(ctx context.Context, observer appv1alpha1.Observer) (err error) {
	err = r.reconcileElasticStack(ctx, observer)
	if err != nil {
		return err
	}
	err = r.reconcileFluentStack(ctx, observer)
	if err != nil {
		return err
	}
	return nil
}

func (r *ObserverReconciler) reconcileFluentStack(ctx context.Context, observer appv1alpha1.Observer) (err error) {
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
			r.Log.Info(err.Error())
			return err
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

	// get elasticsearch service host and user info

	fluentdVersion := "0.2.11"
	if err := r.installRelease(ctx, "fluentd", fluentdVersion, values, &observer, cl); err != nil {
		return err
	}
	fluentBitVersion := "0.18.0"
	if err := r.installRelease(ctx, "fluent-bit", fluentBitVersion, values, &observer, cl); err != nil {
		return err
	}

	elasticFluentdPluginVersion := "12.0.0"
	if err := r.installRelease(ctx, "fluentd-elasticsearch", elasticFluentdPluginVersion, values, &observer, cl); err != nil {
		return err
	}
	return
}

func (r *ObserverReconciler) reconcileElasticStack(ctx context.Context, observer appv1alpha1.Observer) (err error) {
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
			r.Log.Info(err.Error())
			return err
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

	eckOperatorVersion := "1.9.0-SNAPSHOT"
	if err := r.installRelease(ctx, "eck-operator", eckOperatorVersion, values, &observer, cl); err != nil {
		return err
	}

	if err := r.createElasticsearchCluster(ctx); err != nil {
		return err
	}
	return nil
}

func (r *ObserverReconciler) createElasticsearchCluster(ctx context.Context) error {
	elasticsearch := `
---
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: elasticsearch
  namespace: undistro-system 
spec:
  version: 7.15.0
  nodeSets:
  - name: default
    count: 1
    config:
      node.store.allow_mmap: false
`
	objs, err := util.ToUnstructured([]byte(elasticsearch))
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

func (r *ObserverReconciler) reconcileMetrics(ctx context.Context, observer appv1alpha1.Observer) (err error) {
	// check kube-prometheus-stack installation from helm
	values := map[string]interface{}{
		"namespaceOverride": monitoringNS,
		"grafana": map[string]interface{}{
			"enabled": false,
		},
	}
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
			r.Log.Info(err.Error())
			return err
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

	kubeStackVersion := "18.0.4"
	if err := r.installRelease(ctx, "kube-prometheus-stack", kubeStackVersion, values, &observer, cl); err != nil {
		return err
	}
	// after the kube-prometheus-stack crds install
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		err = r.enableUnDistroMetrics(ctx)
		if err != nil {
			return err
		}
	}
	return
}

func (r *ObserverReconciler) installRelease(
	ctx context.Context, name, version string, values map[string]interface{}, observer *appv1alpha1.Observer, cl *appv1alpha1.Cluster) (err error) {
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
	release, err = hr.Prepare(name, monitoringNS, cl.GetNamespace(), version, observer.Spec.ClusterName, values)
	if err != nil {
		return err
	}
	if err := hr.Install(ctx, r.Client, release, cl); err != nil {
		return err
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

// SetupWithManager sets up the controller with the Manager.
func (r *ObserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Observer{}).
		Complete(r)
}
