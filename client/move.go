/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package client

// MoveOptions carries the options supported by move.
type MoveOptions struct {
	// FromKubeconfig defines the kubeconfig to use for accessing the source management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	FromKubeconfig Kubeconfig

	// ToKubeconfig defines the kubeconfig to use for accessing the target management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	ToKubeconfig Kubeconfig

	// Namespace where the objects describing the workload cluster exists. If unspecified, the current
	// namespace will be used.
	Namespace string
}

func (c *undistroClient) Move(options MoveOptions) error {
	// Get the client for interacting with the source management cluster.
	fromCluster, err := c.clusterClientFactory(ClusterClientFactoryInput{kubeconfig: options.FromKubeconfig})
	if err != nil {
		return err
	}

	// Ensures the custom resource definitions required by undistro are in place.
	if err := fromCluster.ProviderInventory().EnsureCustomResourceDefinitions(); err != nil {
		return err
	}

	// Get the client for interacting with the target management cluster.
	toCluster, err := c.clusterClientFactory(ClusterClientFactoryInput{kubeconfig: options.ToKubeconfig})
	if err != nil {
		return err
	}

	// Ensures the custom resource definitions required by undistro are in place
	if err := toCluster.ProviderInventory().EnsureCustomResourceDefinitions(); err != nil {
		return err
	}

	// If the option specifying the Namespace is empty, try to detect it.
	if options.Namespace == "" {
		currentNamespace, err := fromCluster.Proxy().CurrentNamespace()
		if err != nil {
			return err
		}
		options.Namespace = currentNamespace
	}

	if err := fromCluster.ObjectMover().Move(options.Namespace, toCluster); err != nil {
		return err
	}

	return nil
}
