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

package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/hr"
	"github.com/getupio-undistro/undistro/pkg/kube"
	"github.com/getupio-undistro/undistro/pkg/meta"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/pkg/util"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	conciergev1aplha1 "go.pinniped.dev/generated/latest/apis/concierge/authentication/v1alpha1"
	supervisorconfigv1aplha1 "go.pinniped.dev/generated/latest/apis/supervisor/config/v1alpha1"
	supervisoridpv1aplha1 "go.pinniped.dev/generated/latest/apis/supervisor/idp/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	identityRequeueAfter  = time.Minute * 5
	identityManager       = "pinniped"
	conciergeReleaseName  = "pinniped-concierge"
	supervisorReleaseName = "pinniped-supervisor"
)

type PinnipedComponent string

const (
	concierge  PinnipedComponent = "concierge"
	supervisor PinnipedComponent = "supervisor"
)

// IdentityReconciler reconciles a Identity object
type IdentityReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Namespace string
	Log       logr.Logger
	Audience  string
}

// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *IdentityReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	// Fetch the Identity instance.
	instance := &appv1alpha1.Identity{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(instance, r.Client)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to init patch helper")
	}
	defer func() {
		var patchOpts []patch.Option
		if err == nil {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}
		patchErr := patchHelper.Patch(ctx, instance, patchOpts...)
		if patchErr != nil {
			err = kerrors.NewAggregate([]error{patchErr, err})
			r.Log.Info("failed to Patch identity")
		}
	}()
	r.Log.Info("Checking object age")
	if instance.Generation < instance.Status.ObservedGeneration {
		r.Log.Info("Skipping this old version of reconciled object")
		return ctrl.Result{}, nil
	}
	r.Log.Info("Checking paused")
	if instance.Spec.Paused {
		r.Log.Info("Reconciliation is paused for this object")
		instance = appv1alpha1.IdentityPaused(*instance)
		return ctrl.Result{}, nil
	}
	// Add our finalizer if it does not exist
	if !controllerutil.ContainsFinalizer(instance, meta.Finalizer) {
		controllerutil.AddFinalizer(instance, meta.Finalizer)
		return ctrl.Result{}, nil
	}

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, *instance)
	}
	r.Log.Info("Checking if the Pinniped components are installed")
	result, err := r.reconcile(ctx, *instance)

	durationMsg := fmt.Sprintf("Reconcilation finished in %s", time.Since(start).String())
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	r.Log.Info(durationMsg)
	return result, err
}

