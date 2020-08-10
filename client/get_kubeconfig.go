/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package client

import "github.com/pkg/errors"

//GetKubeconfigOptions carries all the options supported by GetKubeconfig
type GetKubeconfigOptions struct {
	// Kubeconfig defines the kubeconfig to use for accessing the management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	Kubeconfig Kubeconfig

	// Namespace is the namespace in which secret is placed.
	Namespace string

	// WorkloadClusterName is the name of the workload cluster.
	WorkloadClusterName string
}

func (c *undistroClient) GetKubeconfig(options GetKubeconfigOptions) (string, error) {
	// gets access to the management cluster
	clusterClient, err := c.clusterClientFactory(ClusterClientFactoryInput{kubeconfig: options.Kubeconfig})
	if err != nil {
		return "", err
	}

	if options.Namespace == "" {
		currentNamespace, err := clusterClient.Proxy().CurrentNamespace()
		if err != nil {
			return "", err
		}
		if currentNamespace == "" {
			return "", errors.New("failed to identify the current namespace. Please specify the namespace where the workload cluster exists")
		}
		options.Namespace = currentNamespace
	}

	return clusterClient.WorkloadCluster().GetKubeconfig(options.WorkloadClusterName, options.Namespace)

}
