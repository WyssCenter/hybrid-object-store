package api

import (
	"net/url"
	"testing"

	"github.com/gigantum/hoss-auth/pkg/test"
)

func TestAllowedHostExternal(t *testing.T) {
	settings, err := SetupLoadConfig(t)
	if err != nil {
		t.Fatal("Failed to setup test: ", err.Error())
	}

	redirectUriParsed, err := url.Parse("http://localhost")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed := isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if !isAllowed {
		t.Fatal("Redirect NOT allowed. Expected to be allowed")
	}

	redirectUriParsed, err = url.Parse("https://localhost")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if isAllowed {
		t.Fatal("Redirect allowed. Expected to be NOT allowed")
	}
}

func TestAllowedHostAdditional(t *testing.T) {
	settings, err := SetupLoadConfig(t)
	if err != nil {
		t.Fatal("Failed to setup test: ", err.Error())
	}

	test.AssertEqual(t, len(settings.AdditionalAllowedServers), 0)
	settings.AdditionalAllowedServers = []string{"https://myserver.com", "http://myserver2.com"}
	test.AssertEqual(t, len(settings.AdditionalAllowedServers), 2)

	redirectUriParsed, err := url.Parse("http://localhost/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed := isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if !isAllowed {
		t.Fatal("Redirect NOT allowed. Expected to be allowed")
	}

	redirectUriParsed, err = url.Parse("http://myserver.com/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if isAllowed {
		t.Fatal("Redirect allowed. Expected to be NOT allowed")
	}

	redirectUriParsed, err = url.Parse("https://myserver.com/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if !isAllowed {
		t.Fatal("Redirect NOT allowed. Expected to be allowed")
	}

	redirectUriParsed, err = url.Parse("https://myserver2.com/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if isAllowed {
		t.Fatal("Redirect allowed. Expected to be NOT allowed")
	}

	redirectUriParsed, err = url.Parse("http://myserver2.com/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if !isAllowed {
		t.Fatal("Redirect NOT allowed. Expected to be allowed")
	}

	redirectUriParsed, err = url.Parse("https://notallowed.com/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if isAllowed {
		t.Fatal("Redirect allowed. Expected to be NOT allowed")
	}
}

func TestAllowedHostDevMode(t *testing.T) {
	settings, err := SetupLoadConfig(t)
	if err != nil {
		t.Fatal("Failed to setup test: ", err.Error())
	}

	redirectUriParsed, err := url.Parse("http://localhost:3000/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed := isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, false)
	if isAllowed {
		t.Fatal("Redirect allowed. Expected to be NOT allowed")
	}

	redirectUriParsed, err = url.Parse("http://localhost:3000/auth/v1/callback")
	if err != nil {
		t.Fatal("Could not parse redirect_uri: ", err.Error())
	}
	isAllowed = isRedirectAllowed(redirectUriParsed, settings.AdditionalAllowedServers, true)
	if !isAllowed {
		t.Fatal("Redirect NOT allowed. Expected to be allowed")
	}
}
