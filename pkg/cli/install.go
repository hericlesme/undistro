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
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	configv1alpha1 "github.com/getupio-undistro/undistro/apis/config/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/cloud"
	"github.com/getupio-undistro/undistro/pkg/helm"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	getters = getter.Providers{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}

	ns = "undistro-system"
)

type InstallOptions struct {
	ConfigPath  string
	ClusterName string
	genericclioptions.IOStreams
}

func NewInstallOptions(streams genericclioptions.IOStreams) *InstallOptions {
	return &InstallOptions{
		IOStreams: streams,
	}
}

func (o *InstallOptions) Complete(f *ConfigFlags, cmd *cobra.Command, args []string) error {
	o.ConfigPath = *f.ConfigFile
	switch len(args) {
	case 0:
		// do nothing
	case 1:
		o.ClusterName = args[0]
	default:
		return cmdutil.UsageErrorf(cmd, "%s", "too many arguments")
	}
	return nil
}

func (o *InstallOptions) Validate() error {
	if o.ConfigPath != "" {
		_, err := os.Stat(o.ConfigPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *InstallOptions) installProviders(ctx context.Context, c client.Client, providers []Provider, secretRef *corev1.LocalObjectReference) error {
	for _, p := range providers {
		secretName := fmt.Sprintf("undistro-%s-config", p.Name)
		secretData := make(map[string][]byte)
		valuesRef := make([]appv1alpha1.ValuesReference, 0)
		for k, v := range p.Configuration {
			vb64 := base64.StdEncoding.EncodeToString([]byte(v))
			secretData[k] = []byte(vb64)
			valuesRef = append(valuesRef, appv1alpha1.ValuesReference{
				Kind:       "Secret",
				Name:       secretName,
				ValuesKey:  k,
				TargetPath: k,
			})
		}
		s := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: ns,
			},
			Data: secretData,
		}
		hasDiff, err := util.CreateOrUpdate(ctx, c, &s)
		if err != nil {
			return err
		}
		provider := configv1alpha1.Provider{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p.Name,
				Namespace: ns,
			},
			Spec: configv1alpha1.ProviderSpec{
				ProviderName:      p.Name,
				ConfigurationFrom: valuesRef,
				Repository: configv1alpha1.Repository{
					SecretRef: secretRef,
				},
			},
		}
		if hasDiff {
			provider, err = cloud.Init(ctx, c, provider)
			if err != nil {
				return err
			}
		}
		_, err = util.CreateOrUpdate(ctx, c, &provider)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *InstallOptions) RunInstall(f cmdutil.Factory, cmd *cobra.Command) error {
	const undistroRepo = "http://repo.undistro.io"
	cfg := Config{}
	if o.ConfigPath != "" {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return errors.Errorf("unable to unmarshal config: %v", err)
		}
	}
	restCfg, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	secretName := "undistro-config"
	c, err := client.New(restCfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	restGetter := kube.NewInClusterRESTClientGetter(restCfg, "")
	if o.ClusterName != "" {
		byt, err := kubeconfig.FromSecret(cmd.Context(), c, util.ObjectKeyFromString(o.ClusterName))
		if err != nil {
			return err
		}
		restGetter = kube.NewMemoryRESTClientGetter(byt, "")
		restCfg, err = restGetter.ToRESTConfig()
		if err != nil {
			return err
		}
		c, err = client.New(restCfg, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return err
		}
	}
	var clientOpts []getter.Option
	var secretRef *corev1.LocalObjectReference
	if cfg.Credentials.Username != "" && cfg.Credentials.Password != "" {
		secretRef = &corev1.LocalObjectReference{
			Name: secretName,
		}
		userb64 := base64.StdEncoding.EncodeToString([]byte(cfg.Credentials.Username))
		passb64 := base64.StdEncoding.EncodeToString([]byte(cfg.Credentials.Password))
		s := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: ns,
			},
			Data: map[string][]byte{
				"username": []byte(userb64),
				"password": []byte(passb64),
			},
		}
		_, err = util.CreateOrUpdate(cmd.Context(), c, &s)
		if err != nil {
			return err
		}
		opts, cleanup, err := helm.ClientOptionsFromSecret(s)
		if err != nil {
			return err
		}
		defer cleanup()
		clientOpts = opts
	}
	err = c.List(cmd.Context(), &configv1alpha1.ProviderList{})
	if err != nil {
		chartRepo, err := helm.NewChartRepository(undistroRepo, getters, clientOpts)
		if err != nil {
			return err
		}
		err = chartRepo.DownloadIndex()
		if err != nil {
			return fmt.Errorf("failed to download repository index: %w", err)
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
		runner, err := helm.NewRunner(restGetter, ns, log.Log)
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
				Repository: configv1alpha1.Repository{
					SecretRef: secretRef,
				},
			},
		}
		err = c.Create(cmd.Context(), &p)
		if err != nil {
			return err
		}
	}
	return o.installProviders(cmd.Context(), c, cfg.Providers, secretRef)
}

func NewCmdInstall(f *ConfigFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewInstallOptions(streams)
	cmd := &cobra.Command{
		Use:                   "install [cluster namespace/cluster name]",
		DisableFlagsInUseLine: true,
		Short:                 "Install UnDistro",
		Long: LongDesc(`Install UnDistro.
		If cluster argument exists UnDistro will be installed in this remote cluster.
		If config file exists UnDistro will be installed using file's configurations`),
		Example: Examples(`
		# Install UnDistro in local cluster
		undistro install
		# Install UnDistro in remote cluster
		undistro install undistro-production/cool-product-cluster
		# Install UnDistro with configuration file
		undistro --config undistro-config.yaml install undistro-production/cool-product-cluster
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.RunInstall(cmdutil.NewFactory(f), cmd))
		},
	}
	return cmd
}
