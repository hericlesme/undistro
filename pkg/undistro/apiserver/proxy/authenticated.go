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
package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/caarlos0/httperr"
	"github.com/getupio-undistro/undistro/pkg/https"
	"github.com/pkg/errors"
)

func Authenticated(w http.ResponseWriter, r *http.Request) error {
	type authInfo struct {
		Cert     string `json:"cert,omitempty"`
		Key      string `json:"key,omitempty"`
		CA       string `json:"ca,omitempty"`
		Endpoint string `json:"endpoint,omitempty"`
	}
	defer r.Body.Close()
	info := authInfo{}
	b64 := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if b64 == "" {
		return httperr.Wrap(errors.New("invalid auth"), http.StatusBadRequest)
	}
	byt, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return httperr.Wrap(err, http.StatusBadRequest)
	}
	err = json.Unmarshal(byt, &info)
	if err != nil {
		return httperr.Wrap(err, http.StatusBadRequest)
	}
	c, err := https.NewClient(info.Cert, info.Key, info.CA)
	if err != nil {
		return httperr.Wrap(err, http.StatusInternalServerError)
	}
	u := fmt.Sprintf("%s%s", info.Endpoint, r.URL.Path)
	buff := &bytes.Buffer{}
	_, err = io.Copy(buff, r.Body)
	if err != nil && err != io.EOF {
		return httperr.Wrap(err, http.StatusInternalServerError)
	}
	k8sReq, err := http.NewRequestWithContext(r.Context(), r.Method, u, buff)
	if err != nil {
		return httperr.Wrap(err, http.StatusInternalServerError)
	}
	resp, err := c.Do(k8sReq)
	if err != nil {
		return httperr.Wrap(err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil && err != io.EOF {
		return httperr.Wrap(err, http.StatusInternalServerError)
	}
	w.WriteHeader(resp.StatusCode)
	return nil
}
