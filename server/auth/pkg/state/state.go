package state

import (
	"math/rand"
	"net/url"
)

// Character set used to generate a state value
const (
	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// Tokens is the JSON response containing both the access and id tokens for the logged in user
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token,omitempty"`
}

// AuthorizeArgs is the arguments to the /authorize endpoint and contains
// fields for the results of parsing the arguments
type AuthorizeArgs struct {
	ResponseType string `form:"response_type"`
	ClientID     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
	Nonce        string `form:"nonce"`

	// set during validation
	ParsedURL *url.URL `form:"-" swaggerignore:"true"`
}

type Scopes struct {
	ScopeOpenID  bool `form:"-"`
	ScopeProfile bool `form:"-"`
	ScopeEmail   bool `form:"-"`
	ScopeHOSS    bool `form:"-"`
}

// TokenArgs is the arguments to the /token endpoint
type TokenArgs struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	RefreshToken string `form:"refresh_token"`
	Scope        string `form:"scope"`
	Nonce        string `form:"nonce"`
}

// State is the internal data containing information about a given login attempt
type State struct {
	AuthArgs AuthorizeArgs
	Scopes   Scopes
	Tokens   Tokens // The user's tokens to return to the caller when requested
}

var (
	state = map[string]*State{}
)

// GenerateStateValue creates a psudo-random string for using as a unique identifier for each loging
func GenerateStateValue() string {
	buf := make([]byte, 12)
	for i := range buf {
		buf[i] = charset[rand.Intn(len(charset))]
	}
	return string(buf)
}

// AddSession adds an object to the internal storage and return the key to retrieve it
func AddSession(session *State) string {
	sv := GenerateStateValue()
	state[sv] = session
	return sv
}

// GetSession gets an object from the internal store
func GetSession(sv string) (*State, bool) {
	session, ok := state[sv]
	return session, ok
}

// RemoveSession removes an object from the internal store and returns the object
func RemoveSession(sv string) (*State, bool) {
	session, ok := state[sv]
	if ok {
		delete(state, sv)
	}

	return session, ok
}
