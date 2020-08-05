/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package v1alpha1

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_Provider_ManifestLabel(t *testing.T) {
	type fields struct {
		provider     string
		providerType ProviderType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "core provider remains the same",
			fields: fields{
				provider:     "cluster-api",
				providerType: CoreProviderType,
			},
			want: "cluster-api",
		},
		{
			name: "kubeadm bootstrap",
			fields: fields{
				provider:     "kubeadm",
				providerType: BootstrapProviderType,
			},
			want: "bootstrap-kubeadm",
		},
		{
			name: "other bootstrap providers gets prefix",
			fields: fields{
				provider:     "xx",
				providerType: BootstrapProviderType,
			},
			want: "bootstrap-xx",
		},
		{
			name: "kubeadm control-plane",
			fields: fields{
				provider:     "kubeadm",
				providerType: ControlPlaneProviderType,
			},
			want: "control-plane-kubeadm",
		},
		{
			name: "other control-plane providers gets prefix",
			fields: fields{
				provider:     "xx",
				providerType: ControlPlaneProviderType,
			},
			want: "control-plane-xx",
		},
		{
			name: "infrastructure providers gets prefix",
			fields: fields{
				provider:     "xx",
				providerType: InfrastructureProviderType,
			},
			want: "infrastructure-xx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			p := &Provider{
				ProviderName: tt.fields.provider,
				Type:         string(tt.fields.providerType),
			}
			g.Expect(p.ManifestLabel()).To(Equal(tt.want))
		})
	}
}
