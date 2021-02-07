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
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type ClusterChanges struct {
	predicate.Funcs
}

func (ClusterChanges) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	cOld, ok := e.ObjectOld.(*appv1alpha1.Cluster)
	if !ok {
		return false
	}
	cn, ok := e.ObjectNew.(*appv1alpha1.Cluster)
	if !ok {
		return false
	}
	old := cOld.DeepCopy()
	n := cn.DeepCopy()
	if meta.InReadyCondition(old.Status.Conditions) && meta.InReadyCondition(n.Status.Conditions) &&
		cmp.Equal(old.Spec, n.Spec) && cmp.Equal(old.Status, n.Status) && n.DeletionTimestamp.IsZero() {
		return false
	}
	return true
}

type ReconcileChanges struct {
	predicate.Funcs
}

func (ReconcileChanges) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	if e.ObjectOld.GetGeneration() == e.ObjectNew.GetGeneration() {
		return false
	}
	if e.ObjectOld.GetResourceVersion() == e.ObjectNew.GetResourceVersion() {
		return false
	}
	return true
}
