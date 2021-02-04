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
package predicate

import (
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	configv1alpha1 "github.com/getupio-undistro/undistro/apis/config/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type ReconcileClusterChanges struct {
	predicate.Funcs
}

func (ReconcileClusterChanges) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	old, ok := e.ObjectOld.(*appv1alpha1.Cluster)
	if !ok {
		return false
	}
	n, ok := e.ObjectNew.(*appv1alpha1.Cluster)
	if !ok {
		return false
	}
	if !cmp.Equal(old.Spec, n.Spec) {
		return true
	}
	if !cmp.Equal(old.Labels, n.Labels) {
		return true
	}
	if !cmp.Equal(old.Annotations, n.Annotations) {
		return true
	}
	if !n.DeletionTimestamp.IsZero() {
		return true
	}

	if !meta.InReadyCondition(n.Status.Conditions) {
		return true
	}
	return false
}

type ReconcileHelmReleaseChanges struct {
	predicate.Funcs
}

func (ReconcileHelmReleaseChanges) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	old, ok := e.ObjectOld.(*appv1alpha1.HelmRelease)
	if !ok {
		return false
	}
	n, ok := e.ObjectNew.(*appv1alpha1.HelmRelease)
	if !ok {
		return false
	}
	if !cmp.Equal(old.Spec, n.Spec) {
		return true
	}
	if !cmp.Equal(old.Labels, n.Labels) {
		return true
	}
	if !cmp.Equal(old.Annotations, n.Annotations) {
		return true
	}
	if !n.DeletionTimestamp.IsZero() {
		return true
	}

	if !meta.InReadyCondition(n.Status.Conditions) {
		return true
	}
	return false
}

type ReconcileProviderChanges struct {
	predicate.Funcs
}

func (ReconcileProviderChanges) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	old, ok := e.ObjectOld.(*configv1alpha1.Provider)
	if !ok {
		return false
	}
	n, ok := e.ObjectNew.(*configv1alpha1.Provider)
	if !ok {
		return false
	}
	if !cmp.Equal(old.Spec, n.Spec) {
		return true
	}
	if !cmp.Equal(old.Labels, n.Labels) {
		return true
	}
	if !cmp.Equal(old.Annotations, n.Annotations) {
		return true
	}
	if !n.DeletionTimestamp.IsZero() {
		return true
	}

	if !meta.InReadyCondition(n.Status.Conditions) {
		return true
	}
	return false
}
