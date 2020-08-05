/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

/*
package controlplane defines the types for a generic control plane provider used for tests
*/
package controlplane

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenericControlPlaneSpec struct {
	InfrastructureTemplate corev1.ObjectReference `json:"infrastructureTemplate"`
}

// +kubebuilder:object:root=true

type GenericControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GenericControlPlaneSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

type GenericControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GenericControlPlane{}, &GenericControlPlaneList{},
	)
}
