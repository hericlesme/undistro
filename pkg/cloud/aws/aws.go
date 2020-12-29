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
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/cloud/aws/cloudformation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	GroupVersion = schema.GroupVersion{Group: "infrastructure.cluster.x-k8s.io", Version: "v1alpha3"}
)

const (
	defaultAWSRegion       = "us-east-1"
	name                   = "undistro-aws-config"
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

func (c awsCredentials) setBase64EncodedAWSDefaultProfile(ctx context.Context, cl client.Client, secret *corev1.Secret) (appv1alpha1.ValuesReference, error) {
	profile, err := c.renderAWSDefaultProfile()
	if err != nil {
		return appv1alpha1.ValuesReference{}, err
	}
	secret.Data[key] = []byte(profile)
	err = cl.Update(ctx, secret)
	if err != nil {
		return appv1alpha1.ValuesReference{}, err
	}
	return appv1alpha1.ValuesReference{
		Kind:       "Secret",
		Name:       name,
		ValuesKey:  key,
		TargetPath: key,
	}, nil
}

// Init providers
func Init(ctx context.Context, c client.Client, cfg []appv1alpha1.ValuesReference, version string) ([]appv1alpha1.ValuesReference, error) {
	cred, secret, err := getCreds(ctx, c)
	if err != nil {
		return cfg, err
	}
	v, err := cred.setBase64EncodedAWSDefaultProfile(ctx, c, secret)
	if err != nil {
		return cfg, err
	}
	cfg = append(cfg, v)
	return cfg, nil
}

// Upgrade providers
func Upgrade(ctx context.Context, c client.Client, cfg []appv1alpha1.ValuesReference, version string) ([]appv1alpha1.ValuesReference, error) {
	cred, _, err := getCreds(ctx, c)
	if err != nil {
		return cfg, err
	}
	err = reconcileCloudformation(cred)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func reconcileCloudformation(cred awsCredentials) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cred.Region),
		Credentials: credentials.NewStaticCredentials(
			cred.AccessKeyID,
			cred.SecretAccessKey,
			cred.SessionToken,
		),
	})
	if err != nil {
		return err
	}
	cfnSvc := cloudformation.NewService(cfn.New(sess))
	return cfnSvc.ReconcileBootstrapStack(cloudformation.Template)
}

func getCreds(ctx context.Context, c client.Client) (awsCredentials, *corev1.Secret, error) {
	secret := corev1.Secret{}
	nm := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}
	err := c.Get(ctx, nm, &secret)
	if err != nil {
		return awsCredentials{}, nil, err
	}
	cred, err := credentialsFromSecret(&secret)
	if err != nil {
		return awsCredentials{}, nil, err
	}
	return cred, &secret, nil
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
	b, ok := secret.Data[key]
	if !ok {
		return ""
	}
	return string(b)
}
