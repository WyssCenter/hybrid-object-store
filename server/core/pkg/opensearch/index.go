package opensearch

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const metadataIndexMappings = `
{
   "mappings":{
      "properties":{
         "core_service_endpoint": {"type": "keyword"},
		 "dataset_extended": {"type": "keyword"},
		 "object_key": {"type": "keyword"},
         "last_modified_date": {"type": "date"},
         "size_bytes": {"type": "double"},
		 "metadata": {
			 "type": "keyword", 
			 "normalizer": "lowercase",
			"fields": {
				"autocomplete": {
					"type": "completion",
					"contexts": [
						{
							"name": "dataset",
							"type": "category",
							"path": "dataset_extended"
						}
					]
				}
			}
		}
      }
   }
 }
`

type MetadataIndexPayload struct {
	// CoreServiceEndpoint is the core service root (e.g. http://localhost/core/v1)
	CoreServiceEndpoint string `json:"core_service_endpoint"`
	// DatasetExtended is a compound string that uniquely identifies a dataset (<object store name>|<bucket name>|<dataset name>)
	DatasetExtended string `json:"dataset_extended"`
	// Object key is the object key in the bucket
	ObjectKey string `json:"object_key"`
	// LastModifiedDate is the datetime string indicating the last modified date in UTC
	LastModifiedDate string `json:"last_modified_date"`
	// SizeBytes is the size of the object in bytes
	SizeBytes int `json:"size_bytes"`
	// Metadata is a list of strings representing key-pairs separated by ':' (e.g. ["fizz:buzz"])
	Metadata []string `json:"metadata"`
}

// searchIndexExists Returns true if the specified index exists and false if it does not
func searchIndexExists(opensearchEndpoint string, indexName string) (bool, error) {
	req, err := http.NewRequest("HEAD", opensearchEndpoint+"/"+indexName, nil)
	if err != nil {
		return false, errors.Wrap(err, "could not create request to check if index exists")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "could not make request to check if index exists")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	} else if resp.StatusCode == 404 {
		return false, nil
	} else {
		return false, errors.New("Unexpected response while checking if search index exists")
	}
}

// CreateMetadataSearchIndex checks if the metadata search index exists, if not, creates it
func CreateMetadataSearchIndex(opensearchEndpoint string) error {
	indexExists, err := searchIndexExists(opensearchEndpoint, "metadata-index")
	if err != nil {
		return err
	}

	if !indexExists {
		logrus.Info("Initializing 'metadata-index' opensearch index...")
		req, err := http.NewRequest("PUT", opensearchEndpoint+"/metadata-index", bytes.NewBuffer([]byte(metadataIndexMappings)))
		if err != nil {
			return errors.New("could not create metadata index request: " + err.Error())
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errors.New("could not make metadata index request: " + err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			d, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return errors.New("failed to parse create metadata index response: " + err.Error())
			}
			return errors.New(fmt.Sprintf("problem with metadata index response: StatusCode != 200: %s", string(d)))
		}
	}

	return nil
}
