/*
Copyright 2020 The UnDistro authors

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

	configv1alpha1 "github.com/getupio-undistro/undistro/apis/config/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/cloud/aws"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Init providers
func Init(ctx context.Context, c client.Client, p configv1alpha1.Provider) (configv1alpha1.Provider, error) {
	var err error
	switch p.Spec.ProviderName {
	case "undistro-aws":
		p.Spec.ConfigurationFrom, err = aws.Init(ctx, c, p.Spec.ConfigurationFrom, p.Spec.ProviderVersion)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

// Upgrade providers
func Upgrade(ctx context.Context, c client.Client, p configv1alpha1.Provider) (configv1alpha1.Provider, error) {
	var err error
	switch p.Spec.ProviderName {
	case "aws":
		p.Spec.ConfigurationFrom, err = aws.Upgrade(ctx, c, p.Spec.ConfigurationFrom, p.Spec.ProviderVersion)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