// reconcile ensures that, if identity is enabled, pinniped is installed in clusters
func (r *IdentityReconciler) reconcile(ctx context.Context, instance appv1alpha1.Identity) (ctrl.Result, error) {
	cl := &appv1alpha1.Cluster{}
	clusterClient := r.Client
	key := client.ObjectKey{
		Name:      instance.Spec.ClusterName,
		Namespace: instance.GetNamespace(),
	}
	err := r.Get(ctx, key, cl)
	if client.IgnoreNotFound(err) != nil {
		r.Log.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	if util.IsMgmtCluster(instance.Spec.ClusterName) {
		cl.Name = "management"
		cl.Namespace = undistro.Namespace
	}
	values := map[string]interface{}{
		"metadata": map[string]interface{}{
			"namespace": undistro.Namespace,
		},
	}
	err = r.reconcileComponentInstallation(ctx, cl, instance, concierge, undistro.Namespace, "0.10.0", values)
	if err != nil {
		r.Log.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	fedo := make(map[string]interface{})
	o, err := util.GetFromConfigMap(
		ctx, clusterClient, "identity-config", undistro.Namespace, "federationdomain.yaml", fedo)
	fedo = o.(map[string]interface{})
	if err != nil {
		r.Log.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	issuer := fedo["issuer"].(string)
	if util.IsMgmtCluster(instance.Spec.ClusterName) {
		r.Log.Info("Installing Pinniped components in cluster ", "cluster-name", instance.Spec.ClusterName)
		// regex to get ip or dns names
		callbackURL := fmt.Sprintf("https://%s/callback", hostFromURL(issuer))
		values["config"] = map[string]interface{}{
			"callbackURL": callbackURL,
		}
		err = r.reconcileComponentInstallation(ctx, cl, instance, supervisor, undistro.Namespace, "0.10.0-undistro", values)
		if err != nil {
			r.Log.Error(err, err.Error())
			return ctrl.Result{}, err
		}
		err = r.reconcileFederationDomain(ctx, fedo)
		if err != nil {
			r.Log.Error(err, err.Error())
			return ctrl.Result{}, err
		}
		err = r.reconcileOIDCProvider(ctx)
		if err != nil {
			r.Log.Error(err, err.Error())
			return ctrl.Result{}, err
		}
	} else {
		clusterClient, err = kube.NewClusterClient(ctx, r.Client, instance.Spec.ClusterName, cl.GetNamespace())
		if err != nil {
			r.Log.Error(err, err.Error())
			return ctrl.Result{}, err
		}
	}
	err = r.reconcileJWTAuthenticator(ctx, clusterClient, issuer)
	if err != nil {
		r.Log.Error(err, err.Error())
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: identityRequeueAfter}, err
}

func (r *IdentityReconciler) reconcileDelete(ctx context.Context, instance appv1alpha1.Identity) (ctrl.Result, error) {
	// Todo check if is required delete resources created in undistro chart installation
	res, err := hr.Uninstall(ctx, r.Client, r.Log, conciergeReleaseName, instance.Spec.ClusterName, instance.GetNamespace())
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	res, err = hr.Uninstall(ctx, r.Client, r.Log, supervisorReleaseName, instance.Spec.ClusterName, instance.GetNamespace())
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	return res, nil
}

func hostFromURL(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return ""
	}
	return u.Host
}

func (r *IdentityReconciler) reconcileFederationDomain(ctx context.Context, federationDomainCfg map[string]interface{}) error {
	r.Log.Info("Reconciling Federation Domain")

	spec := supervisorconfigv1aplha1.FederationDomainSpec{}
	spec.Issuer = federationDomainCfg["issuer"].(string)
	localClus, err := util.IsLocalCluster(ctx, r.Client)
	if err != nil {
		return err
	}
	if localClus != util.NonLocal {
		spec.TLS = &supervisorconfigv1aplha1.FederationDomainTLSSpec{
			SecretName: federationDomainCfg["tlsSecretName"].(string),
		}
	}
	fedo := &supervisorconfigv1aplha1.FederationDomain{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FederationDomain",
			APIVersion: supervisorconfigv1aplha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "undistro-federationdomain",
			Namespace: undistro.Namespace,
		},
		Spec: spec,
	}
	_, err = util.CreateOrUpdate(ctx, r.Client, fedo)
	if err != nil {
		return err
	}
	return nil
}

func (r *IdentityReconciler) getIdentityConfigMap(ctx context.Context) (map[string]interface{}, error) {
	// get oidc related configmap
	tmp := make(map[string]interface{})
	o, err := util.GetFromConfigMap(
		ctx, r.Client, "identity-config", undistro.Namespace, "oidcprovider.yaml", tmp)
	if err != nil {
		return nil, err
	}
	tmp = o.(map[string]interface{})
	return tmp, nil
}

var providersOIDCProviderCfg = map[string]supervisoridpv1aplha1.OIDCIdentityProviderSpec{
	string(appv1alpha1.Google): {
		Issuer: "https://accounts.google.com",
		AuthorizationConfig: supervisoridpv1aplha1.OIDCAuthorizationConfig{
			AdditionalScopes: []string{"email", "profile"},
		},
		Claims: supervisoridpv1aplha1.OIDCClaims{
			Username: "email",
		},
	},
	string(appv1alpha1.Gitlab): {
		Issuer: "https://gitlab.com",
		Claims: supervisoridpv1aplha1.OIDCClaims{
			Username: "nickname",
			Groups:   "groups",
		},
	},
}

func (r *IdentityReconciler) reconcileOIDCProvider(ctx context.Context) error {
	r.Log.Info("Reconciling OIDC provider")
	cfgMap, err := r.getIdentityConfigMap(ctx)
	if err != nil {
		r.Log.Info(err.Error())
		return err
	}
	name := cfgMap["issuer"].(map[string]interface{})["name"].(string)
	fmtName := fmt.Sprintf("undistro-%s-idp", name)
	spec := supervisoridpv1aplha1.OIDCIdentityProviderSpec{}
	spec.Client = supervisoridpv1aplha1.OIDCClient{
		SecretName: "idp-credentials",
	}
	spec.Issuer = providersOIDCProviderCfg[name].Issuer
	spec.AuthorizationConfig = providersOIDCProviderCfg[name].AuthorizationConfig
	spec.Claims = providersOIDCProviderCfg[name].Claims
	oidcProvider := &supervisoridpv1aplha1.OIDCIdentityProvider{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OIDCIdentityProvider",
			APIVersion: supervisoridpv1aplha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmtName,
			Namespace: undistro.Namespace,
		},
		Spec: spec,
	}
	_, err = util.CreateOrUpdate(ctx, r.Client, oidcProvider)
	if err != nil {
		r.Log.Info(err.Error())
		return err
	}
	return nil
}

