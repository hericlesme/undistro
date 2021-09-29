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
	"encoding/json"
	"fmt"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
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
	if err := r.reconcile(ctx, req, *instance); err != nil {
		r.Log.Info(err.Error())
		return ctrl.Result{}, err
	}
	elapsed := time.Since(start)
	msg := fmt.Sprintf("Queueing after %s", elapsed.String())
	r.Log.Info(msg)
	return ctrl.Result{RequeueAfter: observerRequeueAfter}, nil
}

func (r *ObserverReconciler) reconcile(ctx context.Context, req ctrl.Request, observer appv1alpha1.Observer) error {
	// check kube-prometheus-stack installation from helm
	values := make(map[string]interface{})
	defaultValues := map[string]interface{}{
		"namespaceOverride": undistro.Namespace,
		"grafana": map[string]interface{}{
			"enabled": false,
		},
	}
	cl := &appv1alpha1.Cluster{}
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		cl.Name = "management"
		cl.Namespace = undistro.Namespace
		values = defaultValues
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
								"node-role.undistro.io/infra": "true",
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
							"node-role.undistro.io/infra": "true",
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
						"node-role.undistro.io/infra": "true",
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
							"node-role.undistro.io/infra": "true",
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
						"node-role.undistro.io/infra": "true",
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
							"node-role.undistro.io/infra": "true",
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
			values = util.MergeMaps(defaultValues, infraValues)
		} else {
			values = defaultValues
		}
	}
	release := appv1alpha1.HelmRelease{}
	releaseName := "kube-prometheus-stack"
	msg := fmt.Sprintf("Checking if %s is installed", releaseName)
	r.Log.Info(msg)
	err := r.Get(ctx, req.NamespacedName, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		release, err = r.prepareRelease(releaseName, undistro.Namespace, cl.GetNamespace(), "18.0.4", observer, values)
		if err != nil {
			return err
		}
		msg = fmt.Sprintf("Installing %s", releaseName)
		r.Log.Info(msg)
		if err := installComponent(ctx, r.Client, release, cl); err != nil {
			return err
		}
	}
	// after the kube-prometheus-stack crds install
	if util.IsMgmtCluster(observer.Spec.ClusterName) {
		err = r.enableUnDistroMetrics(ctx, r.Client)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ObserverReconciler) enableUnDistroMetrics(ctx context.Context, c client.Client) error {
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
		_, err = util.CreateOrUpdate(ctx, c, &o)
		if err != nil {
			return err
		}
	}
	return nil
}

// prepareRelease fills the Helm Release fields for the given component
func (r *ObserverReconciler) prepareRelease(
	releaseName string,
	targetNs,
	clusterNs,
	version string,
	o appv1alpha1.Observer,
	v map[string]interface{}) (appv1alpha1.HelmRelease, error) {
	byt, err := json.Marshal(v)
	if err != nil {
		return appv1alpha1.HelmRelease{}, err
	}
	values := apiextensionsv1.JSON{
		Raw: byt,
	}
	hrSpec := appv1alpha1.HelmReleaseSpec{
		ReleaseName:     releaseName,
		TargetNamespace: targetNs,
		Values:          &values,
		Chart: appv1alpha1.ChartSource{
			RepoChartSource: appv1alpha1.RepoChartSource{
				RepoURL: undistro.DefaultRepo,
				Name:    releaseName,
				Version: version,
			},
		},
	}
	if !util.IsMgmtCluster(o.Spec.ClusterName) {
		hrSpec.ClusterName = fmt.Sprintf("%s/%s", clusterNs, o.Spec.ClusterName)
	}
	hr := &appv1alpha1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appv1alpha1.GroupVersion.String(),
			Kind:       "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaseName,
			Namespace: clusterNs,
		},
		Spec: hrSpec,
	}
	return *hr, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ObserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.Observer{}).
		Complete(r)
}
