/*
Copyright 2017 The Kubernetes Authors.

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

package patch

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Patch Helper", func() {

	It("Should patch an unstructured object", func() {
		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       "BootstrapMachine",
				"apiVersion": "bootstrap.cluster.x-k8s.io/v1alpha3",
				"metadata": map[string]interface{}{
					"generateName": "test-bootstrap-",
					"namespace":    "default",
				},
			},
		}

		Context("adding an owner reference, preserving its status", func() {
			obj := obj.DeepCopy()

			By("Creating the unstructured object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.GetName(), Namespace: obj.GetNamespace()}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()
			obj.Object["status"] = map[string]interface{}{
				"ready": true,
			}
			Expect(testEnv.Status().Update(ctx, obj)).To(Succeed())

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Modifying the OwnerReferences")
			refs := []metav1.OwnerReference{
				{
					APIVersion: "cluster.x-k8s.io/v1alpha3",
					Kind:       "Cluster",
					Name:       "test",
					UID:        types.UID("fake-uid"),
				},
			}
			obj.SetOwnerReferences(refs)

			By("Patching the unstructured object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating that the status has been preserved")
			ready, err := external.IsReady(obj)
			Expect(err).ToNot(HaveOccurred())
			Expect(ready).To(BeTrue())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}

				return reflect.DeepEqual(obj.GetOwnerReferences(), objAfter.GetOwnerReferences())
			}, timeout).Should(BeTrue())
		})
	})

	Describe("Should patch a clusterv1.Cluster", func() {
		obj := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    "default",
			},
		}

		Specify("add a finalizers", func() {
			obj := obj.DeepCopy()

			By("Creating the object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Adding a finalizer")
			obj.Finalizers = append(obj.Finalizers, clusterv1.ClusterFinalizer)

			By("Patching the object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}

				return reflect.DeepEqual(obj.Finalizers, objAfter.Finalizers)
			}, timeout).Should(BeTrue())
		})

		Specify("removing finalizers", func() {
			obj := obj.DeepCopy()
			obj.Finalizers = append(obj.Finalizers, clusterv1.ClusterFinalizer)

			By("Creating the object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Removing the finalizers")
			obj.SetFinalizers(nil)

			By("Patching the object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}

				return len(objAfter.Finalizers) == 0
			}, timeout).Should(BeTrue())
		})

		Specify("updating spec", func() {
			obj := obj.DeepCopy()
			obj.ObjectMeta.Namespace = "default"

			By("Creating the object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Updating the object spec")
			obj.Spec.Paused = true
			obj.Spec.InfrastructureRef = &corev1.ObjectReference{
				Kind:      "test-kind",
				Name:      "test-ref",
				Namespace: "default",
			}

			By("Patching the object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}

				return objAfter.Spec.Paused == true &&
					reflect.DeepEqual(obj.Spec.InfrastructureRef, objAfter.Spec.InfrastructureRef)
			}, timeout).Should(BeTrue())
		})

		Specify("updating status", func() {
			obj := obj.DeepCopy()

			By("Creating the object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Updating the object status")
			obj.Status.InfrastructureReady = true

			By("Patching the object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}
				return reflect.DeepEqual(objAfter.Status, obj.Status)
			}, timeout).Should(BeTrue())
		})

		Specify("updating both spec, status, and adding a condition", func() {
			obj := obj.DeepCopy()
			obj.ObjectMeta.Namespace = "default"

			By("Creating the object")
			Expect(testEnv.Create(ctx, obj)).ToNot(HaveOccurred())
			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
			defer func() {
				Expect(testEnv.Delete(ctx, obj)).To(Succeed())
			}()

			By("Creating a new patch helper")
			patcher, err := NewHelper(obj, testEnv)
			Expect(err).NotTo(HaveOccurred())

			By("Updating the object spec")
			obj.Spec.Paused = true
			obj.Spec.InfrastructureRef = &corev1.ObjectReference{
				Kind:      "test-kind",
				Name:      "test-ref",
				Namespace: "default",
			}

			By("Updating the object status")
			obj.Status.InfrastructureReady = true

			By("Setting Ready condition")
			conditions.MarkTrue(obj, clusterv1.ReadyCondition)

			By("Patching the object")
			Expect(patcher.Patch(ctx, obj)).To(Succeed())

			By("Validating the object has been updated")
			Eventually(func() bool {
				objAfter := obj.DeepCopy()
				if err := testEnv.Get(ctx, key, objAfter); err != nil {
					return false
				}

				return obj.Status.InfrastructureReady == objAfter.Status.InfrastructureReady &&
					conditions.IsTrue(objAfter, clusterv1.ReadyCondition) &&
					reflect.DeepEqual(obj.Spec, objAfter.Spec)
			}, timeout).Should(BeTrue())
		})
	})
})

func TestNewHelperNil(t *testing.T) {
	var x *appsv1.Deployment
	g := NewWithT(t)
	_, err := NewHelper(x, nil)
	g.Expect(err).ToNot(BeNil())
	_, err = NewHelper(nil, nil)
	g.Expect(err).ToNot(BeNil())
}
