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
package cloudutil

import (
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func RemoveDuplicateEnv(envs []corev1.EnvVar) []corev1.EnvVar {
	nMap := make(map[string]corev1.EnvVar)
	for _, t := range envs {
		nMap[t.Name] = t
	}
	res := make([]corev1.EnvVar, 0)
	for _, v := range nMap {
		res = append(res, v)
	}
	return res
}

func RemoveDuplicateNetwork(n []appv1alpha1.NetworkSpec) []appv1alpha1.NetworkSpec {
	nMap := make(map[appv1alpha1.NetworkSpec]struct{})
	for _, t := range n {
		nMap[t] = struct{}{}
	}
	res := make([]appv1alpha1.NetworkSpec, 0)
	for k := range nMap {
		res = append(res, k)
	}
	return res
}
