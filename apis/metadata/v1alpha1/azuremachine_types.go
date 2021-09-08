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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AzureMachineSpec defines the desired state of AzureMachine
type AzureMachineSpec struct {
	InstanceType      string                  `json:"instanceType,omitempty"`
	AvailabilityZones []string                `json:"availabilityZones,omitempty"`
	Vcpus             string                  `json:"vcpus,omitempty"`
	Memory            string                  `json:"memory,omitempty"`
	ProviderRef       *corev1.ObjectReference `json:"providerRef,omitempty"`
}

// AzureMachineStatus defines the observed state of AzureMachine
type AzureMachineStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// AzureMachine is the Schema for the awsmachines API
type AzureMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureMachineSpec   `json:"spec,omitempty"`
	Status AzureMachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AzureMachineList contains a list of AzureMachine
type AzureMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureMachine{}, &AzureMachineList{})
}
