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
	"fmt"
	"io/ioutil"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/retry"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kind/pkg/cluster"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
)

type SetupOptions struct {
	genericclioptions.IOStreams
	Provider          string
	Name              string
	ConfigPath        string
	Flavor            string
	SSHKeyName        string
	CloudsFile        string
	KubernetesVersion string
	Region            string
	rawCfg            *ConfigFlags
}

func NewSetupOptions(streams genericclioptions.IOStreams) *SetupOptions {
	return &SetupOptions{
		IOStreams: streams,
		Name:      "undistro",
	}
}

func (o *SetupOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.Name, "name", o.Name, "name of the cluster (default: undistro)")
	flags.StringVar(&o.Flavor, "flavor", o.Flavor, "flavor used for management cluster")
	flags.StringVar(&o.Region, "region", o.Region, "region used for management cluster")
	flags.StringVar(&o.SSHKeyName, "ssh-key-name", o.SSHKeyName, "ssh key name used to create the management cluster")
	flags.StringVar(&o.KubernetesVersion, "kubernetes-version", o.KubernetesVersion, "Kubernetes version used to create the management cluster")
	flags.StringVar(&o.CloudsFile, "openstack-clouds-file", o.CloudsFile, "path of clouds.yaml (required by provider openstack)")
}

func (o *SetupOptions) Complete(f *ConfigFlags, args []string) error {
	o.rawCfg = f
	o.ConfigPath = *f.ConfigFile
	switch len(args) {
	case 1:
		o.Provider = args[0]
	default:
		return errors.New("required 1 argument")
	}
	if o.Provider == "openstack" {
		if o.CloudsFile == "" {
			return errors.New("clouds file is required for provider openstack")
		}
		o.Flavor = o.Provider
	}
	if o.Name == "" {
		o.Name = undistro.LocalCluster
	}
	return nil
}

func (o *SetupOptions) RunSetup(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(o.IOStreams.Out, "Setup a bootstrap cluster")
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(kindcmd.NewLogger()),
		cluster.ProviderWithDocker(),
	)
	err := provider.Create(
		o.Name,
		cluster.CreateWithRawConfig([]byte(undistro.KindCfg)),
		cluster.CreateWithNodeImage(""),
		cluster.CreateWithRetain(false),
		cluster.CreateWithWaitForReady(time.Duration(0)),
		cluster.CreateWithKubeconfigPath(""),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
	)
	if err != nil {
		return err
	}
	opts := InstallOptions{
		ConfigPath:  o.ConfigPath,
		IOStreams:   o.IOStreams,
		Remote:      false,
		ClusterName: o.Name,
	}
	err = opts.Complete(o.rawCfg, cmd, args)
	if err != nil {
		return err
	}
	fmt.Fprintln(o.IOStreams.Out, "Install UnDistro into the bootstrap cluster")
	factory := cmdutil.NewFactory(o.rawCfg)
	err = opts.RunInstall(factory, cmd)
	if err != nil {
		return err
	}
	if o.Provider != "kind" {
		cfg, err := factory.ToRESTConfig()
		if err != nil {
			return err
		}
		c, err := client.New(cfg, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return err
		}
		createClusterOpts := ClusterOptions{
			IOStreams:    o.IOStreams,
			Namespace:    undistro.Namespace,
			ClusterName:  undistro.MgmtClusterName,
			Infra:        o.Provider,
			Flavor:       o.Flavor,
			SshKeyName:   o.SSHKeyName,
			GenerateFile: false,
			K8sVersion:   o.KubernetesVersion,
			Region:       o.Region,
			AuthEnabled:  false,
			Addons:       false,
			CloudsFile:   o.CloudsFile,
		}
		err = createClusterOpts.Complete(factory, cmd, args)
		if err != nil {
			return err
		}
		err = createClusterOpts.RunCreateCluster(factory, cmd)
		if err != nil {
			return err
		}
		fmt.Fprintf(o.IOStreams.Out, "Creating management cluster provider=%s flavor=%s\n", o.Provider, o.Flavor)
		fmt.Fprintln(o.IOStreams.Out, "Waiting management cluster to be ready")
		key := client.ObjectKey{
			Name:      undistro.MgmtClusterName,
			Namespace: undistro.Namespace,
		}
		for {
			<-time.After(30 * time.Second)
			cl := appv1alpha1.Cluster{}
			err = c.Get(cmd.Context(), key, &cl)
			if err != nil {
				return err
			}
			fmt.Fprintf(o.IOStreams.Out, "Cluster status = %s\n", util.LastCondition(cl.Status.Conditions).Message)
			if meta.InReadyCondition(cl.Status.Conditions) {
				break
			}
		}
		byt, err := kube.GetKubeconfig(cmd.Context(), c, key)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("undistro-management.kubeconfig", byt, 0644)
		if err != nil {
			return err
		}
		moveOpts := MoveOptions{
			ConfigPath:  o.ConfigPath,
			IOStreams:   o.IOStreams,
			ClusterName: undistro.MgmtClusterName,
			Namespace:   undistro.Namespace,
		}
		err = moveOpts.Complete(o.rawCfg, cmd, args)
		if err != nil {
			return err
		}
		return retry.WithExponentialBackoff(retry.NewBackoff(), func() error {
			ierr := moveOpts.RunMove(factory, cmd)
			if ierr != nil {
				fmt.Fprintf(o.IOStreams.ErrOut, "%s, retrying...\n", ierr)
				return ierr
			}
			return nil
		})
	}
	destroyOpts := DestroyOptions{
		IOStreams: o.IOStreams,
		Name:      o.Name,
		Provider:  "kind",
	}
	destroyOpts.RunDestroy(cmd)
	return nil
}

func NewCmdSetup(f *ConfigFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewSetupOptions(streams)
	cmd := &cobra.Command{
		Use:                   "setup [tool]",
		DisableFlagsInUseLine: true,
		Short:                 "Setup a tool",
		Long:                  LongDesc(`Setup a tool`),
		Example: Examples(`
		undistro setup kind --name undistro-cluster
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, args))
			cmdutil.CheckErr(o.RunSetup(cmd, args))
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}
