package api

import (
	"io/ioutil"
	"os"
	"strings"
)

func getAvailableServices() []string {
	return strings.Fields(os.Getenv("AVAILABLE_SERVICES"))
}

func getVersion() (string, error) {
	v, err := ioutil.ReadFile("/opt/hoss-core/discover_version")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(v), "\n"), nil
}

func getBuildHash() (string, error) {
	b, err := ioutil.ReadFile("/opt/hoss-core/discover_build_hash")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(b), "\n"), nil
}
