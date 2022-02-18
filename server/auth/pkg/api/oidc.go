package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gigantum/hoss-auth/pkg/state"
	"github.com/gigantum/hoss-auth/pkg/userinfo"
	"github.com/gin-gonic/gin"
)

type OidcHosts struct {
	ExternalHost    string
	ConditionalHost string
}

// oidcError Sets a HTTP Bad Request status with the given error message
func oidcError(c *gin.Context, error_type, error_description string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":             error_type,
		"error_description": error_description,
	})
}

// isRedirectAllowed checks to see if the redirect URL is either the current server or an
// "additional allowed servers" set via the config file. If dev mode is 'true', a server
// that is not on the additional allowed servers list will still be allowed with a warning.
func isRedirectAllowed(redirectUriParsed *url.URL, additionalAllowedServers []string, devMode bool) bool {
	allowedServers := append([]string{os.Getenv("EXTERNAL_HOSTNAME")}, additionalAllowedServers...)

	for _, v := range allowedServers {
		if v == (redirectUriParsed.Scheme + "://" + redirectUriParsed.Host) {
			return true
		}
	}

	if devMode {
		log.Printf("WARNING: Redirect '%s' not in authorized server list, but allowed due to dev mode.\n",
			redirectUriParsed.Scheme+"://"+redirectUriParsed.Host)
		return true
	}

	return false
}

// oidcErrorRedirect Sets a HTTP redirect status with the given error message
func oidcErrorRedirect(c *gin.Context, redirectUri *url.URL, error_type, error_description string) {
	q := redirectUri.Query()
	q.Add("error", error_type)
	q.Add("error_description", error_description)
	redirectUri.RawQuery = q.Encode()

	t := redirectUri.String()

	c.Redirect(http.StatusFound, t)
}

// Authorize starts the oauth2 / oidc login workflow
// @Summary Start an OAuth2/OIDC login workflow
// @Schemes
// @Tags OpenID Connect
// @Description This open endpoint starts a login workflow. If successful it will redirect back to the redirectURI with
// @Description the session state variable.
// @Description If an error occurs after the user specified redirect is validated, the redirect will complete with
// @Description error information in the query params of the redirect ('error' and 'error_description')
// @Accept json
// @Produce json
// @Param	authorizeParams		query	state.AuthorizeArgs	true	"Authorize Query Params"
// @Success 302
// @Failure 400 {object} object{internal_error=string}
// @Failure 500 {object} string
// @Router /authorize [get]
func (a *Auth) Authorize(c *gin.Context) {
	log.Println("New auth " + c.Request.URL.String())

	s := &state.State{}

	if err := c.BindQuery(&s.AuthArgs); err != nil {
		oidcError(c, "invalid_request", err.Error())
		return
	}

	// Validate Arguments
	redirectUri, err := url.QueryUnescape(s.AuthArgs.RedirectURI)
	if err != nil {
		oidcError(c, "invalid_request", "Could not decode redirect_uri")
		return
	}

	redirectUriParsed, err := url.Parse(redirectUri)
	if err != nil {
		oidcError(c, "invalid_request", "Could not parse redirect_uri")
		return
	}

	// Since there is no client callback registration, ensure we are only redirecting allowed hosts
	// To be allowed, a host (scheme + hostname) must either be the EXTERNAL_HOSTNAME or in the
	// "additional_allowed_servers" list in the config file
	isAllowed := isRedirectAllowed(redirectUriParsed,
		a.Settings.AdditionalAllowedServers,
		a.Settings.DevServer)

	if !isAllowed {
		log.Println("Redirect URI not allowed: " + redirectUri)
		oidcError(c, "invalid_request", "redirect_uri host is invalid")
		return
	}

	s.AuthArgs.ParsedURL = redirectUriParsed

	if s.AuthArgs.ResponseType != "code" && s.AuthArgs.ResponseType != "id_token" {
		oidcErrorRedirect(c, s.AuthArgs.ParsedURL, "unsupported_response_type", "Only 'code' or 'id_token' response_type supported")
		return
	}

	if s.AuthArgs.ClientID != a.Settings.ClientID {
		oidcErrorRedirect(c, s.AuthArgs.ParsedURL, "invalid_request", "Unknown client_id")
		return
	}

	for _, scope := range strings.Fields(s.AuthArgs.Scope) {
		switch scope {
		case "openid":
			s.Scopes.ScopeOpenID = true
		case "profile":
			s.Scopes.ScopeProfile = true
		case "email":
			s.Scopes.ScopeEmail = true
		case "hoss":
			s.Scopes.ScopeHOSS = true
		default:
			oidcErrorRedirect(c, s.AuthArgs.ParsedURL, "invalid_scope", "Unsupported scope: "+scope)
			return
		}
	}

	if s.AuthArgs.ResponseType == "id_token" && !s.Scopes.ScopeOpenID {
		oidcErrorRedirect(c, s.AuthArgs.ParsedURL, "invalid_request", "Require openid scope for id_token")
		return
	}
	// End Validate Arguments

	sv := state.AddSession(s)

	redirectUri = a.GetAuthConfig(c).AuthCodeURL(sv) // Redirect to Dex for the authentication

	c.Redirect(http.StatusFound, redirectUri)
}

