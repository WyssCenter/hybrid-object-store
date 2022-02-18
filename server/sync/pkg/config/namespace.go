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
type NamespaceResponse struct {
	Name        string `json:"name"`
	ObjectStore struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	} `json:"object_store"`
	BucketName string `json:"bucket_name"`
}

// GetNamespace queries the given core service for the requested namespace data
func GetNamespace(idToken, coreService, namespace string) (*NamespaceResponse, error) {
	// Hack to support running on localhost
	coreService = strings.Replace(coreService, "localhost/core", "core:8080", 1)

	req, err := http.NewRequest("GET", coreService+"/namespace/"+namespace, nil)
	if err != nil {
		return nil, errors.New("could not create namespace request: " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not make namespace request: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with namespace response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with namespace response: StatusCode != 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("problem with reading namespace response: " + err.Error())
	}

	var namespaceDef NamespaceResponse
	if err := json.Unmarshal(body, &namespaceDef); err != nil {
		return nil, errors.New("problem unmarshaling the namespace response: " + err.Error())
	}

	return &namespaceDef, nil
}
