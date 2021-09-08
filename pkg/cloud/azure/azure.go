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

package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
)

const (
	name                     = "undistro-azure-config"
	defaultAzureRegion       = "eastus"
	azureCredentialsTemplate = `[default]
azure_client_id = {{ .ClientID }}
azure_client_secret_key = {{ .ClientSecret }}
azure_subscription_id = {{ .SubscriptionID }}
azure_tenant_id = {{ .TenantID }}`
)

var flavors = map[string][]string{
	//// Check doc for k8s versions
	appv1alpha1.VM.String(): {
		"v1.19.12", "v1.19.13", "v1.20.8", "1.20.9", "v1.21.2", "v1.21.3",
	},
	appv1alpha1.AKS.String(): { //// ManagedCluster>
		"v1.19.8", "v1.20.7", "v1.21.2",
	},
}

//// WIP
func GetFlavors(_ context.Context, p metadatav1alpha1.Provider, subscriptionID string) ([]client.Object, error) {
	apiV := "latest"
	resTyp := "managedClusters"
	///// List by region.
	restUrl := "https://management.azure.com/subscriptions/%s/providers/Microsoft.ContainerService/locations/%s/orchestrators?api-version=%s&resource-type=%s"
	//// List globally.
	// https://management.azure.com/subscriptions/{subscriptionId}/providers/Microsoft.ContainerService/managedClusters?api-version=2021-05-01"
	aksUrl := fmt.Sprintf(restUrl, subscriptionID, defaultAzureRegion, apiV, resTyp)
}

//// Check.
func GetMachineMetadata(_ context.Context, p metadatav1alpha1.Provider) ([]client.Object, error) {
	ref := &corev1.ObjectReference{
		APIVersion: p.APIVersion,
		Kind:       p.Kind,
		Name:       p.Name,
		Namespace:  p.Namespace,
	}
	typeMeta := metav1.TypeMeta{
		APIVersion: metadatav1alpha1.GroupVersion.String(),
		Kind:       "AzureMachine",
	}
	specs := make([]metadatav1alpha1.AzureMachineSpec, 0)
	err := json.Unmarshal(instanceTypes, &specs)
	if err != nil {
		return nil, err
	}
	objs := make([]client.Object, len(specs))
	for i, spec := range specs {
		o := metadatav1alpha1.AzureMachine{
			TypeMeta: typeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name: spec.InstanceType,
			},
			Spec: spec,
		}
		o.Spec.ProviderRef = ref
		objs[i] = &o
	}
	return objs, nil
}

func kindByFlavor(flavor string) string {
	switch flavor {
	case "vm":
		return "AzureMachinePool"
	case "aks":
		return "AzureManagedMachinePool"
	}
	return ""
}

//// Check.
func launchTemplateRef(ctx context.Context, u unstructured.Unstructured) (appv1alpha1.LaunchTemplateReference, error) {
	var (
		ref     appv1alpha1.LaunchTemplateReference
		version int64
		ok      bool
		err     error
	)
	switch u.GetKind() {
	case "AzureMachinePool":
		ref.ID, ok, err = unstructured.NestedString(u.Object, "spec", "azureLaunchTemplate", "id")
		if !ok || err != nil {
			return ref, err
		}
		version, ok, err = unstructured.NestedInt64(u.Object, "spec", "azureLaunchTemplate", "versionNumber")
		if !ok || err != nil {
			return ref, err
		}
		ref.Version = strconv.Itoa(int(version))
	}
	return ref, nil
}

//// Check.
func ReconcileLaunchTemplate(ctx context.Context, r client.Client, cl *appv1alpha1.Cluster) error {
	for i := range cl.Spec.Workers {
		if !cl.Spec.InfrastructureProvider.IsManaged() {
			key := client.ObjectKey{
				Name:      fmt.Sprintf("%s-mp-%d", cl.Name, i),
				Namespace: cl.GetNamespace(),
			}
			u := unstructured.Unstructured{}
			u.SetAPIVersion("infrastructure.cluster.x-k8s.io/v1alpha4")
			u.SetKind(kindByFlavor(cl.Spec.InfrastructureProvider.Flavor))
			err := r.Get(ctx, key, &u)
			if err != nil {
				return client.IgnoreNotFound(err)
			}
			ref, err := launchTemplateRef(ctx, u)
			if err != nil {
				return err
			}
			cl.Spec.Workers[i].LaunchTemplateReference = ref
		}
	}
	return nil
}

type AzureCredentials struct {
	EnvVars        map[string]string
	SubscriptionID string
}

func setupEnv(envVars map[string]string) {
	for envVar, val := range envVars {
		os.Setenv(envVar, val)
	}
}

