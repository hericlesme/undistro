// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package jwtcachefiller implements a controller for filling an authncache.Cache with each
// added/updated JWTAuthenticator.
package jwtcachefiller

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"time"

	coreosoidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-logr/logr"
	"gopkg.in/square/go-jose.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/apiserver/plugin/pkg/authenticator/token/oidc"
	"k8s.io/klog/v2"

	pinnipedcontroller "github.com/getupio-undistro/undistro/third_party/pinniped/internal/controller"
	pinnipedauthenticator "github.com/getupio-undistro/undistro/third_party/pinniped/internal/controller/authenticator"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/controller/authenticator/authncache"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/controllerlib"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/net/phttp"
	auth1alpha1 "go.pinniped.dev/generated/latest/apis/concierge/authentication/v1alpha1"
	authinformers "go.pinniped.dev/generated/latest/client/concierge/informers/externalversions/authentication/v1alpha1"
)

// These default values come from the way that the Supervisor issues and signs tokens. We make these
// the defaults for a JWTAuthenticator so that they can easily integrate with the Supervisor.
const (
	defaultUsernameClaim = "username"
	defaultGroupsClaim   = "groups"
)

// defaultSupportedSigningAlgos returns the default signing algos that this JWTAuthenticator
// supports (i.e., if none are supplied by the user).
func defaultSupportedSigningAlgos() []string {
	return []string{
		// RS256 is recommended by the OIDC spec and required, in some capacity. Since we want the
		// JWTAuthenticator to be able to support many OIDC ID tokens out of the box, we include this
		// algorithm by default.
		string(jose.RS256),
		// ES256 is what the Supervisor does, by default. We want integration with the JWTAuthenticator
		// to be as seamless as possible, so we include this algorithm by default.
		string(jose.ES256),
	}
}

type tokenAuthenticatorCloser interface {
	authenticator.Token
	pinnipedauthenticator.Closer
}

type jwtAuthenticator struct {
	tokenAuthenticatorCloser
	spec *auth1alpha1.JWTAuthenticatorSpec
}

// New instantiates a new controllerlib.Controller which will populate the provided authncache.Cache.
func New(
	cache *authncache.Cache,
	jwtAuthenticators authinformers.JWTAuthenticatorInformer,
	log logr.Logger,
) controllerlib.Controller {
	return controllerlib.New(
		controllerlib.Config{
			Name: "jwtcachefiller-controller",
			Syncer: &controller{
				cache:             cache,
				jwtAuthenticators: jwtAuthenticators,
				log:               log.WithName("jwtcachefiller-controller"),
			},
		},
		controllerlib.WithInformer(
			jwtAuthenticators,
			pinnipedcontroller.MatchAnythingFilter(nil), // nil parent func is fine because each event is distinct
			controllerlib.InformerOption{},
		),
	)
}

type controller struct {
	cache             *authncache.Cache
	jwtAuthenticators authinformers.JWTAuthenticatorInformer
	log               logr.Logger
}

// Sync implements controllerlib.Syncer.
func (c *controller) Sync(ctx controllerlib.Context) error {
	obj, err := c.jwtAuthenticators.Lister().Get(ctx.Key.Name)
	if err != nil && errors.IsNotFound(err) {
		c.log.Info("Sync() found that the JWTAuthenticator does not exist yet or was deleted")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get JWTAuthenticator %s/%s: %w", ctx.Key.Namespace, ctx.Key.Name, err)
	}

	cacheKey := authncache.Key{
		APIGroup: auth1alpha1.GroupName,
		Kind:     "JWTAuthenticator",
		Name:     ctx.Key.Name,
	}

	// If this authenticator already exists, then only recreate it if is different from the desired
	// authenticator. We don't want to be creating a new authenticator for every resync period.
	//
	// If we do need to recreate the authenticator, then make sure we close the old one to avoid
	// goroutine leaks.
	if value := c.cache.Get(cacheKey); value != nil {
		jwtAuthenticator := c.extractValueAsJWTAuthenticator(value)
		if jwtAuthenticator != nil {
			if reflect.DeepEqual(jwtAuthenticator.spec, &obj.Spec) {
				c.log.WithValues("jwtAuthenticator", klog.KObj(obj), "issuer", obj.Spec.Issuer).Info("actual jwt authenticator and desired jwt authenticator are the same")
				return nil
			}
			jwtAuthenticator.Close()
		}
	}

	// Make a deep copy of the spec so we aren't storing pointers to something that the informer cache
	// may mutate!
	jwtAuthenticator, err := newJWTAuthenticator(obj.Spec.DeepCopy())
	if err != nil {
		return fmt.Errorf("failed to build jwt authenticator: %w", err)
	}

	c.cache.Store(cacheKey, jwtAuthenticator)
	c.log.WithValues("jwtAuthenticator", klog.KObj(obj), "issuer", obj.Spec.Issuer).Info("added new jwt authenticator")
	return nil
}

