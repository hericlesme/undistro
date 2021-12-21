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
package hr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"

	"github.com/getupio-undistro/meta"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/retry"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Install(ctx context.Context, c client.Client, log logr.Logger, hr appv1alpha1.HelmRelease, cl *appv1alpha1.Cluster) error {
	msg := fmt.Sprintf("Check condition for %s release", hr.Name)
	log.Info(msg, "lastAppliedVersion", hr.Status.LastAppliedRevision)
	if meta.InReadyCondition(hr.Status.Conditions) {
		meta.SetResourceCondition(cl, meta.ReadyCondition, metav1.ConditionTrue, meta.InstallSucceededReason, msg)
	}
	if !util.IsMgmtCluster(hr.Spec.ClusterName) {
		err := ctrl.SetControllerReference(cl, &hr, scheme.Scheme)
		if err != nil {
			return err
		}
	}
	err := retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
		_, e := util.CreateOrUpdate(ctx, c, &hr)
		return e
	})
	return err
}

func Prepare(releaseName, targetNs, clusterNs, version, clName string, v map[string]interface{}) (appv1alpha1.HelmRelease, error) {
	var values *apiextensionsv1.JSON
	if v != nil {
		byt, err := json.Marshal(v)
		if err != nil {
			return appv1alpha1.HelmRelease{}, err
		}
		values = &apiextensionsv1.JSON{
			Raw: byt,
		}
	}
	hrSpec := appv1alpha1.HelmReleaseSpec{
		ReleaseName:     releaseName,
		TargetNamespace: targetNs,
		Values:          values,
		Chart: appv1alpha1.ChartSource{
			RepoChartSource: appv1alpha1.RepoChartSource{
				RepoURL: undistro.DefaultRepo,
				Name:    releaseName,
				Version: version,
			},
		},
	}
	if !util.IsMgmtCluster(clName) {
		hrSpec.ClusterName = fmt.Sprintf("%s/%s", clusterNs, clName)
	}
	hr := &appv1alpha1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appv1alpha1.GroupVersion.String(),
			Kind:       "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetObjectName(releaseName, clName),
			Namespace: clusterNs,
		},
		Spec: hrSpec,
	}
	return *hr, nil
}

func Uninstall(ctx context.Context, c client.Client, log logr.Logger, releaseName, clusterName, ns string) (ctrl.Result, error) {
	// retrieve helmrelease
	release := appv1alpha1.HelmRelease{}
	key := client.ObjectKey{
		Name:      GetObjectName(releaseName, clusterName),
		Namespace: ns,
	}

	err := c.Get(ctx, key, &release)
	if err != nil {
		log.Info("error getting helm release", "error", err.Error())
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{Requeue: true}, err
		}
	}
	log.Info("deleting helm release", "releaseNameMeta", release.ObjectMeta.Name, "releaseName", release.Name)
	// delete helmrelease
	err = c.Delete(ctx, &release)
	if err != nil {
		log.Info("error deleting helm release", "error", err.Error())
		return ctrl.Result{Requeue: true}, err
	}
	// check if hr yet exists, if y requeue
	err = c.Get(ctx, key, &release)
	if err != nil {
		log.Info("error getting helm release", "error", err.Error())
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{Requeue: true}, err
		}
	}
	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(&release, meta.Finalizer)
	if err := c.Update(ctx, &release); err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

func GetObjectName(release, clName string) string {
	if clName == "" {
		clName = "management"
	}
	return fmt.Sprintf("%s-%s", release, clName)
}
