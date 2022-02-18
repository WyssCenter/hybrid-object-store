package elastic

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/sirupsen/logrus"

	errors "github.com/gigantum/hoss-error"
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

func CreateMetadataSearchIndex(esService string) error {
	req, err := http.NewRequest("PUT", esService+"/metadata-index", bytes.NewBuffer([]byte(metadataIndexMappings)))
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

	if resp.StatusCode == 304 {
		return nil
	}

	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return errors.New("problem with metadata index response: " + err.Error())
		}

		if strings.Contains(string(d), "already exists") {
			logrus.Warning("Metadata search index already exists")
			return nil
		}

		logrus.Error(string(d))
		return errors.New("problem with metadata index response: StatusCode != 200")
	}

	return nil
}
