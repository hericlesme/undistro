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
package cloud

import (
	"context"
	"strconv"

	"github.com/Masterminds/semver/v3"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	metadatav1alpha1 "github.com/getupio-undistro/undistro/apis/metadata/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/cloud/aws"
	"github.com/getupio-undistro/undistro/pkg/cloud/openstack"
	"github.com/go-logr/logr"
	capi "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Account interface {
	GetID() string
	GetUsername() string
	IsRoot() bool
}

type MetadataFunc func(context.Context, metadatav1alpha1.Provider) ([]client.Object, error)

func RegionNames(provider metadatav1alpha1.Provider) []string {
	switch provider.Name {
	case appv1alpha1.Amazon.String():
		return aws.Regions
	}
	return nil
}

func GetFlavors(provider metadatav1alpha1.Provider) MetadataFunc {
	switch provider.Name {
	case appv1alpha1.Amazon.String():
		return aws.GetFlavors
	}
	return nil
}

func GetMachineMetadata(provider metadatav1alpha1.Provider) MetadataFunc {
	switch provider.Name {
	case appv1alpha1.Amazon.String():
		return aws.GetMachineMetadata
	}
	return nil
}

// ReconcileNetwork from clouds
func ReconcileNetwork(ctx context.Context, r client.Client, cl *appv1alpha1.Cluster, capiCluster *capi.Cluster) error {
	if capiCluster.Spec.InfrastructureRef == nil || capiCluster.Spec.ControlPlaneRef == nil {
		return nil
	}
	if capiCluster.Spec.ClusterNetwork != nil {
		cl.Spec.Network.ClusterNetwork = *capiCluster.Spec.ClusterNetwork
	}
	switch cl.Spec.InfrastructureProvider.Name {
	case appv1alpha1.Amazon.String():
		return aws.ReconcileNetwork(ctx, r, cl, capiCluster)
	case appv1alpha1.OpenStack.String():
		return openstack.ReconcileNetwork(ctx, r, cl, capiCluster)
	}
	return nil
}

// ReconcileClusterConf from clouds
func ReconcileClusterConf(ctx context.Context, r client.Client, cl *appv1alpha1.Cluster, capiCluster *capi.Cluster) error {
	switch cl.Spec.InfrastructureProvider.Name {
	case appv1alpha1.Amazon.String():
		return aws.ReconcileLaunchTemplate(ctx, r, cl, capiCluster)
	case appv1alpha1.OpenStack.String():
		return openstack.ReconcileClusterSecret(ctx, r, cl)
	}
	return nil
}

// ReconcileProviderIntegration from clouds
func ReconcileProviderIntegration(ctx context.Context, r client.Client, log logr.Logger, cl *appv1alpha1.Cluster) error {
	switch cl.Spec.InfrastructureProvider.Name {
	case appv1alpha1.OpenStack.String():
		return openstack.ReconcileCloudProvider(ctx, r, log, cl)
	}
	return nil
}

func CalicoValues(cl *appv1alpha1.Cluster) map[string]interface{} {
	values := make(map[string]interface{})
	switch cl.Spec.InfrastructureProvider.Flavor {
	case appv1alpha1.EKS.String():
		values["vxlan"] = true
	default:
		values["vxlan"] = false
	}
	return values
}

func GetAccount(ctx context.Context, c client.Client, cl *appv1alpha1.Cluster) (Account, error) {
	switch cl.Spec.InfrastructureProvider.Name {
	case "aws":
		return aws.NewAccount(ctx, c)
	}
	return nil, nil
}

func DefaultRegion(infra string) string {
	switch infra {
	case "aws":
		return aws.DefaultAWSRegion
	}
	return ""
}

// https://github.com/coredns/deployment/blob/master/kubernetes/CoreDNS-k8s_version.md
func CoreDNSVersion(k8sVersion string) string {
	v := semver.MustParse(k8sVersion)
	minor := strconv.Itoa(int(v.Minor()))
	switch minor {
	case "19", "20":
		return "v1.7.0"
	default:
		return "v1.8.4"
	}
}