func (c *controller) extractValueAsJWTAuthenticator(value authncache.Value) *jwtAuthenticator {
	jwtAuthenticator, ok := value.(*jwtAuthenticator)
	if !ok {
		actualType := "<nil>"
		if t := reflect.TypeOf(value); t != nil {
			actualType = t.String()
		}
		c.log.WithValues("actualType", actualType).Info("wrong JWT authenticator type in cache")
		return nil
	}
	return jwtAuthenticator
}

// newJWTAuthenticator creates a jwt authenticator from the provided spec.
func newJWTAuthenticator(spec *auth1alpha1.JWTAuthenticatorSpec) (*jwtAuthenticator, error) {
	rootCAs, caBundle, err := pinnipedauthenticator.CABundle(spec.TLS)
	if err != nil {
		return nil, fmt.Errorf("invalid TLS configuration: %w", err)
	}

	var caContentProvider oidc.CAContentProvider
	if len(caBundle) != 0 {
		var caContentProviderErr error
		caContentProvider, caContentProviderErr = dynamiccertificates.NewStaticCAContent("ignored", caBundle)
		if caContentProviderErr != nil {
			return nil, caContentProviderErr // impossible since caBundle is validated already
		}
	}
	usernameClaim := spec.Claims.Username
	if usernameClaim == "" {
		usernameClaim = defaultUsernameClaim
	}
	groupsClaim := spec.Claims.Groups
	if groupsClaim == "" {
		groupsClaim = defaultGroupsClaim
	}

	// copied from Kube OIDC code
	issuerURL, err := url.Parse(spec.Issuer)
	if err != nil {
		return nil, err
	}
	if issuerURL.Scheme != "https" {
		return nil, fmt.Errorf("issuer (%q) has invalid scheme (%q), require 'https'", spec.Issuer, issuerURL.Scheme)
	}

	client := phttp.Default(rootCAs)
	client.Timeout = 30 * time.Second // copied from Kube OIDC code

	ctx := coreosoidc.ClientContext(context.Background(), client)

	provider, err := coreosoidc.NewProvider(ctx, spec.Issuer)
	if err != nil {
		return nil, fmt.Errorf("could not initialize provider: %w", err)
	}
	providerJSON := &struct {
		JWKSURL string `json:"jwks_uri"`
	}{}
	if err := provider.Claims(providerJSON); err != nil {
		return nil, fmt.Errorf("could not get provider jwks_uri: %w", err) // should be impossible because coreosoidc.NewProvider validates this
	}
	if len(providerJSON.JWKSURL) == 0 {
		return nil, fmt.Errorf("issuer %q does not have jwks_uri set", spec.Issuer)
	}

	oidcAuthenticator, err := oidc.New(oidc.Options{
		IssuerURL:            spec.Issuer,
		KeySet:               coreosoidc.NewRemoteKeySet(ctx, providerJSON.JWKSURL),
		ClientID:             spec.Audience,
		UsernameClaim:        usernameClaim,
		GroupsClaim:          groupsClaim,
		SupportedSigningAlgs: defaultSupportedSigningAlgos(),
		// this is still needed for distributed claim resolution, meaning this uses a http client that does not honor our TLS config
		// TODO fix when we pick up https://github.com/kubernetes/kubernetes/pull/106141
		CAContentProvider: caContentProvider,
	})
	if err != nil {
		return nil, fmt.Errorf("could not initialize authenticator: %w", err)
	}

	return &jwtAuthenticator{
		tokenAuthenticatorCloser: oidcAuthenticator,
		spec:                     spec,
	}, nil
}
