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
package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/capi"
	"github.com/getupio-undistro/undistro/pkg/certmanager"
	undistrofs "github.com/getupio-undistro/undistro/pkg/fs"
	"github.com/getupio-undistro/undistro/pkg/helm"
	"github.com/getupio-undistro/undistro/pkg/internalautohttps"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/retry"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

var (
	getters = getter.Providers{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}
)

const (
	undistroRepo = "https://registry.undistro.io/chartrepo/library"
	ns           = undistro.Namespace
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

func (o *InstallOptions) installChart(ctx context.Context, c client.Client, restGetter genericclioptions.RESTClientGetter, chartName string, overrideValuesMap map[string]interface{}) (*appv1alpha1.HelmRelease, error) {
	overrideValues := &apiextensionsv1.JSON{}
	if overrideValuesMap == nil {
		overrideValuesMap = make(map[string]interface{})
	}
	if len(overrideValuesMap) > 0 {
		byt, err := json.Marshal(overrideValuesMap)
		if err != nil {
			return nil, err
		}
		overrideValues.Raw = byt
	}
	var (
		file  fs.File
		found bool
	)
	err := fs.WalkDir(undistrofs.ChartFS, "chart", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if found {
			return nil
		}

		if d.IsDir() {
			return nil
		}
		file, err = undistrofs.ChartFS.Open(path)
		if err != nil {
			return err
		}
		if strings.Contains(path, chartName) {
			found = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	chart, err := loader.LoadArchive(file)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(o.IOStreams.Out, "Installing %s version %s\n", chart.Name(), chart.AppVersion())
	for _, dep := range chart.Dependencies() {
		fmt.Fprintf(o.IOStreams.Out, "Installing %s version %s\n", dep.Name(), dep.AppVersion())
	}
	wait := true
	forceUpgrade := true
	reset := false
	history := 0
	_, dev := os.LookupEnv("DEV_ENV")
	hr := appv1alpha1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appv1alpha1.GroupVersion.String(),
			Kind:       "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      chartName,
			Namespace: ns,
		},
		Spec: appv1alpha1.HelmReleaseSpec{
			Paused: dev,
			Chart: appv1alpha1.ChartSource{
				RepoChartSource: appv1alpha1.RepoChartSource{
					RepoURL: undistroRepo,
					Name:    chartName,
					Version: chart.Metadata.Version,
				},
			},
			ReleaseName:     chartName,
			TargetNamespace: ns,
			Values:          overrideValues,
			Wait:            &wait,
			MaxHistory:      &history,
			ResetValues:     &reset,
			ForceUpgrade:    &forceUpgrade,
			Timeout: &metav1.Duration{
				Duration: 5 * time.Minute,
			},
		},
	}
	err = retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
		runner, err := helm.NewRunner(restGetter, ns, log.Log)
		if err != nil {
			return err
		}
		m := make(map[string]interface{})
		if overrideValues.Raw != nil {
			err = json.Unmarshal(overrideValues.Raw, &m)
			if err != nil {
				return err
			}
		}
		chart.Values = util.MergeMaps(chart.Values, m)
		rel, _ := runner.ObserveLastRelease(hr)
		if rel == nil {
			_, err = runner.Install(hr, chart, chart.Values)
			if err != nil {
				return err
			}
		} else if rel.Info.Status == release.StatusDeployed {
			_, err = runner.Upgrade(hr, chart, chart.Values)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &hr, err
}
func (o *InstallOptions) installCore(ctx context.Context, c client.Client, restGetter genericclioptions.RESTClientGetter, cfg map[string]interface{}, charts ...string) ([]appv1alpha1.HelmRelease, error) {
	var depsMap = map[string]string{
		"cert-manager": certmanager.TestResources,
		"cluster-api":  capi.TestResources,
		"undistro":     undistro.TestResources,
	}
	hrs := make([]appv1alpha1.HelmRelease, len(charts))
	for i, chart := range charts {
		values := defaultValues(ctx, c, chart)
		cfgValues, ok := cfg[chart]
		if ok {
			values = util.MergeMaps(cfgValues.(map[string]interface{}), values)
		}
		hr, err := o.installChart(ctx, c, restGetter, chart, values)
		if err != nil {
			return nil, err
		}
		testRes := depsMap[chart]
		if testRes != "" {
			objs, err := util.ToUnstructured([]byte(testRes))
			if err != nil {
				return nil, err
			}
			for _, o := range objs {
				err = retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
					_, err = util.CreateOrUpdate(ctx, c, &o)
					return err
				})
			}
		}
		if chart == "undistro" {
			isLocal, ok := values["local"].(bool)
			if ok && isLocal {
				fmt.Fprintf(o.IOStreams.Out, "Installing local certificates\n")
				err = internalautohttps.InstallLocalCert(ctx, c)
				if err != nil {
					return nil, err
				}
			}
		}
		hrs[i] = *hr
	}
	return hrs, nil
}
func (o *InstallOptions) RunInstall(f cmdutil.Factory, cmd *cobra.Command) error {
	restCfg, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	c, err := client.New(restCfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	restGetter := kube.NewInClusterRESTClientGetter(restCfg, ns)
	if o.ClusterName != "" {
		byt, err := kubeconfig.FromSecret(cmd.Context(), c, util.ObjectKeyFromString(o.ClusterName))
		if err != nil {
			return err
		}
		restGetter = kube.NewMemoryRESTClientGetter(byt, ns)
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
	cfg := make(map[string]interface{})
	byt, err := os.ReadFile(o.ConfigPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(byt, &cfg)
	if err != nil {
		return err
	}
	hrs, err := o.installCore(cmd.Context(), c, restGetter, cfg, "cert-manager", "cluster-api", "undistro", "ingress-nginx", "undistro-aws")
	if err != nil {
		return err
	}
	fmt.Fprint(o.IOStreams.Out, "Waiting all components to be ready")
	ctx, cancel := context.WithDeadline(cmd.Context(), time.Now().Add(10*time.Minute))
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Fprint(o.ErrOut, "\n\n")
			return ctx.Err()
		default:
			ready := true
			runner, err := helm.NewRunner(restGetter, ns, log.Log)
			if err != nil {
				return err
			}
			for _, hr := range hrs {
				_, err = util.CreateOrUpdate(ctx, c, &hr)
				if err != nil {
					ready = false
					continue
				}
				rel, err := runner.Status(hr)
				if err != nil {
					return err
				}
				if rel.Info.Status != release.StatusDeployed {
					ready = false
					continue
				}
			}
			if ready {
				fmt.Fprintln(o.IOStreams.Out, "\n\nManagement cluster is ready to use.")
				return nil
			}
			fmt.Fprint(o.IOStreams.Out, ".")
			<-time.After(15 * time.Second)
		}
	}

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
			err := o.RunInstall(cmdutil.NewFactory(f), cmd)
			if err != nil {
				fmt.Fprintf(o.ErrOut, "error: %v\n", err)
			}
		},
	}
	return cmd
}
