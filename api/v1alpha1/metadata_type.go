/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/version"
)

// +kubebuilder:object:root=true

// Metadata for a provider repository
type Metadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	ReleaseSeries []ReleaseSeries `json:"releaseSeries"`
}

// ReleaseSeries maps a provider release series (major/minor) with a API Version of Cluster API (contract).
type ReleaseSeries struct {
	// Major version of the release series
	Major uint `json:"major,omitempty"`

	// Minor version of the release series
	Minor uint `json:"minor,omitempty"`

	// Contract defines the Cluster API contract supported by this series.
	//
	// The value is an API Version, e.g. `v1alpha3`.
	Contract string `json:"contract,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Metadata{})
}

// GetReleaseSeriesForVersion returns the release series for a given version.
func (m *Metadata) GetReleaseSeriesForVersion(version *version.Version) *ReleaseSeries {
	for _, releaseSeries := range m.ReleaseSeries {
		if version.Major() == releaseSeries.Major && version.Minor() == releaseSeries.Minor {
			return &releaseSeries
		}
	}

	return nil
}
