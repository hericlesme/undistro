/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

/*
package bootstrap defines the types for a generic bootstrap provider used for tests
*/
package bootstrap

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GenericBootstrapConfigStatus struct {
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`
}

// +kubebuilder:object:root=true

type GenericBootstrapConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            GenericBootstrapConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type GenericBootstrapConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericBootstrapConfig `json:"items"`
}

// +kubebuilder:object:root=true

type GenericBootstrapConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

type GenericBootstrapConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericBootstrapConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GenericBootstrapConfig{}, &GenericBootstrapConfigList{},
		&GenericBootstrapConfigTemplate{}, &GenericBootstrapConfigTemplateList{},
	)
}
