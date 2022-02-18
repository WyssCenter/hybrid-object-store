package opensearch

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// IndexSource is the source data for a search index document.
// Note, only a subset of the available fields are represented or used here.
type IndexSource struct {
	Metadata []string `json:"metadata"`
}

// DocumentGetResponse is the response from the Opensearch API when doing a GET
// on a document. Note, only the `_source` field is currently included or used.
// This is to provide direct fetching of an object's metadata stored the index.
type DocumentGetResponse struct {
	Source IndexSource `json:"_source"`
}

// getDocumentID returns the unique string used to identify a document in the search index
// the string is a combination of the core service endpoint, object store, bucket name, datset name
// and object key (including the prefix, which is typically the dataset name). This is then base64
// encoded
func getDocumentID(p *MetadataIndexPayload) string {
	strID := fmt.Sprintf("%s|%s|%s", p.CoreServiceEndpoint, p.DatasetExtended, p.ObjectKey)
	return b64.StdEncoding.EncodeToString([]byte(strID))
}

func CreateOrUpdateDocument(opensearchEndpoint string, documentPayload *MetadataIndexPayload) error {
	payloadBytes, err := json.Marshal(&documentPayload)
	if err != nil {
		return errors.Wrap(err, "unable to marshal payload JSON")
	}

	// update metadata search index
	path := opensearchEndpoint + "/metadata-index/_doc/" + getDocumentID(documentPayload)
	err = makeMetadataIndexRequest("PUT", path, payloadBytes)
	if err != nil {
		return errors.New("could not add or update object in metadata index: " + err.Error())
	}

	return nil
}

func DeleteDocument(opensearchEndpoint string, documentPayload *MetadataIndexPayload) error {
	// update metadata search index
	path := opensearchEndpoint + "/metadata-index/_doc/" + getDocumentID(documentPayload)
	err := makeMetadataIndexRequest("DELETE", path, nil)
	if err != nil {
		return errors.New("could not remove object from metadata index: " + err.Error())
	}

	return nil
}

func GetDocument(opensearchEndpoint string, documentPayload *MetadataIndexPayload) (*DocumentGetResponse, error) {
	path := opensearchEndpoint + "/metadata-index/_doc/" + getDocumentID(documentPayload)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("could not create request when fetching document: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not make request when fetching document: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// 404 means the document ID provided doesn't exist
		return nil, database.ErrNotFound
	} else if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to dump error response. Status code %v, err: %s",
				resp.StatusCode, err.Error()))
		}

		logrus.Debug(string(d))
		return nil, errors.New(fmt.Sprintf("failed to fetch search index document. Status code %v, ", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read response when fetching search index document: " + err.Error())
	}

	var docResponse DocumentGetResponse
	if err := json.Unmarshal(body, &docResponse); err != nil {
		return nil, errors.New("could not parse response when fetching search index document: " + err.Error())
	}

	return &docResponse, nil
}

// makeMetadataIndexRequest is a helper function to make a REST request to the opensearch service
func makeMetadataIndexRequest(verb string, path string, jsonBytes []byte) error {
	client := &http.Client{}

	var req *http.Request
	var err error
	var expectedStatus []int
	switch verb {
	case "PUT":
		expectedStatus = []int{201, 200}
		req, err = http.NewRequest(http.MethodPut, path, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return err
		}

	case "DELETE":
		expectedStatus = []int{200}
		req, err = http.NewRequest(http.MethodDelete, path, nil)
		if err != nil {
			return err
		}
	default:
		return errors.New("Unsupported request type: " + verb)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	for _, code := range expectedStatus {
		if code == resp.StatusCode {
			return nil
		}
	}

	d, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return errors.New("loading metadata index response after a failure resulted in an error: " + err.Error())
	}

	return errors.New(fmt.Sprintf("metadata update request to target `%s` failed, Status Code %v, Response: %s",
		path, resp.Status, string(d)))
}
