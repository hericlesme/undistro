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
package cli

import (
	"encoding/base64"
	"os"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	configv1alpha1 "github.com/getupio-undistro/undistro/apis/config/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/helm"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	getters = getter.Providers{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}
)

type InitOptions struct {
	ConfigPath string
	genericclioptions.IOStreams
}

func NewInitOptions(streams genericclioptions.IOStreams) *InitOptions {
	return &InitOptions{
		IOStreams: streams,
	}
}

func (o *InitOptions) Complete(f *ConfigFlags, cmd *cobra.Command, args []string) error {
	o.ConfigPath = *f.ConfigFile
	if o.ConfigPath == "" {
		return cmdutil.UsageErrorf(cmd, "%s", "a file path need to be passed to --config flag")
	}
	if len(args) > 0 {
		return cmdutil.UsageErrorf(cmd, "%s", "too many arguments")
	}
	return nil
}

func (o *InitOptions) Validate() error {
	_, err := os.Stat(o.ConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func (o *InitOptions) RunInit(f cmdutil.Factory, cmd *cobra.Command) error {
	const undistroRepo = "http://repo.undistro.io"
	cfg := Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return errors.Errorf("unable to unmarshal config", err)
	}
	c, err := f.KubernetesClientSet()
	if err != nil {
		return err
	}
	ns := "undistro-system"
	secretName := "undistro-config"
	var clientOpts []getter.Option
	var private bool
	if cfg.Credentials.Username != "" && cfg.Credentials.Password != "" {
		private = true
		userb64 := base64.StdEncoding.EncodeToString([]byte(cfg.Credentials.Username))
		passb64 := base64.StdEncoding.EncodeToString([]byte(cfg.Credentials.Password))
		s := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: ns,
			},
			Data: map[string][]byte{
				"username": []byte(userb64),
				"password": []byte(passb64),
			},
		}
		s, err = c.CoreV1().Secrets(ns).Create(cmd.Context(), s, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
		opts, cleanup, err := helm.ClientOptionsFromSecret(*s)
		if err != nil {
			return err
		}
		defer cleanup()
		clientOpts = opts
	}
	chartRepo, err := helm.NewChartRepository(undistroRepo, getters, clientOpts)
	if err != nil {
		return err
	}
	err = chartRepo.DownloadIndex()
	if err != nil {
		return errors.Errorf("failed to download repository index: %w", err)
	}
	chartRepo.Index.SortEntries()
	chartName := "undistro"
	versions := chartRepo.Index.Entries[chartName]
	if versions.Len() == 0 {
		return errors.New("undistro chart not found")
	}
	version := versions[0]
	ch, err := chartRepo.Get(chartName, version.Version)
	if err != nil {
		return err
	}
	res, err := chartRepo.DownloadChart(ch)
	if err != nil {
		return err
	}
	chart, err := loader.LoadArchive(res)
	if err != nil {
		return err
	}
	restCfg, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	runner, err := helm.NewRunner(kube.NewInClusterRESTClientGetter(restCfg, ns), ns, log.Log)
	if err != nil {
		return err
	}
	wait := true
	hr := appv1alpha1.HelmRelease{
		Spec: appv1alpha1.HelmReleaseSpec{
			ReleaseName:     chartName,
			TargetNamespace: ns,
			Wait:            &wait,
			Timeout: &metav1.Duration{
				Duration: 5 * time.Minute,
			},
		},
	}
	_, err = runner.Install(hr, chart, chart.Values)
	if err != nil {
		return err
	}
	p := configv1alpha1.Provider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      chartName,
			Namespace: ns,
			Labels: map[string]string{
				meta.LabelProviderType: "core",
			},
		},
		Spec: configv1alpha1.ProviderSpec{
			ProviderName:    chartName,
			ProviderVersion: version.Version,
			AutoUpgrade:     true,
		},
	}
	if private {
		p.Spec.Repository.SecretRef = &corev1.LocalObjectReference{
			Name: secretName,
		}
	}
	return nil
}
