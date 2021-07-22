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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProviderInfo struct {
	SecretName string
	Category   ProviderCategory
}

// +kubebuilder:validation:Enum=core;infra
type ProviderCategory string

const (
	ProviderCore  = ProviderCategory("core")
	ProviderInfra = ProviderCategory("infra")
)

// ProviderSpec defines the desired state of Provider
type ProviderSpec struct {
	Paused          bool                    `json:"paused,omitempty"`
	AutoFetch       bool                    `json:"autoFetch,omitempty"`
	UnDistroVersion string                  `json:"undistroVersion,omitempty"`
	Category        ProviderCategory        `json:"category,omitempty"`
	SecretRef       *corev1.ObjectReference `json:"secretRef,omitempty"`
}

// ProviderStatus defines the observed state of Provider
type ProviderStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions   []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	ChartName    string             `json:"chartName,omitempty"`
	ChartVersion string             `json:"chartVersion,omitempty"`
	RegionNames  []string           `json:"regionNames,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Provider is the Schema for the providers API
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec   `json:"spec,omitempty"`
	Status ProviderStatus `json:"status,omitempty"`
}

func (p *Provider) GetStatusConditions() *[]metav1.Condition {
	return &p.Status.Conditions
}

func UpdatePaused(p Provider) Provider {
	meta.SetResourceCondition(&p, meta.ReadyCondition, metav1.ConditionTrue, meta.ReconciliationPausedReason, meta.ReconciliationPausedReason)
	return p
}

func ProviderNotReady(p Provider, reason, message string) Provider {
	meta.SetResourceCondition(&p, meta.ReadyCondition, metav1.ConditionFalse, reason, message)
	return p
}

func ProviderReady(p Provider) Provider {
	msg := "Completed"
	meta.SetResourceCondition(&p, meta.ReadyCondition, metav1.ConditionTrue, meta.ReconciliationSucceededReason, msg)
	return p
}

//+kubebuilder:object:root=true

// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Provider{}, &ProviderList{})
}
