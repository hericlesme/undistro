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
package e2e_test

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/getupio-undistro/clilib"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	_ "github.com/go-task/slim-sprig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/cluster-api/test/framework/exec"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"sigs.k8s.io/yaml"
)

var (
	undcli    = &clilib.CLI{}
	e2eRun    = flag.Bool("e2e", false, "set true to run e2e tests")
	k8sClient client.Client
)

func TestMain(m *testing.M) {
	flag.Parse()
	fmt.Println("E2E")
	runE2E := *e2eRun
	if !runE2E {
		fmt.Println("Skipping E2E")
		os.Exit(0)
	}
	ctx := context.Background()
	fmt.Println("Build docker image and push")
	sha := os.Getenv("GITHUB_SHA")

	regHost, ok := os.LookupEnv("REG_ADDR")
	if !ok {
		fmt.Println("Environment variable <REG_ADDR> not found, using 'localhost:5000'.")
		regHost = "localhost:5000"
	}
	cmd := exec.NewCommand(
		exec.WithCommand("bash"),
		exec.WithArgs("-c", fmt.Sprintf("../../../testbin/docker-build-e2e.sh %s %s", regHost, sha)),
	)
	stdout, stderr, err := cmd.Run(ctx)
	if err != nil {
		fmt.Println(string(stderr))
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(string(stdout))
	cfg := map[string]interface{}{
		"global": map[string]interface{}{
			"undistroRepository": "localhost:5000",
			"undistroVersion":    sha,
		},
		"undistro-aws": map[string]interface{}{
			"enabled": true,
			"credentials": map[string]interface{}{
				"accessKeyID":     os.Getenv("E2E_AWS_ACCESS_KEY_ID"),
				"secretAccessKey": os.Getenv("E2E_AWS_SECRET_ACCESS_KEY"),
			},
		},
	}
	byt, _ := yaml.Marshal(cfg)
	err = ioutil.WriteFile("undistro-config.yaml", byt, 0700)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	k8sClient, err = client.New(config, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Install UnDistro")
	sout, serr, _ := undcli.Install("--config", "undistro-config.yaml")
	fmt.Println(sout)
	if !strings.Contains(sout, "Management cluster is ready to use.") {
		msg := "failed to install undistro: " + serr
		fmt.Println(msg)
		cmd = exec.NewCommand(
			exec.WithCommand("kubectl"),
			exec.WithArgs("get", "pods", "--all-namespaces"),
		)
		stdout, stderr, _ = cmd.Run(ctx)
		fmt.Println(string(stdout))
		fmt.Println("err:", string(stderr))
		cmd = exec.NewCommand(
			exec.WithCommand("kubectl"),
			exec.WithArgs("describe", "nodes"),
		)
		stdout, stderr, _ = cmd.Run(ctx)
		fmt.Println(string(stdout))
		fmt.Println("err:", string(stderr))
		cmd = exec.NewCommand(
			exec.WithCommand("kubectl"),
			exec.WithArgs("describe", "pods", "-n", "undistro-system"),
		)
		stdout, stderr, _ = cmd.Run(ctx)
		fmt.Println(string(stdout))
		fmt.Println("err:", string(stderr))
		podList := corev1.PodList{}
		err = k8sClient.List(ctx, &podList, client.InNamespace("undistro-system"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, p := range podList.Items {
			if strings.Contains(p.Name, "undistro") {
				sout, serr, _ = undcli.Logs(p.Name, "-n", "undistro-system", "-c", "manager", "--previous")
				fmt.Println(sout)
				fmt.Println("err:", stderr)
			}
		}
		cmd = exec.NewCommand(
			exec.WithCommand("helm"),
			exec.WithArgs("get", "values", "undistro", "-n", "undistro-system"),
		)
		stdout, _, err = cmd.Run(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(stdout))
		cmd = exec.NewCommand(
			exec.WithCommand("helm"),
			exec.WithArgs("ls", "-n", "undistro-system"),
		)
		stdout, _, err = cmd.Run(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(stdout))
		cmd = exec.NewCommand(
			exec.WithCommand("helm"),
			exec.WithArgs("status", "undistro", "--show-desc", "-n", "undistro-system"),
		)
		stdout, _, err = cmd.Run(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(stdout))
		os.Exit(1)
	}

	sout, _, err = undcli.Get("pods", "-n", "undistro-system")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(sout)
	code := m.Run()
	os.Exit(code)
}

func TestGinkgoSuite(t *testing.T) {
	SetDefaultEventuallyPollingInterval(1 * time.Minute)
	SetDefaultEventuallyTimeout(120 * time.Minute)
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"E2E Suite",
		[]Reporter{printer.NewlineReporter{}})
}
