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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:Enum=HTTP;Helm
type RepositoryType string

const (
	HTTPRepository RepositoryType = "HTTP"
	HelmRepository RepositoryType = "Helm"
)

// ProviderSpec defines the desired state of Provider
type ProviderSpec struct {
	// +kubebuilder:default=false
	Paused            bool                         `json:"paused,omitempty"`
	ProviderName      string                       `json:"providerName,omitempty"`
	ProviderVersion   string                       `json:"providerVersion,omitempty"`
	Repository        Repository                   `json:"repository,omitempty"`
	ConfigurationFrom *corev1.LocalObjectReference `json:"configurationFrom,omitempty"`
	// +kubebuilder:default=false
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`
}

type Repository struct {
	// +kubebuilder:default="https://charts.undistro.io"
	URL string `json:"url,omitempty"`
	// +kubebuilder:default=Helm
	Type      RepositoryType               `json:"type,omitempty"`
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`
}

// ProviderStatus defines the observed state of Provider
type ProviderStatus struct {
	// ObservedGeneration is the last observed generation.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions         []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	HelmReleaseName    string             `json:"helmReleaseName,omitempty"`
	LastAppliedVersion string             `json:"lastAppliedVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=providers,scope=Namespaced
// +kubebuilder:printcolumn:name="Provider Name",type="string",JSONPath=".spec.providerName"
// +kubebuilder:printcolumn:name="Provider Version",type="string",JSONPath=".spec.providerVersion"
// +kubebuilder:printcolumn:name="Provider Type",type="string",JSONPath=".spec.repository.type"
// +kubebuilder:printcolumn:name="Provider URL",type="string",JSONPath=".status.url"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// Provider is the Schema for the providers API
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec   `json:"spec,omitempty"`
	Status ProviderStatus `json:"status,omitempty"`
}

func (p *Provider) GetNamespace() string {
	if p.Namespace == "" {
		return "default"
	}
	return p.Namespace
}

// +kubebuilder:object:root=true

// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Provider{}, &ProviderList{})
}