// Callback handles the response from the authentication server and retrieves the user's tokens
// @Summary Handle response from authentication service and creates JWTs
// @Schemes
// @Tags OpenID Connect
// @Description This endpoint will handle the response from the authentication service, generate Hoss JWTs, and then redirect
// @Description back to the user specified redirect location from the originating `/authorize` request with the JWTs.
// @Accept json
// @Produce json
// @Param	state		query	string	true	"Session state value"
// @Param	code		query	string	true	"Session code value"
// @Success 302
// @Failure 400 {object} object{internal_error=string}
// @Failure 500 {object} string
// @Router /callback [get]
func (a *Auth) Callback(c *gin.Context) {
	sv, ok := c.GetQuery("state")
	if !ok || sv == "" {
		oidcError(c, "invalid_request", "No 'state' value provided")
		return
	}

	log.Println("Auth callback: state = " + sv)

	sess, ok := state.GetSession(sv)
	if !ok {
		oidcError(c, "invalid_request", "Invalid 'state' value provided")
		return
	}

	code, ok := c.GetQuery("code")
	if !ok || code == "" {
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "invalid_request", "No 'code' value provided")
		return
	}

	_, claims, err := a.GetProviderTokens(c, code)
	if err != nil {
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "server_error", err.Error())
		return
	}

	userInfo, errClaims := userinfo.ParseUserClaims(claims,
		a.Settings.AdminGroup, a.Settings.PrivilegedGroup, a.Settings.UsernameClaim)

	if err := a.Database.CreateOrUpdateUser(userInfo); err != nil {
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "server_error", err.Error())
	}

	account := AccountInformation{
		Email: userInfo.Email,
	}
	errAllow := a.Allowlist.CheckAllowlist(account)

	if errClaims != nil || errAllow != nil {
		// Delete the session object
		_, _ = state.RemoveSession(sv)
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "access_denied", "Claims not valid or account is not on the allowlist")
		return
	}

	sess.Tokens, err = a.CreateGigantumTokens(userInfo, sess.Scopes, sess.AuthArgs.Nonce, false)
	if err != nil {
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "server_error", err.Error())
		return
	}

	// Add user to the automated "admin" and "public" groups
	err = a.addToAutoGroups(userInfo)
	if err != nil {
		oidcErrorRedirect(c, sess.AuthArgs.ParsedURL, "server_error", err.Error())
		return
	}

	if sess.AuthArgs.ResponseType == "code" {
		// Return a code that is used with the /token endpoint to get the tokens
		q := sess.AuthArgs.ParsedURL.Query()
		q.Add("code", sv)
		if sess.AuthArgs.State != "" {
			q.Add("state", sess.AuthArgs.State)
		}
		sess.AuthArgs.ParsedURL.RawQuery = q.Encode()

		u := sess.AuthArgs.ParsedURL.String()

		c.Redirect(http.StatusFound, u)
	} else if sess.AuthArgs.ResponseType == "id_token" {
		// Return the id_token in the URL fragment of the redirect response
		state.RemoveSession(sv)

		f := url.Values{}
		f.Add("id_token", sess.Tokens.IDToken)
		if sess.AuthArgs.State != "" {
			f.Add("state", sess.AuthArgs.State)
		}
		sess.AuthArgs.ParsedURL.Fragment = f.Encode()

		u := sess.AuthArgs.ParsedURL.String()
		c.Redirect(http.StatusFound, u)
	} else {
		panic("Not Implemented") // Should never happen
	}
}

