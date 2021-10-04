package authnz

import (
	"fmt"
	"net/http"

	appv1alpha1 "github.com/getupio-undistro/undistro/apis/app/v1alpha1"
	"github.com/getupio-undistro/undistro/pkg/scheme"
	"github.com/getupio-undistro/undistro/pkg/undistro"
	"github.com/getupio-undistro/undistro/third_party/pinniped/internal/httputil/httperr"
	"go.pinniped.dev/pkg/oidcclient/nonce"
	"go.pinniped.dev/pkg/oidcclient/pkce"
	"go.pinniped.dev/pkg/oidcclient/state"
	"golang.org/x/oauth2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	supervisorAuthorizeUpstreamNameParam = "pinniped_idp_name"
	supervisorAuthorizeUpstreamTypeParam = "pinniped_idp_type"
)

func (h HandlerState) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	idpName := params.Get("idp")
	if idpName == "" {
		return httperr.Newf(http.StatusBadRequest, "missing 'idp' parameter or query string")
	}
	idpFmtName := fmt.Sprintf("undistro-%s-idp", idpName)
	h.upstreamIdentityProviderName = idpFmtName
	h.Ctx = r.Context()
	// Initialize login parameters.
	var err error
	h.State, err = state.Generate()
	if err != nil {
		return err
	}
	h.Nonce, err = nonce.Generate()
	if err != nil {
		return err
	}
	h.PKCE, err = pkce.Generate()
	if err != nil {
		return err
	}
	// Prepare the common options for the authorization URL. We don't have the redirect URL yet though.
	authorizeOptions := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		h.Nonce.Param(),
		h.PKCE.Challenge(),
		h.PKCE.Method(),
	}
	// Get the callback url from helm release
	c, err := client.New(h.RestConf, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return err
	}
	hr := appv1alpha1.HelmRelease{}
	hrKey := client.ObjectKey{
		Namespace: undistro.Namespace,
		Name:      "pinniped-supervisor-management",
	}
	// get supervisor helm release
	err = c.Get(h.Ctx, hrKey, &hr)
	if err != nil {
		return err
	}
	if h.upstreamIdentityProviderName != "" {
		authorizeOptions = append(authorizeOptions, oauth2.SetAuthURLParam(supervisorAuthorizeUpstreamNameParam, h.upstreamIdentityProviderName))
		authorizeOptions = append(authorizeOptions, oauth2.SetAuthURLParam(supervisorAuthorizeUpstreamTypeParam, h.upstreamIdentityProviderType))
	}
	// Perform OIDC discovery.
	if h, err = h.initOIDCDiscovery(); err != nil {
		return err
	}
	config := hr.ValuesAsMap()["config"].(map[string]interface{})
	h.OAuth2Config.RedirectURL = config["callbackURL"].(string)
	session, err := h.SessionStore.Get(r, "undistro-login")
	if err != nil {
		return err
	}
	session.Values["state"] = h.State.String()
	session.Values["nonce"] = h.Nonce.String()
	session.Values["pkce"] = string(h.PKCE)
	session.Values["redirectURL"] = h.OAuth2Config.RedirectURL
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	authorizeURL := h.OAuth2Config.AuthCodeURL(h.State.String(), authorizeOptions...)
	http.Redirect(w, r, authorizeURL, http.StatusTemporaryRedirect)
	return nil
}
