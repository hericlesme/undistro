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
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/getupio-undistro/controllerlib"
	"github.com/getupio-undistro/meta"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/helm"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/getupio-undistro/undistro/pkg/version"
	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"helm.sh/helm/v3/pkg/strvals"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	getters = getter.Providers{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}
)

// HelmReleaseReconciler reconciles a HelmRelease object
type HelmReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	config *rest.Config
}

func (r *HelmReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	hr := appv1alpha1.HelmRelease{}
	if err := r.Get(ctx, req.NamespacedName, &hr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	keysAndValues := []interface{}{
		"helmrelease", req.NamespacedName,
		"chartRepo", hr.Spec.Chart.RepoURL,
		"chartName", hr.Spec.Chart.Name,
		"chartVersion", hr.Spec.Chart.Version,
	}
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}
	log.WithValues(keysAndValues...)

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(&hr, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer controllerlib.PatchInstance(ctx, controllerlib.InstanceOpts{
		Controller: "HelmReleaseController",
		Request:    req.String(),
		Object:     &hr,
		Error:      err,
		Helper:     patchHelper,
	})

	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(&hr, meta.Finalizer) {
		controllerutil.AddFinalizer(&hr, meta.Finalizer)
		return ctrl.Result{}, nil
	}
	if hr.Spec.Paused {
		log.Info("Reconciliation is paused for this object")
		hr = appv1alpha1.HelmReleasePaused(hr)
		return ctrl.Result{}, nil
	}
	if !hr.DeletionTimestamp.IsZero() {
		hr = appv1alpha1.HelmReleaseDeleting(hr)
		return r.reconcileDelete(ctx, hr)
	}
	if !util.IsMgmtCluster(hr.Spec.ClusterName) {
		key := util.ObjectKeyFromString(hr.Spec.ClusterName)
		cl := appv1alpha1.Cluster{}
		err = r.Get(ctx, key, &cl)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = ctrl.SetControllerReference(&cl, &hr, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}
		isCNIChart := false
		isKyverno := false
		if hr.Annotations != nil {
			_, isCNIChart = hr.Annotations[meta.SetupAnnotation]
			_, isKyverno = hr.Annotations[meta.KyvernoAnnotation]
		}
		if !meta.InReadyCondition(cl.Status.Conditions) && !isCNIChart {
			hr = appv1alpha1.HelmReleaseNotReady(hr, meta.WaitProvisionReason, "Wait cluster to be ready")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}

		policies := appv1alpha1.DefaultPoliciesList{}
		opts := make([]client.ListOption, 0)
		if cl.Namespace != "" {
			opts = append(opts, client.InNamespace(cl.Namespace))
		}
		err = r.List(ctx, &policies, opts...)
		if err != nil {
			return ctrl.Result{}, err
		}
		for _, p := range policies.Items {
			if p.Spec.ClusterName == cl.Name && !isCNIChart && !isKyverno {
				if !meta.InReadyCondition(p.Status.Conditions) {
					hr = appv1alpha1.HelmReleaseNotReady(hr, meta.WaitProvisionReason, "Wait cluster policies to be applied")
				}
			}
		}
	}
	hr, result, err := r.reconcile(ctx, hr)

	// Log reconciliation duration
	durationMsg := fmt.Sprintf("Reconcilation finished in %s", time.Since(start).String())
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	log.Info(durationMsg)
	return result, err
}

func (r *HelmReleaseReconciler) reconcile(ctx context.Context, hr appv1alpha1.HelmRelease) (appv1alpha1.HelmRelease, ctrl.Result, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	var clientOpts []getter.Option
	if hr.Spec.Chart.SecretRef != nil {
		resourceNamespacedName := types.NamespacedName{
			Name:      hr.Spec.Chart.SecretRef.Name,
			Namespace: hr.GetNamespace(),
		}
		var secret corev1.Secret
		err := r.Client.Get(ctx, resourceNamespacedName, &secret)
		if err != nil {
			err = fmt.Errorf("auth secret error: %w", err)
			hr = appv1alpha1.HelmReleaseNotReady(hr, meta.AuthenticationFailedReason, err.Error())
			return hr, ctrl.Result{}, err
		}
		opts, cleanup, err := helm.ClientOptionsFromSecret(secret)
		if err != nil {
			err = fmt.Errorf("auth options error: %w", err)
			hr = appv1alpha1.HelmReleaseNotReady(hr, meta.AuthenticationFailedReason, err.Error())
			return hr, ctrl.Result{}, err
		}
		defer cleanup()
		clientOpts = opts
	}
	if hr.Spec.Timeout != nil {
		clientOpts = append(clientOpts, getter.WithTimeout(hr.Spec.Timeout.Duration))
	}
	chartRepo, err := helm.NewChartRepository(hr.Spec.Chart.RepoURL, getters, clientOpts)
	if err != nil {
		switch err.(type) {
		default:
			hr = appv1alpha1.HelmReleaseNotReady(hr, meta.IndexationFailedReason, err.Error())
		case *url.Error:
			hr = appv1alpha1.HelmReleaseNotReady(hr, meta.URLInvalidReason, err.Error())
		}
		return hr, ctrl.Result{}, err
	}
	if err := chartRepo.DownloadIndex(); err != nil {
		err = fmt.Errorf("failed to download repository index: %w", err)
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.IndexationFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	chartRepo.Index.SortEntries()
	versions := chartRepo.Index.Entries[hr.Spec.Chart.Name]
	if versions.Len() > 0 {
		latestVersion := versions[0]
		lv, err := version.ParseVersion(latestVersion.Version)
		if err != nil {
			return hr, ctrl.Result{}, err
		}
		if hr.Spec.Chart.Version == "" {
			hr.Spec.Chart.Version = lv.String()
			return hr, ctrl.Result{}, nil
		}
		if hr.Spec.AutoUpgrade {
			acv, err := version.ParseVersion(hr.Spec.Chart.Version)
			if err != nil {
				return hr, ctrl.Result{}, err
			}
			if lv.GreaterThan(acv) && lv.Major() == acv.Major() {
				hr.Spec.Chart.Version = lv.String()
				return hr, ctrl.Result{}, nil
			}
		}
	}
	ch, err := chartRepo.Get(hr.Spec.Chart.Name, hr.Spec.Chart.Version)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.ChartPullFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	res, err := chartRepo.DownloadChart(ch)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.ChartPullFailedReason, err.Error())
		return hr, ctrl.Result{Requeue: true}, err
	}
	// Check dependencies
	if len(hr.Spec.Dependencies) > 0 {
		if err := r.checkDependencies(ctx, hr); err != nil {
			msg := fmt.Sprintf("dependencies do not meet ready condition (%s), retrying in 5s", err.Error())
			log.Info(msg)

			// Exponential backoff would cause execution to be prolonged too much,
			// instead we requeue on a fixed interval.
			return appv1alpha1.HelmReleaseNotReady(hr,
				meta.DependencyNotReadyReason, err.Error()), ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		log.Info("all dependencies are ready, proceeding with release")
	}
	restClientGetter, err := r.getRESTClientGetter(ctx, hr)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return hr, ctrl.Result{}, err
		}
		if _, isSetupChart := hr.Annotations[meta.SetupAnnotation]; isSetupChart {
			return appv1alpha1.HelmReleaseNotReady(hr, meta.WaitProvisionReason, "wait kubeconfig to be available"), ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		return hr, ctrl.Result{}, nil
	}
	restCfg, err := restClientGetter.ToRESTConfig()
	if err != nil {
		if _, isCNIChart := hr.Annotations[meta.SetupAnnotation]; isCNIChart {
			return appv1alpha1.HelmReleaseNotReady(hr, meta.WaitProvisionReason, "wait kubeconfig to be available"), ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		return hr, ctrl.Result{}, err
	}
	workloadClient, err := client.New(restCfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		if _, isCNIChart := hr.Annotations[meta.SetupAnnotation]; isCNIChart {
			return appv1alpha1.HelmReleaseNotReady(hr, meta.WaitProvisionReason, "wait kubeconfig to be available"), ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		return hr, ctrl.Result{}, err
	}
	_, err = util.CreateOrUpdate(ctx, workloadClient, &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: hr.GetNamespace(),
		},
	})
	if err != nil {
		return hr, ctrl.Result{}, err
	}
	// install helm secrets in undistro-system so this namespace need to exist in workload clusters
	_, err = util.CreateOrUpdate(ctx, workloadClient, &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "undistro-system",
		},
	})
	if err != nil {
		return hr, ctrl.Result{}, err
	}
	err = r.applyObjs(ctx, workloadClient, hr.Spec.BeforeApplyObjects)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.ObjectsApliedFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	meta.SetResourceCondition(&hr, meta.ObjectsAppliedCondition, metav1.ConditionTrue, meta.ObjectsAppliedSuccessReason, "objects successfully applied before install")
	// Compose values
	values, err := r.composeValues(ctx, hr)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.InitFailedReason, err.Error())
		return hr, ctrl.Result{Requeue: true}, nil
	}
	hc, err := loader.LoadArchive(res)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.StorageOperationFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	hr, err = r.reconcileRelease(ctx, restClientGetter, workloadClient, hr, hc, values)
	if err != nil {
		if errors.Is(err, driver.ErrNoDeployedReleases) {
			return hr, ctrl.Result{Requeue: true}, nil
		}
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.ReconciliationFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	err = r.applyObjs(ctx, workloadClient, hr.Spec.AfterApplyObjects)
	if err != nil {
		hr = appv1alpha1.HelmReleaseNotReady(hr, meta.ObjectsApliedFailedReason, err.Error())
		return hr, ctrl.Result{}, err
	}
	meta.SetResourceCondition(&hr, meta.ObjectsAppliedCondition, metav1.ConditionTrue, meta.ObjectsAppliedSuccessReason, "objects successfully applied after install")
	return hr, ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *HelmReleaseReconciler) applyObjs(ctx context.Context, c client.Client, objs []apiextensionsv1.JSON) error {
	for _, raw := range objs {
		uobjs, err := util.ToUnstructured(raw.Raw)
		if err != nil {
			return err
		}
		for _, o := range uobjs {
			_, err = util.CreateOrUpdate(ctx, c, &o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *HelmReleaseReconciler) reconcileRelease(ctx context.Context, restClientGetter genericclioptions.RESTClientGetter, workloadClient client.Client,
	hr appv1alpha1.HelmRelease, chart *chart.Chart, values chartutil.Values) (appv1alpha1.HelmRelease, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	log.Info("Reconciling release", "release name", chart.Name())
	// Initialize Helm action runner
	runner, err := helm.NewRunner(restClientGetter, hr.Spec.TargetNamespace, log)
	if err != nil {
		return appv1alpha1.HelmReleaseNotReady(hr, meta.InitFailedReason, "failed to initialize Helm action runner"), err
	}
	rel, err := runner.ObserveLastRelease(hr)
	if err != nil {
		return appv1alpha1.HelmReleaseNotReady(hr, meta.GetLastReleaseFailedReason, "failed to get last release revision"), err
	}
	revision := chart.Metadata.Version
	releaseRevision := util.ReleaseRevision(rel)
	valuesChecksum := util.ValuesChecksum(values)
	hr, hasNewState := appv1alpha1.HelmReleaseAttempted(hr, revision, releaseRevision, valuesChecksum)
	if hasNewState {
		hr = appv1alpha1.HelmReleaseProgressing(hr)
	}
	if hr.Labels[meta.LabelProviderType] == "core" && rel != nil {
		if rel.Chart.Metadata.Version == hr.Spec.Chart.Version && rel.Info.Status == release.StatusDeployed {
			return appv1alpha1.HelmReleaseReady(hr), nil
		}
	}
	// Check status of any previous release attempt.
	if meta.InReadyCondition(hr.Status.Conditions) && !hasNewState && rel != nil && rel.Info.Deleted.IsZero() {
		return appv1alpha1.HelmReleaseReady(hr), nil
	}
	var isInstallation bool
	var installErr error
	if rel == nil || rel.Version == 0 {
		isInstallation = true
		rel, installErr = runner.Install(hr, chart, values)
		installErr = r.handleHelmActionResult(ctx, &hr, revision, installErr, "install", meta.ReleasedCondition, meta.InstallSucceededReason, meta.InstallFailedReason)
	} else if (rel.Info.Status == release.StatusDeployed && hasNewState) || rel.Info.Status == release.StatusUninstalled || rel.Info.Status == release.StatusPendingUpgrade {
		if rel.Info.Status == release.StatusPendingUpgrade {
			err := runner.UpdateState(rel)
			if err != nil {
				return hr, err
			}
			return hr, nil
		}
		if rel.Info.Status == release.StatusUninstalled {
			hr.Spec.ForceUpgrade = pointer.Bool(true)
		}
		rel, err = runner.Upgrade(hr, chart, values)
		err = r.handleHelmActionResult(ctx, &hr, revision, err, "upgrade", meta.ReleasedCondition, meta.UpgradeSucceededReason, meta.UpgradeFailedReason)
	}
	if util.ReleaseRevision(rel) > releaseRevision {
		if err == nil && hr.Spec.Test.Enable {
			_, err = runner.Test(hr)
			err = r.handleHelmActionResult(ctx, &hr, revision, err, "test", meta.TestSuccessCondition, meta.TestSucceededReason, meta.TestFailedReason)
			if err != nil && hr.Spec.Test.IgnoreFailures {
				err = nil
			}
		}
	}
	rolledBack := false
	var rollErr error
	if err != nil {
		if errors.Is(err, driver.ErrNoDeployedReleases) {
			slist := corev1.SecretList{}
			serr := workloadClient.List(ctx, &slist, client.InNamespace(hr.Spec.TargetNamespace))
			if serr != nil {
				return hr, err
			}
			for _, i := range slist.Items {
				if strings.Contains(i.Name, hr.Spec.ReleaseName) {
					serr = workloadClient.Delete(ctx, &i)
					if serr != nil {
						return hr, err
					}
				}
			}
			return appv1alpha1.ResetHelmReleaseStatus(hr), err
		}
		if util.ReleaseRevision(rel) <= releaseRevision {
			log.Info("skip, no new revision created")
		} else if !isInstallation {
			rolledBack = true
			rollErr = err
			err = runner.Rollback(hr)
			err = r.handleHelmActionResult(ctx, &hr, revision, err, "rollback", meta.RemediatedCondition, meta.RollbackSucceededReason, meta.RollbackSucceededReason)
		}
	}
	if !isInstallation {
		rel, observeLastReleaseErr := runner.ObserveLastRelease(hr)
		if observeLastReleaseErr != nil {
			err = observeLastReleaseErr
		}
		hr.Status.LastReleaseRevision = util.ReleaseRevision(rel)
		if err != nil {
			return appv1alpha1.HelmReleaseNotReady(hr, meta.ReconciliationFailedReason, err.Error()), err
		}
		hr.Spec.Chart.RepoChartSource.Version = rel.Chart.Metadata.Version
		if rolledBack {
			return appv1alpha1.HelmReleaseNotReady(hr, meta.RollbackSucceededReason, rollErr.Error()), nil
		}
	}
	if installErr != nil {
		return appv1alpha1.HelmReleaseNotReady(hr, meta.InstallFailedReason, installErr.Error()), nil
	}
	return appv1alpha1.HelmReleaseReady(hr), nil
}

func (r *HelmReleaseReconciler) checkDependencies(ctx context.Context, hr appv1alpha1.HelmRelease) error {
	for _, d := range hr.Spec.Dependencies {
		if d.Namespace == "" {
			d.Namespace = hr.GetNamespace()
		}
		var dHr appv1alpha1.HelmRelease
		nm := types.NamespacedName{
			Name:      d.Name,
			Namespace: d.Namespace,
		}
		err := r.Get(ctx, nm, &dHr)
		if err != nil {
			return fmt.Errorf("unable to get '%v' dependency: %w", nm, err)
		}

		if len(dHr.Status.Conditions) == 0 || dHr.Generation != dHr.Status.ObservedGeneration {
			return fmt.Errorf("dependency '%v' is not ready", nm)
		}

		if !apimeta.IsStatusConditionTrue(dHr.Status.Conditions, meta.ReadyCondition) {
			return fmt.Errorf("dependency '%v' is not ready", nm)
		}
	}
	return nil
}

func (r *HelmReleaseReconciler) handleHelmActionResult(ctx context.Context, hr *appv1alpha1.HelmRelease, revision string, err error, action string, condition string, succeededReason string, failedReason string) error {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	log.Info("Release revision", "revision", revision)
	if err != nil {
		msg := fmt.Sprintf("Helm %s failed: %s", action, err.Error())
		meta.SetResourceCondition(hr, condition, metav1.ConditionFalse, failedReason, msg)
		return err
	} else {
		msg := fmt.Sprintf("Helm %s succeeded", action)
		meta.SetResourceCondition(hr, condition, metav1.ConditionTrue, succeededReason, msg)
		return nil
	}
}

// composeValues attempts to resolve all \ValuesReference resources
// and merges them as defined. Referenced resources are only retrieved once
// to ensure a single version is taken into account during the merge.
func (r *HelmReleaseReconciler) composeValues(ctx context.Context, hr appv1alpha1.HelmRelease) (chartutil.Values, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	result := chartutil.Values{}
	configMaps := make(map[string]*corev1.ConfigMap)
	secrets := make(map[string]*corev1.Secret)

	for _, v := range hr.Spec.ValuesFrom {
		namespacedName := types.NamespacedName{Namespace: hr.Namespace, Name: v.Name}
		var valuesData []byte

		switch v.Kind {
		case "ConfigMap":
			resource, ok := configMaps[namespacedName.String()]
			if !ok {
				// The resource may not exist, but we want to act on a single version
				// of the resource in case the values reference is marked as optional.
				configMaps[namespacedName.String()] = nil

				resource = &corev1.ConfigMap{}
				if err := r.Get(ctx, namespacedName, resource); err != nil {
					if apierrors.IsNotFound(err) {
						if v.Optional {
							log.Info("could not find optional %s '%s'", v.Kind, namespacedName)
							continue
						}
						return nil, fmt.Errorf("could not find %s '%s'", v.Kind, namespacedName)
					}
					return nil, err
				}
				configMaps[namespacedName.String()] = resource
			}
			if resource == nil {
				if v.Optional {
					log.Info("could not find optional %s '%s'", v.Kind, namespacedName)
					continue
				}
				return nil, fmt.Errorf("could not find %s '%s'", v.Kind, namespacedName)
			}
			if data, ok := resource.Data[v.ValuesKey]; !ok {
				return nil, fmt.Errorf("missing key '%s' in %s '%s'", v.ValuesKey, v.Kind, namespacedName)
			} else {
				valuesData = []byte(data)
			}
		case "Secret":
			resource, ok := secrets[namespacedName.String()]
			if !ok {
				// The resource may not exist, but we want to act on a single version
				// of the resource in case the values reference is marked as optional.
				secrets[namespacedName.String()] = nil

				resource = &corev1.Secret{}
				if err := r.Get(ctx, namespacedName, resource); err != nil {
					if apierrors.IsNotFound(err) {
						if v.Optional {
							log.Info("could not find optional %s '%s'", v.Kind, namespacedName)
							continue
						}
						return nil, fmt.Errorf("could not find %s '%s'", v.Kind, namespacedName)
					}
					return nil, err
				}
				secrets[namespacedName.String()] = resource
			}
			if resource == nil {
				if v.Optional {
					log.Info("could not find optional %s '%s'", v.Kind, namespacedName)
					continue
				}
				return nil, fmt.Errorf("could not find %s '%s'", v.Kind, namespacedName)
			}
			if data, ok := resource.Data[v.ValuesKey]; !ok {
				return nil, fmt.Errorf("missing key '%s' in %s '%s'", v.ValuesKey, v.Kind, namespacedName)
			} else {
				valuesData = data
			}
		default:
			return nil, fmt.Errorf("unsupported ValuesReference kind '%s'", v.Kind)
		}
		switch v.TargetPath {
		case "":
			values, err := chartutil.ReadValues(valuesData)
			if err != nil {
				return nil, fmt.Errorf("unable to read values from key '%s' in %s '%s': %w", v.ValuesKey, v.Kind, namespacedName, err)
			}
			result = util.MergeMaps(result, values)
		default:
			// TODO(hidde): this is a bit of hack, as it mimics the way the option string is passed
			// 	to Helm from a CLI perspective. Given the parser is however not publicly accessible
			// 	while it contains all logic around parsing the target path, it is a fair trade-off.
			singleValue := v.TargetPath + "=" + string(valuesData)
			if err := strvals.ParseInto(singleValue, result); err != nil {
				return nil, fmt.Errorf("unable to merge value from key '%s' in %s '%s' into target path '%s': %w", v.ValuesKey, v.Kind, namespacedName, v.TargetPath, err)
			}
		}
	}
	m := map[string]interface{}{}
	if hr.Spec.Values != nil {
		json.Unmarshal(hr.Spec.Values.Raw, &m)
	}
	return util.MergeMaps(result, m), nil
}

func (r *HelmReleaseReconciler) getRESTClientGetter(ctx context.Context, hr appv1alpha1.HelmRelease) (genericclioptions.RESTClientGetter, error) {
	if hr.Annotations == nil {
		hr.Annotations = make(map[string]string)
	}

	_, localChart := hr.Annotations[meta.HelmReleaseLocation]
	if util.IsMgmtCluster(hr.Spec.ClusterName) || localChart {
		return kube.NewInClusterRESTClientGetter(r.config, hr.Spec.TargetNamespace), nil
	}

	key := util.ObjectKeyFromString(hr.Spec.ClusterName)
	kubeConfig, err := kubeconfig.FromSecret(ctx, r.Client, key)
	if err != nil {
		return nil, err
	}
	return kube.NewMemoryRESTClientGetter(kubeConfig, hr.Spec.TargetNamespace), nil
}

func (r *HelmReleaseReconciler) reconcileDelete(ctx context.Context, hr appv1alpha1.HelmRelease) (ctrl.Result, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	restClient, err := r.getRESTClientGetter(ctx, hr)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	runner, err := helm.NewRunner(restClient, hr.Spec.TargetNamespace, log)
	if err != nil {
		return ctrl.Result{}, err
	}

	rel, err := runner.ObserveLastRelease(hr)
	if err != nil {
		controllerutil.RemoveFinalizer(&hr, meta.Finalizer)
		_, err = util.CreateOrUpdate(ctx, r.Client, &hr)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if rel == nil {
		controllerutil.RemoveFinalizer(&hr, meta.Finalizer)
		_, err = util.CreateOrUpdate(ctx, r.Client, &hr)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	err = runner.Uninstall(hr)
	if err != nil {
		return ctrl.Result{}, err
	}
	l, err := runner.List()
	if err != nil {
		return ctrl.Result{}, err
	}
	for _, i := range l {
		if i.Name == hr.Spec.ReleaseName {
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(&hr, meta.Finalizer)
	_, err = util.CreateOrUpdate(ctx, r.Client, &hr)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.config = mgr.GetConfig()
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.HelmRelease{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