func getAuthorizer() (autorest.Authorizer, error) {
	auth, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func getLocations(ctx context.Context, subscriptionID string) (subscriptions.LocationListResult, error) {
	var locations subscriptions.LocationListResult
	client := subscriptions.NewClient()
	auth, err := getAuthorizer()
	if err != nil {
		return locations, errors.New("unable to authenticate to Azure")
	}
	client.Authorizer = auth
	if locations, err := client.ListLocations(ctx, scriptionID, false); err != nil {
		return locations, errors.New("unable to get Azure locations")
	}
	return locations, nil
}

func locationsAsString(locations []subscriptions.LocationListResult) []string {
	locNames := make([]string, len(locations))
	for c := 0; c < len(locations); c++ {
		locNames[c] = &locations[c].Value.Name

	}
	return locNames
}

func GetRegions(ctx context.Context, c client.Client, log logr.Logger) []string {
	s := corev1.Secret{}
	key := client.ObjectKey{
		name:      name,
		namespace: undistro.Namespace,
	}
	err := c.Get(ctx, key, &s)
	if err != nil {
		log.Info("unable to get azure credentials", "err", err)
		return nil
	}
	ac := AzureCredentials{
		map[string]string{
			"AZURE_TENANT_ID":     util.GetData(s, "tenantID"),
			"AZURE_CLIENT_ID":     util.GetData(s, "clientID"),
			"AZURE_CLIENT_SECRET": util.GetData(s, "clientSecret"),
		}, util.GetData(s, "subscriptionID"),
	}
	setupEnv(ac.EnvVars)
	locations, err := getLocations(ctx, ac.SubscriptionID)
	if err != nil {
		log.Info(err)
		return nil
	}
	locNames := locationsAsString(locations)
	return locNames
}

func ReconcileNetwork(ctx context.Context, r client.Client, cl *appv1alpha1.Cluster, capiCluster *capi.Cluster) error {
	u := unstructured.Unstructured{}
	key := client.ObjectKey{}
	if cl.Spec.InfrastructureProvider.IsManaged() {
		u.SetGroupVersionKind(capiCluster.Spec.ControlPlaneRef.GroupVersionKind())
		key = client.ObjectKey{
			Name:      capiCluster.Spec.ControlPlaneRef.Name,
			Namespace: capiCluster.Spec.ControlPlaneRef.Namespace,
		}
	} else {
		u.SetGroupVersionKind(capiCluster.Spec.InfrastructureRef.GroupVersionKind())
		key = client.ObjectKey{
			Name:      capiCluster.Spec.InfrastructureRef.Name,
			Namespace: capiCluster.Spec.InfrastructureRef.Namespace,
		}
	}
	err := r.Get(ctx, key, &u)
	if err != nil {
		return client.IgnoreNotFound(err)
	}
	return clusterNetwork(ctx, cl, u)
}

func clusterNetwork(ctx context.Context, cl *appv1alpha1.Cluster, u unstructured.Unstructured) error {
	host, ok, err := unstructured.NestedString(u.Object, "spec", "controlPlaneEndpoint", "host")
	if err != nil {
		return err
	}
	if ok && host != "" {
		cl.Spec.ControlPlane.Endpoint.Host = host
	}
	port, ok, err := unstructured.NestedInt64(u.Object, "spec", "controlPlaneEndpoint", "port")
	if err != nil {
		return err
	}
	if ok && port != 0 {
		cl.Spec.ControlPlane.Endpoint.Port = int32(port)
	}
	vnet, ok, err := unstructured.NestedMap(u.Object, "spec", "networkSpec", "vnet")
	if err != nil {
		return err
	}
	if ok {
		id, ok := vnet["id"]
		if ok {
			cl.Spec.Network.VPC.ID = id.(string)
		}
		cidrs, ok := vnet["cidrBlocks"]
		if ok {
			cl.Spec.Network.VPC.CIDRBlock = cidrs.([]string)[0]
		}
	}
	subnets, ok, err := unstructured.NestedSlice(u.Object, "spec", "networkSpec", "subnets")
	if err != nil {
		return err
	}
	if ok {
		for i, s := range subnets {
			subnet, ok := s.(map[string]interface{})
			if !ok {
				return errors.Errorf("unable to reconcile subnets for cluster %s/%s", cl.Namespace, cl.Name)
			}
			n := appv1alpha1.NetworkSpec{
				ID:        subnet["id"].(string),
				CIDRBlock: subnet["cidrBlock"].([]string)[i],
				IsPublic:  subnet["isPublic"].(bool),
				Role:      subnet["role"].(string),
			}
			cl.Spec.Network.Subnets = append(cl.Spec.Network.Subnets, n)
		}
		cl.Spec.Network.Subnets = util.RemoveDuplicateNetwork(cl.Spec.Network.Subnets)
	}
	return nil
}
