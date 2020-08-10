/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

/*
package infrastructure defines the types for a generic infrastructure provider used for tests
*/
package infrastructure

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type GenericInfrastructureCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

type GenericInfrastructureClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericInfrastructureCluster `json:"items"`
}

// +kubebuilder:object:root=true

type GenericInfrastructureMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

type GenericInfrastructureMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericInfrastructureMachine `json:"items"`
}

// +kubebuilder:object:root=true

type GenericInfrastructureMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

type GenericInfrastructureMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericInfrastructureMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GenericInfrastructureCluster{}, &GenericInfrastructureClusterList{},
		&GenericInfrastructureMachine{}, &GenericInfrastructureMachineList{},
		&GenericInfrastructureMachineTemplate{}, &GenericInfrastructureMachineTemplateList{},
	)
}
