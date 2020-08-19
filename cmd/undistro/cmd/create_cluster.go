package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/getupcloud/undistro/client"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	utilresource "sigs.k8s.io/cluster-api/util/resource"
	"sigs.k8s.io/cluster-api/util/yaml"
)

type createClusterOptions struct {
	url        string
	kubeconfig string
}

var ccOpts = &createClusterOptions{}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a cluster",
	Long:  "Create a cluster",

	Example: Examples(`
		# Creates a cluster with variable values using
		a template from a specific URL.
		undistro create cluster --from https://github.com/foo-org/foo-repository/blob/master/cluster-template.yaml

		# Creates a cluster with variable values using
		a template stored locally.
		undistro create cluster --from ~/workspace/cluster-template.yaml

		# Creates a cluster from template passed in via stdin
		cat ~/workspace/cluster-template.yaml | undistro create cluster
`),

	RunE: func(cmd *cobra.Command, args []string) error {
		return createCluster(os.Stdin, os.Stdout)
	},
}

func init() {
	// flags for the url source
	createClusterCmd.Flags().StringVarP(&ccOpts.url, "from", "f", "-",
		"The URL to read the template from. It defaults to '-' which reads from stdin.")
	createClusterCmd.Flags().StringVar(&ccOpts.kubeconfig, "kubeconfig", "",
		"Path to a kubeconfig file to use for the management cluster. If empty, default discovery rules apply.")

	createCmd.AddCommand(createClusterCmd)
}

func createCluster(r io.Reader, w io.Writer) error {
	c, err := client.New(cfgFile)
	if err != nil {
		return err
	}
	options := client.ProcessYAMLOptions{
		ListVariablesOnly: false,
	}

	options.URLSource = &client.URLSourceOptions{
		URL: ccOpts.url,
	}
	if ccOpts.url == "-" {
		options.ReaderSource = &client.ReaderSourceOptions{
			Reader: r,
		}
		options.URLSource = nil
	}
	printer, err := c.ProcessYAML(options)
	if err != nil {
		return err
	}
	b, err := printer.Yaml()
	if err != nil {
		return err
	}
	objs, err := yaml.ToUnstructured(b)
	if err != nil {
		return err
	}
	proxy, err := c.GetProxy()
	if err != nil {
		return err
	}
	k8sClient, err := proxy.NewClient()
	if err != nil {
		return err
	}
	objs = utilresource.SortForCreate(objs)
	nm := types.NamespacedName{}
	for i, o := range objs {
		if i == 0 {
			nm.Name = o.GetName()
			nm.Namespace = o.GetNamespace()
		}
		err = k8sClient.Create(context.Background(), &o)
		if err != nil {
			return err
		}
	}
	logStreamer, err := c.GetLogs(client.Kubeconfig{
		Path: ccOpts.kubeconfig,
	})
	if err != nil {
		return err
	}
	logChan := make(chan string)
	cfg, err := proxy.GetConfig()
	if err != nil {
		return err
	}
	go logStreamer.Stream(context.Background(), cfg, logChan, nm)
	for str := range logChan {
		fmt.Fprint(os.Stdout, str)
	}
	if nm.Namespace == "" {
		nm.Namespace = "default"
	}
	fmt.Fprintf(os.Stdout, "\n\nCluster %s is ready. \nRun undistro get kubeconfig %s -n %s to get the Kubeconfig\n\n", nm.String(), nm.Name, nm.Namespace)
	return nil
}
