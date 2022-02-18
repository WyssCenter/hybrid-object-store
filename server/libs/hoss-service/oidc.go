package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	errors "github.com/gigantum/hoss-error"
)

type OpenIDConfiguration struct {
	Issuer string `json:"issuer"`

	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	JwksUri               string `json:"jwks_uri"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
}

func GetOpenIDConfiguration(provider string) (OpenIDConfiguration, error) {
	config := OpenIDConfiguration{}

	wellknown := provider
	if !strings.HasSuffix(wellknown, "/") {
		wellknown = wellknown + "/"
	}
	wellknown = wellknown + ".well-known/openid-configuration"
	resp, err := http.Get(wellknown)

	if err != nil {
		return config, err
	}

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return config, err
	}

	if resp.StatusCode != 200 {
		return config, errors.New("Could not retrieve OpenID configuration: " + resp.Status)
	}

	err = json.Unmarshal(responseBytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
