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
package coredns

import (
	"context"

	"github.com/getupio-undistro/undistro/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cfg = `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: system:coredns
rules:
  - apiGroups:
    - ""
    resources:
    - endpoints
    - services
    - pods
    - namespaces
    verbs:
    - list
    - watch
  - apiGroups:
    - discovery.k8s.io
    resources:
    - endpointslices
    verbs:
    - list
    - watch`

func EnsureComponentsConfig(ctx context.Context, c client.Client) error {
	objs, err := util.ToUnstructured([]byte(cfg))
	if err != nil {
		return err
	}
	for _, o := range objs {
		_, err = util.CreateOrUpdate(ctx, c, &o)
		if err != nil {
			return err
		}
	}
	return nil
}
