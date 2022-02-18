package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/sirupsen/logrus"

	errors "github.com/gigantum/hoss-error"
)

// NamespaceResponse defines the Namespace returned by the Core Service.
// Only fields that are used are defined.
type ObjectStore struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

// GetObjectStores queries the given core service for a list of object stores
func GetObjectStores(idToken, coreService string) ([]*ObjectStore, error) {
	// Hack to support running on localhost
	coreService = strings.Replace(coreService, "localhost/core", "core:8080", 1)

	req, err := http.NewRequest("GET", coreService+"/object_store/", nil)
	if err != nil {
		return nil, errors.New("could not create object store request: " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not make object store request: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with object store response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with object store response: StatusCode != 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("problem with reading object store response: " + err.Error())
	}

	var objectStoreDef []*ObjectStore
	if err := json.Unmarshal(body, &objectStoreDef); err != nil {
		return nil, errors.New("problem unmarshaling the object store response: " + err.Error())
	}

	return objectStoreDef, nil
}
