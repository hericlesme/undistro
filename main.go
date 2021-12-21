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

package main

import (
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/getupio-undistro/record"
	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	appcontrollers "github.com/getupio-undistro/undistro/controllers/app"
	metadatacontrollers "github.com/getupio-undistro/undistro/controllers/metadata"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/version"
	// +kubebuilder:scaffold:imports
)

var (
	setupLog        = ctrl.Log.WithName("setup")
	requestAudience = undistro.GetRequestAudience()
)

func main() {
	var metricsAddr string
	var undistroApiAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&undistroApiAddr, "undistro-api-addr", ":2020", "The address and port of the UnDistro API server")
	klog.InitFlags(nil)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	cfg := ctrl.GetConfigOrDie()

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                 scheme.Scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "19acc88f.undistro.io",
		Namespace:              "",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()

	record.InitFromRecorder(mgr.GetEventRecorderFor("undistro-controller"))

	if err = (&appcontrollers.ClusterReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Cluster")
		os.Exit(1)
	}
	if err = (&appcontrollers.HelmReleaseReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HelmRelease")
		os.Exit(1)
	}
	if err = (&appcontrollers.DefaultPoliciesReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DefaultPolicies")
		os.Exit(1)
	}
	if err = (&metadatacontrollers.ProviderReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Provider")
		os.Exit(1)
	}
	if err = (&appcontrollers.IdentityReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Audience: requestAudience,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Identity")
		os.Exit(1)
	}
	if err = (&appcontrollers.ObserverReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Observer")
		os.Exit(1)
	}

	if err = (&appv1alpha1.Cluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Cluster")
		os.Exit(1)
	}
	if err = (&appv1alpha1.HelmRelease{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HelmRelease")
		os.Exit(1)
	}

	if err = (&appv1alpha1.DefaultPolicies{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "DefaultPolicies")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
	setupLog.Info("Starting Manager", "version", version.Get())
	err = mgr.Start(ctx)
	if err != nil {
		setupLog.Error(err, "failed to start UnDistro")
		os.Exit(1)
	}
}
