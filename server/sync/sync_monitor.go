package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	errors "github.com/gigantum/hoss-error"
	service "github.com/gigantum/hoss-service"

	"github.com/gigantum/hoss-sync/pkg/config"
)

// SyncConfigurations is the set of SyncConfigurations returned by a single Core Service
type SyncConfigurations struct {
	mu     sync.RWMutex
	reload chan struct{}

	configs []config.SyncConfiguration
}

// ForceReload requests that the sync configuration information be polled outside of the normal interval
func (scs *SyncConfigurations) ForceReload() {
	scs.reload <- struct{}{}
}

// GetConfigs returns the current set of SyncConfigurations that have been queried from a Core Service
func (scs *SyncConfigurations) GetConfigs() []config.SyncConfiguration {
	scs.mu.RLock()
	defer scs.mu.RUnlock()
	return scs.configs
}

// Monitor periodically queries the given Core Service for the current set of SyncConfigurations stored in the Core Service's database
func (scs *SyncConfigurations) Monitor(tokens service.RenewingTokens, coreService string, interval time.Duration, notify chan<- struct{}) {
	lastChecked := time.Time{}
	ticker := time.NewTicker(interval)
	scs.reload = make(chan struct{})

	logrus.Infof("Starting to monitor %s for sync config changes (interval: %v)", coreService, interval)

	for {
		now := time.Now()

		configs, err := QuerySyncConfigurations(tokens, coreService, lastChecked)
		if err != nil {
			logrus.Warnf("Could not query sync configuration from %s: %s", coreService, err.Error())
		}

		lastChecked = now

		if configs != nil {
			scs.mu.Lock()
			scs.configs = configs
			scs.mu.Unlock()
			notify <- struct{}{}
		}

		select {
		case <-ticker.C:
		case <-scs.reload:
		}
	}
}

// QuerySyncConfigurations makes the HTTP query to the Core Service and decodes the results
func QuerySyncConfigurations(tokens service.RenewingTokens, coreService string, lastChecked time.Time) ([]config.SyncConfiguration, error) {
	// Hack to support running on localhost
	coreService = strings.Replace(coreService, "localhost/core", "core:8080", 1)

	req, err := http.NewRequest("GET", coreService+"/configuration/sync", nil)
	if err != nil {
		return nil, errors.New("could not create sync configuration request: " + err.Error())
	}

	idToken, err := tokens.GetIDToken()
	if err != nil {
		return nil, errors.New("could not get service ID Token for authentication: " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+idToken)
	req.Header.Set("If-Modified-Since", lastChecked.Format(time.RFC1123))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not make sync configuration request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 304 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with sync configuration response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with sync configuration response: StatusCode != 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("problem with reading sync configuration response: " + err.Error())
	}

	syncConfigurations := []config.SyncConfiguration{}
	if err := json.Unmarshal(body, &syncConfigurations); err != nil {
		return nil, errors.New("problem unmarshaling the sync configurations response: " + err.Error())
	}

	logrus.Debugf("New sync configs: %+v", syncConfigurations)

	return syncConfigurations, nil
}
