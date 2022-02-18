package api

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	yaml "github.com/ghodss/yaml"
	"github.com/gigantum/hoss-auth/pkg/test"
	"github.com/mitchellh/go-homedir"
)

// SetupLoadConfig is a simple fixture that loads the .env and config files
func SetupLoadConfig(t *testing.T) (*Settings, error) {
	err := test.LoadEnvFile("~/.hoss/.env")
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	configPath := filepath.Join("/opt", "config.yaml")
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// Config does not exist, so you are running the test locally.
		configPath, err = homedir.Expand("~/.hoss/auth/config.yaml")
		if err != nil {
			t.Fatal(err)
		}
	}

	settingsBytes, err := ioutil.ReadFile(configPath)
	if err == nil {
		if err = yaml.Unmarshal(settingsBytes, &settings); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}

	return settings, nil
}