// Token allows the client to exchange the code for the user's tokens, or to exchange a refresh
// token for a new set of tokens
// @Summary Exchange authorization code for tokens
// @Schemes
// @Tags OpenID Connect
// @Description This open endpoint allows a client to exchange an authorization code for JWTs if using the
// @Description authorization_code workflow.
// @Description Note, this endpoint also works with refresh tokens, but they are not fully implemented and supported
// @Description by the Hoss system at this time.
// @Accept json
// @Produce json
// @Param	tokenArgs		body	state.TokenArgs	true	"Token Exchange Args"
// @Success 200 {object} state.Tokens
// @Failure 400 {object} object{internal_error=string}
// @Failure 500 {object} string
// @Router /token [post]
func (a *Auth) Tokens(c *gin.Context) {
	var args state.TokenArgs
	if err := c.Bind(&args); err != nil {
		oidcError(c, "invalid_request", "Couldn't bind parameters")
		return
	}

	switch args.GrantType {
	case "authorization_code":

		// Validate arguments
		if args.ClientID != a.Settings.ClientID {
			oidcError(c, "invalid_request", "Unknown client_id")
			return
		}

		log.Println("Get tokens: state = " + args.Code)

		sess, ok := state.RemoveSession(args.Code)
		if !ok {
			oidcError(c, "invalid_grant", "invalid code")
			return
		}

		if args.RedirectURI != "" && args.RedirectURI != sess.AuthArgs.RedirectURI {
			oidcError(c, "invalid_request", "invalid redirect_uri")
			return
		}
		// End Validate Arguments

		c.JSON(http.StatusOK, sess.Tokens)

	case "refresh_token":

		// validate refresh token
		userClaims, err := a.ValidateToken(args.RefreshToken)
		if err != nil {
			oidcError(c, "server_error", err.Error())
			return
		}

		var scopes state.Scopes
		originalScopes := userClaims["scopes"].(string)
		for _, scope := range strings.Fields(args.Scope) {
			switch scope {
			case "openid":
				if !strings.Contains(originalScopes, "openid") {
					oidcError(c, "invalid_scope", "Scope `openid` was not one of the original authorized scopes")
				}
				scopes.ScopeOpenID = true
			case "profile":
				if !strings.Contains(originalScopes, "profile") {
					oidcError(c, "invalid_scope", "Scope `profile` was not one of the original authorized scopes")
				}
				scopes.ScopeProfile = true
			case "email":
				if !strings.Contains(originalScopes, "email") {
					oidcError(c, "invalid_scope", "Scope `email` was not one of the original authorized scopes")
				}
				scopes.ScopeEmail = true
			case "hoss":
				if !strings.Contains(originalScopes, "hoss") {
					oidcError(c, "invalid_scope", "Scope `hoss` was not one of the original authorized scopes")
				}
				scopes.ScopeHOSS = true
			default:
				oidcError(c, "invalid_scope", "Unsupported scope: "+scope)
				return
			}
		}

		// get user info from database based on username
		userRecord, err := a.Database.GetUser(userClaims["nickname"].(string))
		if err != nil {
			oidcError(c, "server_error", err.Error())
			return
		}

		userInfo := userinfo.UserInfo{
			Subject:       userRecord.Subject,
			FullName:      userRecord.FullName,
			GivenName:     userRecord.GivenName,
			FamilyName:    userRecord.FamilyName,
			Username:      userRecord.Username,
			Email:         userRecord.Email,
			EmailVerified: &userRecord.EmailVerified,
			Role:          userRecord.Role,
		}

		tokens, err := a.CreateGigantumTokens(userInfo, scopes, args.Nonce, false)
		if err != nil {
			oidcError(c, "server_error", err.Error())
			return
		}

		c.JSON(http.StatusOK, tokens)

	default:
		oidcError(c, "invalid_request", "grant_type must be one of `authorization_grant` or `refresh_token`")
	}
}

// UserInfo allows the client to get the decoded user's profile
// @Summary Get user profile information
// @Schemes
// @Tags OpenID Connect
// @Description This open endpoint allows a client fetch the user's profile information
// @Accept json
// @Produce json
// @Success 200 {object} userinfo.UserInfo
// @Failure 400 {object} object{internal_error=string}
// @Failure 500 {object} string
// @Router /userinfo [get]
func (a *Auth) UserInfo(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "invalid_token",
			"error_description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// JWKS returns the JWKS JSON object containing the JWT signing key information
// @Summary Fetch JWKS
// @Schemes
// @Tags OpenID Connect
// @Description This open endpoint returns the JWKS for this OIDC provider
// @Accept json
// @Produce json
// @Success 200
// @Router /keys [get]
func (a *Auth) JWKS(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, a.JwksJSON)
}

// OpenID provides the openid configuration json
// @Summary OpenID Connect Configuration
// @Schemes
// @Tags OpenID Connect
// @Description This open endpoint returns the OIDC configuration information for clients
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} object{internal_error=string}
// @Failure 500 {object} string
// @Router /.well-known/openid-configuration [get]
func (a *Auth) OpenID(c *gin.Context) {
	oh := OidcHosts{}
	if c.Request.Host == "auth:8080" {
		// internal route
		oh = OidcHosts{
			ExternalHost:    os.Getenv("EXTERNAL_HOSTNAME"),
			ConditionalHost: "http://auth:8080",
		}
	} else {
		// external route
		oh = OidcHosts{
			ExternalHost:    os.Getenv("EXTERNAL_HOSTNAME"),
			ConditionalHost: fmt.Sprintf("%s/auth", os.Getenv("EXTERNAL_HOSTNAME")),
		}
	}
	var w http.ResponseWriter = c.Writer

	t, err := template.New("oidchosts").Parse(a.OpenIDConfig)
	if err != nil {
		oidcError(c, "internal_error", "Failed to parse oidc .well-known data: "+err.Error())
	}
	err = t.Execute(w, oh)
	if err != nil {
		oidcError(c, "internal_error", "Failed to parse oidc .well-known data: "+err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
