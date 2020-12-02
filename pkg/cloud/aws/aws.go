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
package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"text/template"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	defaultAWSRegion       = "us-east-1"
	name                   = "aws-provider-config"
	namespace              = "undistro-system"
	key                    = "credentials"
	awsCredentialsTemplate = `[default]
aws_access_key_id = {{ .AccessKeyID }}
aws_secret_access_key = {{ .SecretAccessKey }}
region = {{ .Region }}
{{if .SessionToken }}
aws_session_token = {{ .SessionToken }}
{{end}}`
)

type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

func (c awsCredentials) renderAWSDefaultProfile() (string, error) {
	tmpl, err := template.New("AWS Credentials").Parse(awsCredentialsTemplate)
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

func (c awsCredentials) setBase64EncodedAWSDefaultProfile(secret *corev1.Secret) (appv1alpha1.ValuesReference, error) {
	profile, err := c.renderAWSDefaultProfile()
	if err != nil {
		return appv1alpha1.ValuesReference{}, err
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(profile))
	secret.Data[key] = []byte(b64)
	return appv1alpha1.ValuesReference{
		Kind:       "Secret",
		Name:       name,
		ValuesKey:  key,
		TargetPath: key,
	}, nil
}

// Init providers
func Init(ctx context.Context, c client.Client, cfg []appv1alpha1.ValuesReference, version string) ([]appv1alpha1.ValuesReference, error) {
	secret := corev1.Secret{}
	nm := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}
	err := c.Get(ctx, nm, &secret)
	if err != nil {
		return cfg, err
	}
	cred, err := credentialsFromSecret(&secret)
	if err != nil {
		return cfg, err
	}
	v, err := cred.setBase64EncodedAWSDefaultProfile(&secret)
	if err != nil {
		return cfg, err
	}
	cfg = append(cfg, v)
	return cfg, nil
}

func credentialsFromSecret(s *corev1.Secret) (awsCredentials, error) {
	cred := awsCredentials{
		AccessKeyID:     getData(s, "accessKeyID"),
		SecretAccessKey: getData(s, "secretAccessKey"),
		Region:          getData(s, "region"),
		SessionToken:    getData(s, "sessionToken"),
	}
	if cred.Region == "" {
		cred.Region = defaultAWSRegion
	}
	return cred, nil
}

func getData(secret *corev1.Secret, key string) string {
	b64, ok := secret.Data[key]
	if !ok {
		return ""
	}
	s, err := base64.StdEncoding.DecodeString(string(b64))
	if err != nil {
		return string(b64)
	}
	return string(s)
}
