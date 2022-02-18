package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/gigantum/hoss-auth/pkg/database"
	service "github.com/gigantum/hoss-service"
)

// Settings contains user-defined settings read from config yaml
type Settings struct {
	DevServer                bool                 `yaml:"dev_server"`
	ClientID                 string               `yaml:"client_id"`
	ClientSecret             string               `yaml:"client_secret"`
	Issuer                   string               `yaml:"issuer"`
	Core                     string               `yaml:"core"`
	OpenIDConfigFile         string               `yaml:"open_id_config_file"`
	AdminGroup               string               `yaml:"admin_group"`
	PrivilegedGroup          string               `yaml:"privileged_group"`
	TokenExpirationHours     TokenExpirationHours `yaml:"token_expiration_hours"`
	AdditionalAllowedServers []string             `yaml:"additional_allowed_servers"`
	PasswordPolicy           PasswordPolicy       `yaml:"password_policy"`
	UsernameClaim            string               `yaml:"username_claim"`
}

type TokenExpirationHours struct {
	Access  int `yaml:"access"`
	Id      int `yaml:"id"`
	Refresh int `yaml:"refresh"`
}

type PasswordPolicy struct {
	MinLength        int  `yaml:"min_length"`
	RequireUppercase bool `yaml:"require_uppercase"`
	RequireSpecial   bool `yaml:"require_special"`
}

// Auth contains auth info
type Auth struct {
	Settings *Settings

	Hostname  string
	JWTIssuer string

	Ctx      context.Context
	Config   oauth2.Config
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier

	KeyID        string
	JwtKey       *rsa.PrivateKey
	JwksJSON     string
	OpenIDConfig string

	Allowlist *AccountAllowlist
	Database  *database.Database
}

// LoadConfig creates a default config and then initializes it with values from
// the default config file location.
func LoadConfig() *Auth {
	settings := &Settings{}
	settingsBytes, err := ioutil.ReadFile(filepath.Join("/opt", "config.yaml"))
	if err == nil {
		if err = yaml.Unmarshal(settingsBytes, &settings); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	auth := &Auth{
		Settings: settings,
		Hostname: os.Getenv("EXTERNAL_HOSTNAME"),
		// Context for exchanging the OAuth2 code for the user's token. Need to explicitly create it so that
		// SSL verification can be disabled if running under dev mode
		Ctx: context.WithValue(context.TODO(), oauth2.HTTPClient, &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: settings.DevServer,
				},
			},
		}),

		Provider: &oidc.Provider{},
		Verifier: (&oidc.Provider{}).Verifier(&oidc.Config{ClientID: settings.ClientID}),

		// keyID is the kid value for the current RSA key in use by the service
		// future work may support multiple keys/key rolling.
		KeyID:        "",
		JwtKey:       &rsa.PrivateKey{},
		JwksJSON:     "{}",
		OpenIDConfig: "{}",

		Allowlist: &AccountAllowlist{},
		Database:  database.Load(),
	}

	// Initialize the OIDC provider and verifier
	auth.Provider, err = oidc.NewProvider(auth.Ctx, auth.Settings.Issuer)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	auth.Verifier = auth.Provider.Verifier(&oidc.Config{ClientID: settings.ClientID})

	pemEncoded, err := ioutil.ReadFile("/secrets/private.pem")
	if err == nil {
		log.Println("Loading JWT key from disk")
		block, _ := pem.Decode([]byte(pemEncoded))
		if block == nil {
			log.Println("Could not decode key PEM")
			os.Exit(2)
		}
		if block.Type != "RSA PRIVATE KEY" {
			log.Println("Unexpected key type")
			os.Exit(2)
		}

		auth.JwtKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Println(err.Error())
			os.Exit(2)
		}
	} else {
		log.Println("Generating new JWT key")
		// Initialize the jwtKey and jwksJSON variables
		auth.JwtKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatal(err)
		}

		pemEncoded := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(auth.JwtKey),
			},
		)

		err = ioutil.WriteFile("/secrets/private.pem", pemEncoded, 0600)
		if err != nil {
			log.Println("Couldn't save JWT private key: " + err.Error())
		}
	}

	jwk := service.JSONWebKeyFromRSA(auth.JwtKey)
	auth.KeyID = jwk.Kid

	jwks := service.Jwks{
		Keys: []service.JSONWebKey{jwk},
	}

	// Since the JWKS information is static, render it once instead of with each request
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	e.Encode(jwks)
	auth.JwksJSON = buf.String()

	// Load the openid config
	log.Println(settings.OpenIDConfigFile)
	openIDConfigFile, err := os.Open(settings.OpenIDConfigFile)
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}
	defer openIDConfigFile.Close()
	openIDConfigBytes, _ := ioutil.ReadAll(openIDConfigFile)
	auth.OpenIDConfig = string(openIDConfigBytes)

	// Parse out and use the issuer set in the oidc config, so that both
	// the config and the tokens match.
	// This allows an admin to configure multiple auth services to work together
	var oidcConfig map[string]interface{}
	if err := json.Unmarshal(openIDConfigBytes, &oidcConfig); err != nil {
		log.Fatal("Could not unmarshal oidc config: " + err.Error())
	}

	issuer, ok := oidcConfig["issuer"].(string)
	if !ok {
		log.Fatal("Could not locate oidc issuer")
	}

	auth.JWTIssuer = strings.ReplaceAll(issuer, "{{.ExternalHost}}", os.Getenv("EXTERNAL_HOSTNAME"))

	// Initialize the allowlist
	auth.Allowlist, err = LoadAllowlist()
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}

	// Initialize admin and public groups
	_, err = auth.Database.GetGroup("admin")
	if err != nil {
		// If error group not created yet.
		_, err = auth.Database.CreateGroup("admin",
			"Group for all admin users. This group is automatically added to all datasets and can only be removed by an administrator.")
		if err != nil {
			log.Println("Warning: Failed to create admin group")
		}
	}

	_, err = auth.Database.GetGroup("public")
	if err != nil {
		// If error group not created yet.
		_, err = auth.Database.CreateGroup("public",
			"Group containing *all* users. You can grant this group read-only permissions to a dataset.")
		if err != nil {
			log.Println("Warning: Failed to create public group")
		}
	}

	return auth
}

// GetAuthConfig gets the OAuth2 configuration based on the request host header,
// allowing the service to work without specifying the hostname at start time
func (a *Auth) GetAuthConfig(c *gin.Context) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.Settings.ClientID,
		ClientSecret: a.Settings.ClientSecret,
		Scopes:       []string{"openid", "profile", "email", "groups"},
		RedirectURL:  a.Hostname + "/auth/v1/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  a.Hostname + "/dex/auth",
			TokenURL: a.Settings.Issuer + "/token",
		},
	}
}

// GetLocalAuthConfig gets the OAuth2 configuration, for the Auth Service,
// based on the request host header, allowing the service to work without
// specifying the hostname at start time
func (a *Auth) GetLocalAuthConfig(c *gin.Context) *oauth2.Config {
	return &oauth2.Config{
		ClientID:    a.Settings.ClientID,
		Scopes:      []string{"openid", "profile", "email", "hoss"},
		RedirectURL: a.Hostname + "/auth/v1/login",
		Endpoint: oauth2.Endpoint{
			AuthURL:  a.Hostname + "/auth/v1/authorize",
			TokenURL: "http://auth:8080/v1/token", // use the internal hostname, as localhost will fail
		},
	}
}
