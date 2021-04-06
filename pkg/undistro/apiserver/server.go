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
package apiserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/getupio-undistro/undistro/pkg/fs"
	"github.com/getupio-undistro/undistro/pkg/undistro/apiserver/proxy"
	"github.com/gorilla/mux"
	"gocloud.dev/server"
	"gocloud.dev/server/requestlog"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type Server struct {
	genericclioptions.IOStreams
	*server.Server
	K8sFactory cmdutil.Factory
}

func NewServer(f cmdutil.Factory, in io.Reader, out, errOut io.Writer) *Server {
	streams := genericclioptions.IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}
	opts := &server.Options{
		RequestLogger: requestlog.NewNCSALogger(streams.Out, func(err error) {
			fmt.Fprintln(streams.ErrOut, err.Error())
		}),
	}
	router := mux.NewRouter()
	s := server.New(router, opts)
	deamonServer := &Server{
		IOStreams:  streams,
		Server:     s,
		K8sFactory: f,
	}
	deamonServer.routes(router)
	return deamonServer
}

func (s *Server) routes(router *mux.Router) {
	cfg, err := s.K8sFactory.ToRESTConfig()
	if err != nil {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			klog.Fatal(err)
		}
	}
	frontFS, err := fs.GetFrontendFS()
	if err != nil {
		klog.Fatal(err)
	}
	router.PathPrefix("/v1/namespaces/{namespece}/clusters/{cluster}/proxy").Handler(proxy.NewHandler(cfg))
	router.Handle("/", http.FileServer(http.FS(frontFS)))
}

func (s *Server) GracefullyStart(ctx context.Context, addr string) error {
	cerr := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func(ctx context.Context) {
		klog.Infof("listen on %s", addr)
		cerr <- s.ListenAndServe(addr)
	}(ctx)
	select {
	case <-sigCh:
		return s.Shutdown(ctx)
	case err := <-cerr:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
