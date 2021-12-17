// Copyright 2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package upstreamwatchers

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"

	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/constable"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/oidc/provider"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/plog"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/upstreamldap"
	"go.pinniped.dev/generated/latest/apis/supervisor/idp/v1alpha1"
)

const (
	ReasonNotFound         = "SecretNotFound"
	ReasonWrongType        = "SecretWrongType"
	ReasonMissingKeys      = "SecretMissingKeys"
	ReasonSuccess          = "Success"
	ReasonInvalidTLSConfig = "InvalidTLSConfig"

	ErrNoCertificates = constable.Error("no certificates found")

	LDAPBindAccountSecretType = corev1.SecretTypeBasicAuth
	probeLDAPTimeout          = 90 * time.Second

	// Constants related to conditions.
	typeBindSecretValid              = "BindSecretValid"
	typeTLSConfigurationValid        = "TLSConfigurationValid"
	typeLDAPConnectionValid          = "LDAPConnectionValid"
	TypeSearchBaseFound              = "SearchBaseFound"
	reasonLDAPConnectionError        = "LDAPConnectionError"
	noTLSConfigurationMessage        = "no TLS configuration provided"
	loadedTLSConfigurationMessage    = "loaded TLS configuration"
	ReasonUsingConfigurationFromSpec = "UsingConfigurationFromSpec"
	ReasonErrorFetchingSearchBase    = "ErrorFetchingSearchBase"
)

// An in-memory cache with an entry for each ActiveDirectoryIdentityProvider, to keep track of which ResourceVersion
// of the bind Secret, which TLS/StartTLS setting was used and which search base was found during the most recent successful validation.
type SecretVersionCacheI interface {
	Get(upstreamName, resourceVersion string, generation int64) (ValidatedSettings, bool)
	Set(upstreamName, resourceVersion string, generation int64, settings ValidatedSettings)
}

type SecretVersionCache struct {
	ValidatedSettingsByName map[string]ValidatedSettings
}

func (s *SecretVersionCache) Get(upstreamName, resourceVersion string, generation int64) (ValidatedSettings, bool) {
	validatedSettings := s.ValidatedSettingsByName[upstreamName]
	if validatedSettings.BindSecretResourceVersion == resourceVersion &&
		validatedSettings.Generation == generation && validatedSettings.UserSearchBase != "" &&
		validatedSettings.GroupSearchBase != "" && validatedSettings.LDAPConnectionProtocol != "" {
		return validatedSettings, true
	}
	return ValidatedSettings{}, false
}

func (s *SecretVersionCache) Set(upstreamName, resourceVersion string, generation int64, settings ValidatedSettings) {
	s.ValidatedSettingsByName[upstreamName] = settings
}

type ValidatedSettings struct {
	Generation                int64
	BindSecretResourceVersion string
	LDAPConnectionProtocol    upstreamldap.LDAPConnectionProtocol
	UserSearchBase            string
	GroupSearchBase           string
}

func NewSecretVersionCache() SecretVersionCacheI {
	cache := SecretVersionCache{ValidatedSettingsByName: map[string]ValidatedSettings{}}
	return &cache
}

// read only interface for sharing between ldap and active directory.
type UpstreamGenericLDAPIDP interface {
	Spec() UpstreamGenericLDAPSpec
	Name() string
	Namespace() string
	Generation() int64
	Status() UpstreamGenericLDAPStatus
}

type UpstreamGenericLDAPSpec interface {
	Host() string
	TLSSpec() *v1alpha1.TLSSpec
	BindSecretName() string
	UserSearch() UpstreamGenericLDAPUserSearch
	GroupSearch() UpstreamGenericLDAPGroupSearch
	DetectAndSetSearchBase(ctx context.Context, config *upstreamldap.ProviderConfig) *v1alpha1.Condition
}

type UpstreamGenericLDAPUserSearch interface {
	Base() string
	Filter() string
	UsernameAttribute() string
	UIDAttribute() string
}

