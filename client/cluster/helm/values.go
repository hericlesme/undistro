package helm

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	undistrov1 "github.com/getupcloud/undistro/api/v1alpha1"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chartutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// readURL attempts to read a file from an HTTP(S) URL.
func readURL(URL string) ([]byte, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return []byte{}, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return []byte{}, errors.Errorf("URL scheme should be HTTP(S), got '%s'", u.Scheme)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	default:
		return []byte{}, errors.Errorf("failed to retrieve file from URL, status '%s (%d)'", resp.Status, resp.StatusCode)
	}
}

// readLocalChartFile attempts to read a file from the chart path.
func readLocalChartFile(filePath string) ([]byte, error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}
	return f, nil
}

// MmrgeValues merges source and destination map, preferring values
// from the source values. This is slightly adapted from:
// https://github.com/helm/helm/blob/2332b480c9cb70a0d8a85247992d6155fbe82416/cmd/helm/install.go#L359
func mergeValues(dest, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

func ComposeValues(ctx context.Context, r client.Client, hr *undistrov1.HelmRelease, chartPath string) ([]byte, error) {
	var res chartutil.Values = make(map[string]interface{})
	for _, v := range hr.GetValuesFromSources() {
		nm := types.NamespacedName{
			Namespace: hr.Namespace,
		}
		var result chartutil.Values
		switch {
		case v.ConfigMapKeyRef != nil:
			cm := v.ConfigMapKeyRef
			nm.Name = cm.Name
			if cm.Namespace != "" {
				nm.Namespace = cm.Namespace
			}
			key := cm.Key
			if key == "" {
				key = "values.yaml"
			}
			var configMap corev1.ConfigMap
			err := r.Get(ctx, nm, &configMap)
			if err != nil {
				if apierrors.IsNotFound(err) && cm.Optional {
					continue
				}
				return []byte{}, err
			}
			d, ok := configMap.Data[key]
			if !ok {
				if cm.Optional {
					continue
				}
				return []byte{}, errors.Errorf("could not find key %v in ConfigMap %s", key, nm.String())
			}
			result, err = chartutil.ReadValues([]byte(d))
			if err != nil {
				if cm.Optional {
					continue
				}
				return []byte{}, errors.Errorf("unable to yaml.Unmarshal %v from %s in ConfigMap %s", d, key, nm.String())
			}
		case v.SecretKeyRef != nil:
			s := v.SecretKeyRef
			nm.Name = s.Name
			if s.Namespace != "" {
				nm.Namespace = s.Namespace
			}
			key := s.Key
			if key == "" {
				key = "values.yaml"
			}
			var secret corev1.Secret
			err := r.Get(ctx, nm, &secret)
			if err != nil {
				if apierrors.IsNotFound(err) && s.Optional {
					continue
				}
				return []byte{}, err
			}
			d, ok := secret.Data[key]
			if !ok {
				if s.Optional {
					continue
				}
				return []byte{}, errors.Errorf("could not find key %s in Secret %s", key, nm.String())
			}
			result, err = chartutil.ReadValues(d)
			if err != nil {
				return []byte{}, errors.Errorf("unable to yaml.Unmarshal %v from %s in Secret %s", d, key, nm.String())
			}
		case v.ExternalSourceRef != nil:
			es := v.ExternalSourceRef
			u := es.URL
			optional := es.Optional != nil && *es.Optional
			b, err := readURL(u)
			if err != nil {
				if optional {
					continue
				}
				return []byte{}, errors.Errorf("unable to read value file from URL %s", u)
			}
			result, err = chartutil.ReadValues(b)
			if err != nil {
				if optional {
					continue
				}
				return []byte{}, errors.Errorf("unable to yaml.Unmarshal %v from URL %s", b, u)
			}
		case v.ChartFileRef != nil:
			cf := v.ChartFileRef
			filePath := cf.Path
			optional := cf.Optional != nil && *cf.Optional
			f, err := readLocalChartFile(filepath.Join(chartPath, filePath))
			if err != nil {
				if optional {
					continue
				}
				return []byte{}, errors.Errorf("unable to read value file from path %s", filePath)
			}
			result, err = chartutil.ReadValues(f)
			if err != nil {
				if optional {
					continue
				}
				return []byte{}, errors.Errorf("unable to yaml.Unmarshal %v from path %s", f, filePath)
			}
		}
		res = mergeValues(res, result)
	}
	res = mergeValues(res, hr.GetValues())
	y, err := res.YAML()
	if err != nil {
		return []byte{}, err
	}
	return []byte(y), nil
}