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
package openstack

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"text/template"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	secretName        = "undistro-openstack-config"
	cloudConfTemplate = `[Global]
auth-url={{.AuthURL}}
tenant-name={{.ProjectName}}
domain-name=Default
region="RegionOne"
application-credential-id="{{.SecretID}}"
application-credential-secret="{{.SecretKey}}"
ca-file=/etc/config/cacert

[BlockStorage]
bs-version=v2

[Networking]
public-network-name=public
ipv6-support-disabled=false`
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type config struct {
	Clouds struct {
		Openstack struct {
			Auth struct {
				AuthURL                     string `json:"auth_url"`
				ApplicationCredentialID     string `json:"application_credential_id"`
				ApplicationCredentialSecret string `json:"application_credential_secret"`
				ProjectName                 string `json:"project_name"`
			} `json:"auth"`
		} `json:"openstack"`
	} `json:"clouds"`
}

type CloudConf struct {
	AuthURL     string
	ProjectName string
	SecretID    string
	SecretKey   string
}

func (c CloudConf) renderConf() (string, error) {
	tmpl, err := template.New("cloud.conf").Parse(cloudConfTemplate)
	if err != nil {
		return "", err
	}
	var credsFileStr bytes.Buffer
	err = tmpl.Execute(&credsFileStr, c)
	if err != nil {
		return "", err
	}
	return credsFileStr.String(), nil
}

func ReconcileCloudProvider(ctx context.Context, c client.Client, cl *appv1alpha1.Cluster) error {
	cfg := config{}
	err := json.Unmarshal(cl.Spec.InfrastructureProvider.ExtraConfiguration.Raw, &cfg)
	if err != nil {
		return err
	}
	conf := CloudConf{
		AuthURL:     cfg.Clouds.Openstack.Auth.AuthURL,
		ProjectName: cfg.Clouds.Openstack.Auth.ProjectName,
		SecretID:    cfg.Clouds.Openstack.Auth.ApplicationCredentialID,
		SecretKey:   cfg.Clouds.Openstack.Auth.ApplicationCredentialSecret,
	}
	cfgFile, err := conf.renderConf()
	if err != nil {
		return err
	}
	secret := corev1.Secret{}
	nm := client.ObjectKey{
		Name:      secretName,
		Namespace: undistro.Namespace,
	}
	err = c.Get(ctx, nm, &secret)
	if err != nil {
		return err
	}
	const (
		cloudHelm = "cloud-provider-openstack"
		version   = "1.22.0"
	)

	m := map[string]interface{}{
		"cloudconf": base64.StdEncoding.EncodeToString([]byte(cfgFile)),
		"cacert":    base64.StdEncoding.EncodeToString(secret.Data["caFile"]),
	}
	key := client.ObjectKey{
		Name:      hr.GetObjectName(cloudHelm, cl.Name),
		Namespace: cl.GetNamespace(),
	}
	release := appv1alpha1.HelmRelease{}
	err = c.Get(ctx, key, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		release, err = hr.Prepare(cloudHelm, "kube-system", cl.GetNamespace(), version, cl.Name, m)
		if err != nil {
			return err
		}
	}
	if release.Annotations == nil {
		release.Annotations = make(map[string]string)
	}
	release.Annotations[meta.SetupAnnotation] = cloudHelm
	err = hr.Install(ctx, c, release, cl)
	if err != nil {
		return err
	}
	if meta.InReadyCondition(release.Status.Conditions) {
		meta.SetResourceCondition(cl, meta.CloudProviderInstalledCondition, metav1.ConditionTrue, meta.CNIInstalledSuccessReason, "opestack cloud integration installed")
	}
	return nil
}

func ReconcileClusterSecret(ctx context.Context, c client.Client, cl *appv1alpha1.Cluster) error {
	secret := corev1.Secret{}
	nm := client.ObjectKey{
		Name:      secretName,
		Namespace: undistro.Namespace,
	}
	err := c.Get(ctx, nm, &secret)
	if err != nil {
		return err
	}
	dnsNameServers := string(secret.Data["dnsNameServers"])
	if cl.Spec.InfrastructureProvider.ExtraConfiguration == nil {
		return errors.New("clouds.yaml is required")
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(cl.Spec.InfrastructureProvider.ExtraConfiguration.Raw, &m)
	if err != nil {
		return err
	}
	for _, v := range m {
		clouds, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("invalid clouds.yaml")
		}
		for k := range clouds {
			e := corev1.EnvVar{
				Name:  "CLOUD_NAME",
				Value: k,
			}
			cl.Spec.InfrastructureProvider.Env = append(cl.Spec.InfrastructureProvider.Env, e)
		}
	}
	byt, err := yaml.JSONToYAML(cl.Spec.InfrastructureProvider.ExtraConfiguration.Raw)
	if err != nil {
		return err
	}
	clusterSecret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-cloud-config", cl.Name),
			Namespace: cl.Namespace,
		},
		Data: map[string][]byte{
			"clouds.yaml": byt,
			"cacert":      secret.Data["caFile"],
		},
	}
	_, err = util.CreateOrUpdate(ctx, c, &clusterSecret)
	if err != nil {
		return err
	}
	e := corev1.EnvVar{
		Name:  "DNS_NAME_SERVER",
		Value: dnsNameServers,
	}
	cl.Spec.InfrastructureProvider.Env = append(cl.Spec.InfrastructureProvider.Env, e)
	return nil
}
