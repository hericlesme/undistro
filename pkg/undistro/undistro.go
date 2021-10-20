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
package undistro

import (
	"context"
	"sync"

	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Namespace           = "undistro-system"
	MgmtClusterName     = "undistro"
	DefaultRepo         = "https://registry.undistro.io/chartrepo/library"
	loginAudienceSecret = "undistro-login-audience"
)

var (
	once            sync.Once
	requestAudience string
)

func GetRequestAudience() string {
	once.Do(func() {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			klog.Fatal(err)
		}
		c, err := client.New(cfg, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			klog.Fatal(err)
		}
		sec := corev1.Secret{}
		key := client.ObjectKey{
			Name:      loginAudienceSecret,
			Namespace: Namespace,
		}
		err = c.Get(context.Background(), key, &sec)
		if err != nil {
			if client.IgnoreNotFound(err) != nil {
				klog.Fatal(err)
				return
			}
			if sec.Data == nil {
				sec.Data = make(map[string][]byte)
			}
			requestAudience = util.RandomString(24)
			sec.Data["audience"] = []byte(requestAudience)
			err = c.Create(context.Background(), &sec)
			if err != nil {
				klog.Error(err)
			}
		}
		requestAudience = string(sec.Data["audience"])

	})
	return requestAudience
}

const LocalCluster = "undistro"

var KindCmdDestroy = `kind delete cluster --name "%s"`

var KindCfg = `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerPort: 6443
  apiServerAddress: 0.0.0.0
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: ClusterConfiguration
    apiServer:
      extraArgs:
        cors-allowed-origins: "http://*,https://*"
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP`

var TestResources = `---
apiVersion: v1
kind: Namespace
metadata:
  name: undistro-test
---
apiVersion: app.undistro.io/v1alpha1
kind: Cluster
metadata:
  name: undistro-test
  namespace: undistro-test
spec:
  paused: true
  kubernetesVersion: v1.18.2
  controlPlane:
    replicas: 3
    machineType: t3.large
  workers:
    - replicas: 3
      machineType: t3.large
  infrastructureProvider:
    name: aws
    sshKey: undistro
    region: us-east1
    flavor: ec2`