func (r *IdentityReconciler) reconcileJWTAuthenticator(ctx context.Context, c client.Client, issuer string) (err error) {
	local, err := util.IsLocalCluster(ctx, c)
	if err != nil {
		return
	}
	const caSecretName = "ca-secret"
	const caName = "ca.crt"
	spec := conciergev1aplha1.JWTAuthenticatorSpec{
		Issuer:   issuer,
		Audience: r.Audience,
	}
	jwtAuth := conciergev1aplha1.JWTAuthenticator{
		TypeMeta: metav1.TypeMeta{
			Kind:       "JWTAuthenticator",
			APIVersion: conciergev1aplha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "supervisor-jwt-authenticator",
			Namespace: undistro.Namespace,
		},
		Spec: spec,
	}
	if local != util.NonLocal {
		secretByt, err := util.GetCaFromSecret(ctx, c, caSecretName, caName, undistro.Namespace)
		if err != nil {
			return err
		}
		secretData := base64.StdEncoding.EncodeToString(secretByt)
		jwtAuth.Spec.TLS = &conciergev1aplha1.TLSSpec{
			CertificateAuthorityData: secretData,
		}
	}
	_, err = util.CreateOrUpdate(ctx, c, &jwtAuth)
	if err != nil {
		return
	}
	return
}

func (r *IdentityReconciler) reconcileComponentInstallation(
	ctx context.Context,
	cl *appv1alpha1.Cluster,
	i appv1alpha1.Identity,
	pc PinnipedComponent,
	targetNs, version string,
	values map[string]interface{},
) (err error) {
	releaseName := fmt.Sprintf("%s-%s", "pinniped", pc)
	release := appv1alpha1.HelmRelease{}
	msg := fmt.Sprintf("Checking if %s is installed", pc)
	r.Log.Info(msg)
	key := client.ObjectKey{
		Name:      hr.GetObjectName(releaseName, i.Spec.ClusterName),
		Namespace: i.GetNamespace(),
	}
	err = r.Get(ctx, key, &release)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}
	release, err = hr.Prepare(releaseName, targetNs, cl.GetNamespace(), version, i.Spec.ClusterName, values)
	if err != nil {
		return err
	}
	if release.Labels == nil {
		release.Labels = make(map[string]string)
	}
	release.Labels[meta.LabelUndistroMove] = ""
	msg = fmt.Sprintf("Installing %s component: %s", identityManager, pc)
	r.Log.Info(msg)
	if err := hr.Install(ctx, r.Client, r.Log, release, cl); err != nil {
		return err
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *IdentityReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		For(&appv1alpha1.Identity{}).
		Complete(r)
}