type UpstreamGenericLDAPGroupSearch interface {
	Base() string
	Filter() string
	GroupNameAttribute() string
}

type UpstreamGenericLDAPStatus interface {
	Conditions() []v1alpha1.Condition
}

func ValidateTLSConfig(tlsSpec *v1alpha1.TLSSpec, config *upstreamldap.ProviderConfig) *v1alpha1.Condition {
	if tlsSpec == nil {
		return validTLSCondition(noTLSConfigurationMessage)
	}
	if len(tlsSpec.CertificateAuthorityData) == 0 {
		return validTLSCondition(loadedTLSConfigurationMessage)
	}

	bundle, err := base64.StdEncoding.DecodeString(tlsSpec.CertificateAuthorityData)
	if err != nil {
		return invalidTLSCondition(fmt.Sprintf("certificateAuthorityData is invalid: %s", err.Error()))
	}

	ca := x509.NewCertPool()
	ok := ca.AppendCertsFromPEM(bundle)
	if !ok {
		return invalidTLSCondition(fmt.Sprintf("certificateAuthorityData is invalid: %s", ErrNoCertificates))
	}

	config.CABundle = bundle
	return validTLSCondition(loadedTLSConfigurationMessage)
}

func TestConnection(
	ctx context.Context,
	bindSecretName string,
	config *upstreamldap.ProviderConfig,
	currentSecretVersion string,
) *v1alpha1.Condition {
	// First try using TLS.
	config.ConnectionProtocol = upstreamldap.TLS
	tlsLDAPProvider := upstreamldap.New(*config)
	err := tlsLDAPProvider.TestConnection(ctx)
	if err != nil {
		plog.InfoErr("testing LDAP connection using TLS failed, so trying again with StartTLS", err, "host", config.Host)
		// If there was any error, try again with StartTLS instead.
		config.ConnectionProtocol = upstreamldap.StartTLS
		startTLSLDAPProvider := upstreamldap.New(*config)
		startTLSErr := startTLSLDAPProvider.TestConnection(ctx)
		if startTLSErr == nil {
			plog.Info("testing LDAP connection using StartTLS succeeded", "host", config.Host)
			// Successfully able to fall back to using StartTLS, so clear the original
			// error and consider the connection test to be successful.
			err = nil
		} else {
			plog.InfoErr("testing LDAP connection using StartTLS also failed", err, "host", config.Host)
			// Falling back to StartTLS also failed, so put TLS back into the config
			// and consider the connection test to be failed.
			config.ConnectionProtocol = upstreamldap.TLS
		}
	}

	if err != nil {
		return &v1alpha1.Condition{
			Type:   typeLDAPConnectionValid,
			Status: v1alpha1.ConditionFalse,
			Reason: reasonLDAPConnectionError,
			Message: fmt.Sprintf(`could not successfully connect to "%s" and bind as user "%s": %s`,
				config.Host, config.BindUsername, err.Error()),
		}
	}

	return &v1alpha1.Condition{
		Type:   typeLDAPConnectionValid,
		Status: v1alpha1.ConditionTrue,
		Reason: ReasonSuccess,
		Message: fmt.Sprintf(`successfully able to connect to "%s" and bind as user "%s" [validated with Secret "%s" at version "%s"]`,
			config.Host, config.BindUsername, bindSecretName, currentSecretVersion),
	}
}

func validTLSCondition(message string) *v1alpha1.Condition {
	return &v1alpha1.Condition{
		Type:    typeTLSConfigurationValid,
		Status:  v1alpha1.ConditionTrue,
		Reason:  ReasonSuccess,
		Message: message,
	}
}

func invalidTLSCondition(message string) *v1alpha1.Condition {
	return &v1alpha1.Condition{
		Type:    typeTLSConfigurationValid,
		Status:  v1alpha1.ConditionFalse,
		Reason:  ReasonInvalidTLSConfig,
		Message: message,
	}
}

