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
	"os"

	"github.com/getupio-undistro/undistro/pkg/graph"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MoveOptions struct {
	ConfigPath  string
	ClusterName string
	Namespace   string
	genericclioptions.IOStreams
}

func NewMoveOptions(streams genericclioptions.IOStreams) *MoveOptions {
	return &MoveOptions{
		IOStreams: streams,
	}
}

func (o *MoveOptions) Complete(f *ConfigFlags, cmd *cobra.Command, args []string) error {
	o.ConfigPath = *f.ConfigFile
	var err error
	o.Namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
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

func (o *MoveOptions) Validate() error {
	if o.ConfigPath != "" {
		_, err := os.Stat(o.ConfigPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *MoveOptions) RunMove(f cmdutil.Factory, cmd *cobra.Command) error {
	cfg := Config{}
	if o.ConfigPath != "" {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			return errors.Errorf("unable to unmarshal config: %v", err)
		}
	}
	key := client.ObjectKey{
		Namespace: o.Namespace,
		Name:      o.ClusterName,
	}
	iopts := &InstallOptions{
		ConfigPath:  o.ConfigPath,
		ClusterName: key.String(),
		IOStreams:   o.IOStreams,
	}
	err := iopts.RunInstall(f, cmd)
	if err != nil {
		return err
	}
	localCfg, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	localClient, err := client.New(localCfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	byt, err := kubeconfig.FromSecret(cmd.Context(), localClient, key)
	if err != nil {
		return err
	}
	restGetter := kube.NewMemoryRESTClientGetter(byt, "")
	remoteCfg, err := restGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	remoteClient, err := client.New(remoteCfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	_ = remoteClient
	objectGraph := graph.NewObjectGraph(localClient)
	// Gets all the types defines by the CRDs installed by clusterctl plus the ConfigMap/Secret core types.
	err = objectGraph.GetDiscoveryTypes()
	if err != nil {
		return err
	}
	if err := objectGraph.Discovery(""); err != nil {
		return err
	}
	// Check whether nodes are not included in GVK considered for move
	objectGraph.CheckVirtualNode()
	return nil
}

func NewCmdMove(f *ConfigFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewMoveOptions(streams)
	cmd := &cobra.Command{
		Use:                   "move [cluster name]",
		DisableFlagsInUseLine: true,
		Short:                 "Move UnDistro resources to another cluster",
		Long: LongDesc(`Install UnDistro.
		IMove UnDistro resources to cluster passed as argument`),
		Example: Examples(`
		# Move
		undistro --config undistro-config.yaml move cool-product-cluster -n undistro-production
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.RunMove(cmdutil.NewFactory(f), cmd))
		},
	}
	return cmd
}
