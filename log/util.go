/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package log

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// UnstructuredToValues provide a utility function for creation values describing an Unstructured objects. e.g.
// Deployment="capd-controller-manager" Namespace="capd-system"  (<Kind>=<name> Namespace=<Namespace>)
// CustomResourceDefinition="dockerclusters.infrastructure.cluster.x-k8s.io" (omit Namespace if it does not apply)
func UnstructuredToValues(obj unstructured.Unstructured) []interface{} {
	values := []interface{}{
		obj.GetKind(), obj.GetName(),
	}
	if obj.GetNamespace() != "" {
		values = append(values, "Namespace", obj.GetNamespace())
	}
	return values
}
