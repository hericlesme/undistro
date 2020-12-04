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
package util

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"math"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// MergeMaps merges map b into given map a and returns the result.
// It allows overwrites of map values with flat values, and vice versa.
// This is copied from https://github.com/helm/helm/blob/v3.3.0/pkg/cli/values/options.go#L88,
// as the public chartutil.CoalesceTables function does not allow
// overwriting maps with flat values.
func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// ValuesChecksum calculates and returns the SHA1 checksum for the
// given chartutil.Values.
func ValuesChecksum(values chartutil.Values) string {
	var s string
	if len(values) != 0 {
		s, _ = values.YAML()
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// ReleaseRevision returns the revision of the given release.Release.
func ReleaseRevision(rel *release.Release) int {
	if rel == nil {
		return 0
	}
	return rel.Version
}

func CreateOrUpdate(ctx context.Context, r client.Client, o client.Object) (bool, error) {
	old := unstructured.Unstructured{}
	old.SetGroupVersionKind(o.GetObjectKind().GroupVersionKind())
	nm := client.ObjectKey{
		Name:      o.GetName(),
		Namespace: o.GetNamespace(),
	}
	err := r.Get(ctx, nm, &old)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return false, err
		}
		err = r.Create(ctx, o)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	o.SetResourceVersion(old.GetResourceVersion())
	merge := client.MergeFrom(&old)
	byt, err := merge.Data(o)
	if err != nil {
		return false, err
	}
	err = r.Patch(ctx, o, merge)
	if err != nil {
		return false, err
	}
	return len(byt) > 0, nil
}

// ToUnstructured takes a YAML and converts it to a list of Unstructured objects
func ToUnstructured(rawyaml []byte) ([]unstructured.Unstructured, error) {
	var ret []unstructured.Unstructured

	reader := apiyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(rawyaml)))
	count := 1
	for {
		// Read one YAML document at a time, until io.EOF is returned
		b, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrapf(err, "failed to read yaml")
		}
		if len(b) == 0 {
			break
		}

		var m map[string]interface{}
		if err := yaml.Unmarshal(b, &m); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal the %s yaml document: %q", Ordinalize(count), string(b))
		}

		var u unstructured.Unstructured
		u.SetUnstructuredContent(m)

		// Ignore empty objects.
		// Empty objects are generated if there are weird things in manifest files like e.g. two --- in a row without a yaml doc in the middle
		if u.Object == nil {
			continue
		}

		ret = append(ret, u)
		count++
	}

	return ret, nil
}

// Ordinalize takes an int and returns the ordinalized version of it.
// Eg. 1 --> 1st, 103 --> 103rd
func Ordinalize(n int) string {
	m := map[int]string{
		0: "th",
		1: "st",
		2: "nd",
		3: "rd",
		4: "th",
		5: "th",
		6: "th",
		7: "th",
		8: "th",
		9: "th",
	}

	an := int(math.Abs(float64(n)))
	if an < 10 {
		return fmt.Sprintf("%d%s", n, m[an])
	}
	return fmt.Sprintf("%d%s", n, m[an%10])
}
