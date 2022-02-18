package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	errors "github.com/gigantum/hoss-error"
)

// Tokens contains the OIDC tokens returned by the Auth Service
type Tokens struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	IDToken      string `json:"id_token" binding:"required"`
}

// GetServiceJWT exchanges the SERVICE_AUTH_SECRET env var token for a service JWT
func GetServiceJWT(authService string) (*Tokens, error) {
	bearer_token := os.Getenv("SERVICE_AUTH_SECRET")
	if bearer_token == "" {
		return nil, errors.New("SERVICE_AUTH_SECRET not defined")
	}

	req, err := http.NewRequest("POST", authService+"/sat/exchange/jwt", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearer_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with SAT exchange response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with SAT exchange response: StatusCode != 201")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokens Tokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

// RenewingTokens defines an object that has a background refresh routine that can be started
// and a way to get the current tokens. This allows a service to keep the tokens fresh while
// also sharing the same token information and not having to query the tokens multiple times.
type RenewingTokens interface {
	GetAccessToken() (string, error)
	GetIDToken() (string, error)

	RefreshRoutine()
}

// RenewingTokensImpl is the implementation of the RenewingTokens interface using the existing
// GetServiceJWT method
type RenewingTokensImpl struct {
	mu sync.RWMutex

	tokens    *Tokens
	lastError error

	authService string
	interval    time.Duration
}

// GetAccessToken gets the current access token or returns the last error
func (i *RenewingTokensImpl) GetAccessToken() (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.lastError != nil {
		return "", i.lastError
	}

	if i.tokens == nil {
		return "", errors.New("No access token is available")
	}

	return i.tokens.AccessToken, nil
}

// GetIDToken gets the current ID token or returns the last error
func (i *RenewingTokensImpl) GetIDToken() (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.lastError != nil {
		return "", i.lastError
	}

	if i.tokens == nil {
		return "", errors.New("No ID token is available")
	}

	return i.tokens.IDToken, nil
}

// RefreshRoutine is a never ending method that will periodically refresh the tokens
func (i *RenewingTokensImpl) RefreshRoutine() {
	ticker := time.NewTicker(i.interval)
	logrus.Infof("Starting Service JWT Refresh routine. Interval: %s", i.interval)

	for {
		select {
		case <-ticker.C:
			i.mu.Lock()
			i.tokens, i.lastError = GetServiceJWT(i.authService)
			if i.lastError != nil {
				logrus.Infof("Failed to fefresh Service JWT: %s", i.lastError.Error())
			} else {
				logrus.Infof("Succesfully refreshed Service JWT")
			}
			i.mu.Unlock()
		}
	}
}

// GetRenewingServiceJWT returns a RewnewingTokensImpl that will perodically refresh
// the stored tokens.
// Note: it is up to the caller to start the refresh routine
func GetRenewingServiceJWT(authService string, interval time.Duration) RenewingTokens {
	impl := &RenewingTokensImpl{
		authService: authService,
		interval:    interval,
	}

	impl.tokens, impl.lastError = GetServiceJWT(authService)

	return impl
}
