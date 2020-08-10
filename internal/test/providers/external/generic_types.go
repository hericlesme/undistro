/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

/*
package external defines the types for a generic external provider used for tests
*/
package external

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenericExternalObject is an object which is not actually managed by CAPI, but we wish to move with clusterctl
// using the "move" label on the resource.
// +kubebuilder:object:root=true
type GenericExternalObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

type GenericExternalObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericExternalObject `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&GenericExternalObject{}, &GenericExternalObjectList{},
	)
}
