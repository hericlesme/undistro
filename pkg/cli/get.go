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
	"github.com/getupio-undistro/undistro/pkg/cloud/aws"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeconfigOptions struct {
	genericclioptions.IOStreams
	Namespace   string
	ClusterName string
}

func NewKubeconfigOptions(streams genericclioptions.IOStreams) *KubeconfigOptions {
	return &KubeconfigOptions{
		IOStreams: streams,
	}
}

func (o *KubeconfigOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.Namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	if len(args) != 1 {
		return errors.New("required 1 argument")
	}
	o.ClusterName = args[0]
	return nil
}

func (o *KubeconfigOptions) RunGetKubeconfig(f cmdutil.Factory, cmd *cobra.Command) error {
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return errors.Errorf("unable to get kubeconfig: %v", err)
	}
	c, err := client.New(cfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return errors.Errorf("unable to create client: %v", err)
	}
	byt, err := kube.GetKubeconfig(cmd.Context(), c, client.ObjectKey{
		Namespace: o.Namespace,
		Name:      o.ClusterName,
	})
	if err != nil {
		return errors.Errorf("unable to get kubeconfig: %v", err)
	}
	_, err = o.IOStreams.Out.Write(byt)
	return err
}

func NewCmdKubeconfig(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewKubeconfigOptions(streams)
	cmd := &cobra.Command{
		Use:                   "kubeconfig [cluster name]",
		DisableFlagsInUseLine: true,
		Short:                 "Get kubeconfig of a cluster",
		Long:                  LongDesc(`Get kubeconfig of a cluster created or imported by UnDistro.`),
		Example: Examples(`
		# Get kubeconfig of a cluster in default namespace
		undistro get kubeconfig cool-cluster
		# Get kubeconfig of a cluster in others namespace
		undistro get kubeconfig cool-cluster -n cool-namespace
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.RunGetKubeconfig(f, cmd))
		},
	}
	return cmd
}

type EksTokenOptions struct {
	genericclioptions.IOStreams
	Namespace   string
	ClusterName string
}

func NewEksTokenOptions(streams genericclioptions.IOStreams) *EksTokenOptions {
	return &EksTokenOptions{
		IOStreams: streams,
	}
}

func (o *EksTokenOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.Namespace, _, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	if len(args) != 1 {
		return errors.New("required 1 argument")
	}
	o.ClusterName = args[0]
	return nil
}

func (o *EksTokenOptions) RunGetKubeconfig(f cmdutil.Factory, cmd *cobra.Command) error {
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return errors.Errorf("unable to get kubeconfig: %v", err)
	}
	c, err := client.New(cfg, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return errors.Errorf("unable to create client: %v", err)
	}
	out, err := aws.GenerateEksToken(cmd.Context(), c, o.ClusterName, o.Namespace)
	if err != nil {
		return errors.Errorf("unable to generate token: %v", err)
	}
	_, err = o.IOStreams.Out.Write([]byte(out))
	return err
}

func NewCmdEksToken(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewKubeconfigOptions(streams)
	cmd := &cobra.Command{
		Use:                   "eks-token [cluster name]",
		DisableFlagsInUseLine: true,
		Short:                 "Get AWS EKS token of a cluster",
		Long:                  LongDesc(`Get AWS EKS token of a cluster created or imported by UnDistro.`),
		Example: Examples(`
		# Get AWS EKS token of a cluster in default namespace
		undistro get eks-token cool-cluster
		# Get AWS EKS token of a cluster in others namespace
		undistro get eks-token cool-cluster -n cool-namespace
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.RunGetKubeconfig(f, cmd))
		},
	}
	return cmd
}

func NewCmdGet(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := get.NewCmdGet("undistro", f, streams)
	cmd.AddCommand(NewCmdKubeconfig(f, streams))
	cmd.AddCommand(NewCmdEksToken(f, streams))
	return cmd
}
