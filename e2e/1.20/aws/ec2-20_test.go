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
	"fmt"
	"os"
	"time"

	"github.com/getupio-undistro/meta"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Create EC2 cluster 1.20", func() {
	var (
		clusterClient client.Client
	)
	It("Should generate recommend cluster spec 1.20", func() {
		_, _, err := undcli.Create("cluster", "ec2-20-e2e", "-n", "e2e", "--infra", "aws", "--flavor", "ec2", "--ssh-key-name", "undistro", "--generate-file")
		Expect(err).ToNot(HaveOccurred())
		_, err = os.Stat("ec2-20-e2e.yaml")
		Expect(err).ToNot(HaveOccurred())

	})
	It("Should create EC2 cluster 1.20", func() {
		sout, _, err := undcli.Apply("-f", "../../testdata/ec2-20.yaml")
		fmt.Println(err)
		Expect(err).ToNot(HaveOccurred())
		fmt.Println(sout)
		Eventually(func() bool {
			cl := appv1alpha1.Cluster{}
			key := client.ObjectKey{
				Name:      "ec2-20-e2e",
				Namespace: "e2e",
			}
			err = k8sClient.Get(context.Background(), key, &cl)
			if err != nil {
				fmt.Println(err)
				return false
			}
			fmt.Println(cl)
			return meta.InReadyCondition(cl.Status.Conditions)
		}, 240*time.Minute, 2*time.Minute).Should(BeTrue())
		fmt.Println("Get Kubeconfig")
		sout, _, err = undcli.Get("kubeconfig", "ec2-20-e2e", "-n", "e2e", "--admin")
		Expect(err).ToNot(HaveOccurred())
		getter := kube.NewMemoryRESTClientGetter([]byte(sout), "")
		cfg, err := getter.ToRESTConfig()
		Expect(err).ToNot(HaveOccurred())
		Expect(cfg).ToNot(BeNil())
		clusterClient, err = client.New(cfg, client.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(clusterClient).ToNot(BeNil())
		Eventually(func() []corev1.Node {
			nodes := corev1.NodeList{}
			err = clusterClient.List(context.Background(), &nodes)
			if err != nil {
				return []corev1.Node{}
			}
			fmt.Println(nodes.Items)
			fmt.Println(len(nodes.Items))
			return nodes.Items
		}, 240*time.Minute, 2*time.Minute).Should(HaveLen(7))
		Eventually(func() []corev1.Node {
			cpNodes := make([]corev1.Node, 0)
			nodes := corev1.NodeList{}
			err = clusterClient.List(context.Background(), &nodes)
			if err != nil {
				return []corev1.Node{}
			}
			fmt.Println(nodes.Items)
			fmt.Println(len(nodes.Items))
			for _, n := range nodes.Items {
				labels := n.GetLabels()
				if labels != nil {
					_, okCP := labels[meta.LabelK8sCP]
					_, okMaster := labels[meta.LabelK8sMaster]
					if okCP || okMaster {
						cpNodes = append(cpNodes, n)
					}
				}
			}
			fmt.Println(cpNodes)
			fmt.Println(len(cpNodes))
			return cpNodes
		}, 240*time.Minute, 2*time.Minute).Should(HaveLen(3))
		fmt.Println("check kyverno")
		Eventually(func() []unstructured.Unstructured {
			list := unstructured.UnstructuredList{}
			list.SetGroupVersionKind(schema.FromAPIVersionAndKind("kyverno.io/v1", "ClusterPolicyList"))
			err = clusterClient.List(context.Background(), &list)
			if err != nil {
				fmt.Println(err)
				return []unstructured.Unstructured{}
			}
			fmt.Println(list.Items)
			fmt.Println(len(list.Items))
			if undistroPodName != "" {
				sout, serr, err := undcli.Logs(undistroPodName, "-n", "undistro-system", "-c", "manager")
				if err != nil {
					fmt.Println(serr)
					fmt.Println(err)
					fmt.Println(sout)
					return list.Items
				}
				fmt.Println(sout)
			}
			return list.Items
		}, 240*time.Minute, 2*time.Minute).Should(HaveLen(16))
		podList := corev1.PodList{}
		err = clusterClient.List(context.Background(), &podList, client.InNamespace("kube-system"))
		Expect(err).ToNot(HaveOccurred())
		for _, p := range podList.Items {
			for _, container := range p.Spec.Containers {
				Expect(container.Image).To(HavePrefix("registry.undistro.io"))
			}
		}
		fmt.Println("delete cluster")
		sout, _, err = undcli.Delete("-f", "../../testdata/ec2-20.yaml")
		Expect(err).ToNot(HaveOccurred())
		fmt.Println(sout)
	}, float64(480*time.Minute))
})