func ValidateSecret(secretInformer corev1informers.SecretInformer, secretName string, secretNamespace string, config *upstreamldap.ProviderConfig) (*v1alpha1.Condition, string) {
	secret, err := secretInformer.Lister().Secrets(secretNamespace).Get(secretName)
	if err != nil {
		return &v1alpha1.Condition{
			Type:    typeBindSecretValid,
			Status:  v1alpha1.ConditionFalse,
			Reason:  ReasonNotFound,
			Message: err.Error(),
		}, ""
	}

	if secret.Type != corev1.SecretTypeBasicAuth {
		return &v1alpha1.Condition{
			Type:   typeBindSecretValid,
			Status: v1alpha1.ConditionFalse,
			Reason: ReasonWrongType,
			Message: fmt.Sprintf("referenced Secret %q has wrong type %q (should be %q)",
				secretName, secret.Type, corev1.SecretTypeBasicAuth),
		}, secret.ResourceVersion
	}

	config.BindUsername = string(secret.Data[corev1.BasicAuthUsernameKey])
	config.BindPassword = string(secret.Data[corev1.BasicAuthPasswordKey])
	if len(config.BindUsername) == 0 || len(config.BindPassword) == 0 {
		return &v1alpha1.Condition{
			Type:   typeBindSecretValid,
			Status: v1alpha1.ConditionFalse,
			Reason: ReasonMissingKeys,
			Message: fmt.Sprintf("referenced Secret %q is missing required keys %q",
				secretName, []string{corev1.BasicAuthUsernameKey, corev1.BasicAuthPasswordKey}),
		}, secret.ResourceVersion
	}

	return &v1alpha1.Condition{
		Type:    typeBindSecretValid,
		Status:  v1alpha1.ConditionTrue,
		Reason:  ReasonSuccess,
		Message: "loaded bind secret",
	}, secret.ResourceVersion
}

type GradatedConditions struct {
	gradatedConditions []GradatedCondition
}

func (g *GradatedConditions) Conditions() []*v1alpha1.Condition {
	conditions := []*v1alpha1.Condition{}
	for _, gc := range g.gradatedConditions {
		conditions = append(conditions, gc.condition)
	}
	return conditions
}

func (g *GradatedConditions) Append(condition *v1alpha1.Condition, isFatal bool) {
	g.gradatedConditions = append(g.gradatedConditions, GradatedCondition{condition: condition, isFatal: isFatal})
}

// A condition and a boolean that tells you whether it's fatal or just a warning.
type GradatedCondition struct {
	condition *v1alpha1.Condition
	isFatal   bool
}

func ValidateGenericLDAP(ctx context.Context, upstream UpstreamGenericLDAPIDP, secretInformer corev1informers.SecretInformer, validatedSecretVersionsCache SecretVersionCacheI, config *upstreamldap.ProviderConfig) GradatedConditions {
	conditions := GradatedConditions{}
	secretValidCondition, currentSecretVersion := ValidateSecret(secretInformer, upstream.Spec().BindSecretName(), upstream.Namespace(), config)
	conditions.Append(secretValidCondition, true)
	tlsValidCondition := ValidateTLSConfig(upstream.Spec().TLSSpec(), config)
	conditions.Append(tlsValidCondition, true)

	// No point in trying to connect to the server if the config was already determined to be invalid.
	var ldapConnectionValidCondition *v1alpha1.Condition
	var searchBaseFoundCondition *v1alpha1.Condition
	if secretValidCondition.Status == v1alpha1.ConditionTrue && tlsValidCondition.Status == v1alpha1.ConditionTrue {
		ldapConnectionValidCondition, searchBaseFoundCondition = validateAndSetLDAPServerConnectivityAndSearchBase(ctx, validatedSecretVersionsCache, upstream, config, currentSecretVersion)
		if ldapConnectionValidCondition != nil {
			conditions.Append(ldapConnectionValidCondition, false)
		}
		if searchBaseFoundCondition != nil {
			conditions.Append(searchBaseFoundCondition, true)
		}
	}
	return conditions
}

