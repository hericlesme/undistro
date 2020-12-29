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
	"fmt"

	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
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
	byt, err := getKubeconfig(cmd.Context(), c, client.ObjectKey{
		Namespace: o.Namespace,
		Name:      o.ClusterName,
	})
	if err != nil {
		return errors.Errorf("unable to get kubeconfig: %v", err)
	}
	_, err = o.IOStreams.Out.Write(byt)
	return err
}

// Purpose is the name to append to the secret generated for a cluster.
type Purpose string

const (
	// KubeconfigDataName is the key used to store a Kubeconfig in the secret's data field.
	KubeconfigDataName = "value"
	// Kubeconfig is the secret name suffix storing the Cluster Kubeconfig.
	Kubeconfig = Purpose("kubeconfig")
	// UserKubeconfig is the secret name suffix storing the Cluster Kubeconfig for user usage.
	UserKubeconfig = Purpose("user-kubeconfig")
)

func getKubeconfig(ctx context.Context, c client.Reader, cluster client.ObjectKey) ([]byte, error) {
	out, err := getSecret(ctx, c, cluster, UserKubeconfig)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		out, err = getSecret(ctx, c, cluster, Kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return toKubeconfigBytes(out)
}

// Get retrieves the specified Secret (if any) from the given
// cluster name and namespace.
func getSecret(ctx context.Context, c client.Reader, cluster client.ObjectKey, purpose Purpose) (*corev1.Secret, error) {
	return getFromNamespacedName(ctx, c, cluster, purpose)
}

// GetFromNamespacedName retrieves the specified Secret (if any) from the given
// cluster name and namespace.
func getFromNamespacedName(ctx context.Context, c client.Reader, clusterName client.ObjectKey, purpose Purpose) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: clusterName.Namespace,
		Name:      name(clusterName.Name, purpose),
	}

	if err := c.Get(ctx, secretKey, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

// Name returns the name of the secret for a cluster.
func name(cluster string, suffix Purpose) string {
	return fmt.Sprintf("%s-%s", cluster, suffix)
}

func toKubeconfigBytes(out *corev1.Secret) ([]byte, error) {
	data, ok := out.Data[KubeconfigDataName]
	if !ok {
		return nil, errors.Errorf("missing key %q in secret data", KubeconfigDataName)
	}
	return data, nil
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

func NewCmdGet(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := get.NewCmdGet("undistro", f, streams)
	cmd.AddCommand(NewCmdKubeconfig(f, streams))
	return cmd
}
