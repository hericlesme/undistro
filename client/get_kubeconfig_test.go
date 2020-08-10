/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package client

import (
	"testing"

	"github.com/getupcloud/undistro/client/cluster"
	"github.com/getupcloud/undistro/internal/test"
	. "github.com/onsi/gomega"
)

func Test_clusterctlClient_GetKubeconfig(t *testing.T) {

	configClient := newFakeConfig()
	kubeconfig := cluster.Kubeconfig{Path: "kubeconfig", Context: "mgmt-context"}
	clusterClient := &fakeClusterClient{
		kubeconfig: kubeconfig,
	}
	// create a clusterctl client where the proxy returns an empty namespace
	clusterClient.fakeProxy = test.NewFakeProxy().WithNamespace("")
	badClient := newFakeClient(configClient).WithCluster(clusterClient)

	tests := []struct {
		name      string
		client    *fakeClient
		options   GetKubeconfigOptions
		expectErr bool
	}{
		{
			name:      "returns error if unable to get client for mgmt cluster",
			client:    fakeEmptyCluster(),
			expectErr: true,
		},
		{
			name:      "returns error if unable namespace is empty",
			client:    badClient,
			options:   GetKubeconfigOptions{Kubeconfig: Kubeconfig(kubeconfig)},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			config, err := tt.client.GetKubeconfig(tt.options)
			if tt.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(config).ToNot(BeEmpty())
		})
	}
}
