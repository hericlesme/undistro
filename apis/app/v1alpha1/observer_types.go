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

package v1alpha1

import (
	"github.com/getupio-undistro/undistro/pkg/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ObserverSpec defines the desired state of Observer
type ObserverSpec struct {
	// Pause Observer reconciliation.
	Paused bool `json:"paused,omitempty"`

	// ClusterName is the name of the cluster to which this Observer belongs.
	ClusterName string `json:"clusterName,omitempty"`
}

// ObserverStatus defines the observed state of Observer
type ObserverStatus struct {
	// ObservedGeneration is the last observed generation.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Observer is the Schema for the observers API
type Observer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObserverSpec   `json:"spec,omitempty"`
	Status ObserverStatus `json:"status,omitempty"`
}

// IdentityPaused registers a paused reconciliation of the given Cluster.
func ObserverPaused(i Observer) *Observer {
	meta.SetResourceCondition(&i, meta.ReadyCondition, metav1.ConditionTrue, meta.ReconciliationPausedReason, meta.ReconciliationPausedReason)
	return &i
}

func (i *Observer) GetStatusConditions() *[]metav1.Condition {
	return &i.Status.Conditions
}

//+kubebuilder:object:root=true

// ObserverList contains a list of Observer
type ObserverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Observer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Observer{}, &ObserverList{})
}