func validateAndSetLDAPServerConnectivityAndSearchBase(ctx context.Context, validatedSecretVersionsCache SecretVersionCacheI, upstream UpstreamGenericLDAPIDP, config *upstreamldap.ProviderConfig, currentSecretVersion string) (*v1alpha1.Condition, *v1alpha1.Condition) {
	// previouslyValidatedSecretVersion := validatedSecretVersionsCache.ValidatedSettingsByName[upstream.Name()].BindSecretResourceVersion
	// doesn't have an existing entry for ValidatedSettingsByName with this secret version ->
	// lets double check tls connection
	// if we can connect, put it in the secret cache
	// also we KNOW we need to recheck the search base stuff too... so they should all be one function?
	// but if tls validation fails no need to also try to get search base stuff?

	validatedSettings, hasPreviousValidatedSettings := validatedSecretVersionsCache.Get(upstream.Name(), currentSecretVersion, upstream.Generation())
	var ldapConnectionValidCondition, searchBaseFoundCondition *v1alpha1.Condition
	if !hasPreviousValidatedSettings {
		testConnectionTimeout, cancelFunc := context.WithTimeout(ctx, probeLDAPTimeout)
		defer cancelFunc()

		ldapConnectionValidCondition = TestConnection(testConnectionTimeout, upstream.Spec().BindSecretName(), config, currentSecretVersion)

		searchBaseTimeout, cancelFunc := context.WithTimeout(ctx, probeLDAPTimeout)
		defer cancelFunc()
		searchBaseFoundCondition = upstream.Spec().DetectAndSetSearchBase(searchBaseTimeout, config)

		if ldapConnectionValidCondition.Status == v1alpha1.ConditionTrue {
			// if it's nil, don't worry about the search base condition. But if it exists make sure the status is true.
			if searchBaseFoundCondition == nil || (searchBaseFoundCondition.Status == v1alpha1.ConditionTrue) {
				// Remember (in-memory for this pod) that the controller has successfully validated the LDAP provider
				// using this version of the Secret. This is for performance reasons, to avoid attempting to connect to
				// the LDAP server more than is needed. If the pod restarts, it will attempt this validation again.
				validatedSettings.LDAPConnectionProtocol = config.ConnectionProtocol
				validatedSettings.BindSecretResourceVersion = currentSecretVersion
				validatedSettings.Generation = upstream.Generation()
				validatedSettings.UserSearchBase = config.UserSearch.Base
				validatedSettings.GroupSearchBase = config.GroupSearch.Base
				validatedSecretVersionsCache.Set(upstream.Name(), currentSecretVersion, upstream.Generation(), validatedSettings)
			}
		}
	} else {
		config.ConnectionProtocol = validatedSettings.LDAPConnectionProtocol
		config.UserSearch.Base = validatedSettings.UserSearchBase
		config.GroupSearch.Base = validatedSettings.GroupSearchBase
	}

	return ldapConnectionValidCondition, searchBaseFoundCondition
}

func EvaluateConditions(conditions GradatedConditions, config *upstreamldap.ProviderConfig) (provider.UpstreamLDAPIdentityProviderI, bool) {
	for _, gradatedCondition := range conditions.gradatedConditions {
		if gradatedCondition.condition.Status != v1alpha1.ConditionTrue && gradatedCondition.isFatal {
			// Invalid provider, so do not load it into the cache.
			return nil, true
		}
	}

	for _, gradatedCondition := range conditions.gradatedConditions {
		if gradatedCondition.condition.Status != v1alpha1.ConditionTrue && !gradatedCondition.isFatal {
			// Error but load it into the cache anyway, treating this condition failure more like a warning.
			// Try again hoping that the condition will improve.
			return upstreamldap.New(*config), true
		}
	}
	// Fully validated provider, so load it into the cache.
	return upstreamldap.New(*config), false
}
