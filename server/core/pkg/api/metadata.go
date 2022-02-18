package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gigantum/hoss-core/pkg/opensearch"
	"github.com/gin-gonic/gin"
)

func getDatasetExtended(objStoreName string, bucketName string, rootDir string) string {
	return fmt.Sprintf(
		"%s|%s|%s",
		objStoreName,
		bucketName,
		strings.Replace(rootDir, "/", "", 1),
	)
}

type DatasetTerm struct {
	Term struct {
		DatasetExtended string `json:"dataset_extended"`
	} `json:"term"`
}

type MetadataTerms struct {
	Metadata []string `json:"metadata"`
}

type TimeRangeFilter struct {
	Range struct {
		LastModifiedDate map[string]string `json:"last_modified_date"`
	} `json:"range"`
}

type CoreServiceFilter struct {
	Term struct {
		CoreServiceEndpoint string `json:"core_service_endpoint"`
	} `json:"term"`
}

type MetadataSearchPayload struct {
	Size  int `json:"size"`
	From  int `json:"from"`
	Query struct {
		Bool struct {
			Must struct {
				Terms    *MetadataTerms `json:"terms,omitempty"`
				MatchAll map[string]int `json:"match_all,omitempty"`
			} `json:"must"`
			Filter struct {
				Bool struct {
					Must               []interface{} `json:"must"`
					Should             []DatasetTerm `json:"should"`
					MinimumShouldMatch int           `json:"minimum_should_match"`
				} `json:"bool"`
			} `json:"filter"`
		} `json:"bool"`
	} `json:"query"`
}

type MetadataSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ObjectKey        string   `json:"object_key"`
				DatasetExtended  string   `json:"dataset_extended"`
				LastModifiedDate string   `json:"last_modified_date"`
				SizeBytes        int      `json:"size_bytes"`
				Metadata         []string `json:"metadata"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type MetadataSearchResult struct {
	// URI is the Hoss URI that can be loaded via the client library
	URI string `json:"uri"`
	// FilePath is the full path in the object store to the file
	FilePath string `json:"file_path"`
	// Dataset is the dataset this object is in
	Dataset string `json:"dataset"`
	// Namespace is the namespace this object is in
	Namespace string `json:"namespace"`
	// LastModifiedDate is the datetime when this object was last modified
	LastModifiedDate string `json:"last_modified_date"`
	// SizeBytes is the size of the object in bytes
	SizeBytes int `json:"size_bytes"`
	// Metadata is a map of key-value pairs of metadata written to the object.
	Metadata []map[string]string `json:"metadata"`
}

// SearchMetadata searched the elasticsearch metadata index for the given key pairs
// @Summary Search object metadata
// @Schemes
// @Description Search object metadata based on key pairs or modified dates. The search process
// @Description will apply permissions to the results, only showing results in datasets to which
// @Description the authorized user has access. If no metadata key-value pairs are provided all
// @Description objects will be returned
// @Tags Search
// @Accept json
// @Produce json
// @Param	size  query  int  false  "Number of results to return" default(540)
// @Param	from  query  int  false  "Result index to start from if paging results" default(0)
// @Param	namespace  query  string  false  "If set, restrict results to this namespace"
// @Param	dataset  query  string  false  "If set, restrict results to this dataset. `namespace` must be set."
// @Param	metadata  query  string  false  "A comma separated list of key-value pairs to search for. (e.g. foo:bar,fizz:buzz)"
// @Param	modified_after  query  string  false  "Filter results to include only objects modified after the specified datetime string in the format '2006-01-02T15:04:05.000Z'"
// @Param	modified_before  query  string  false  "Filter results to include only objects modified before the specified datetime string in the format '2006-01-02T15:04:05.000Z'"
// @Success 200 {object} object{result=[]MetadataSearchResult}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search [get]
func SearchMetadata(c *gin.Context) {
	config, db := getAppConfig(c)

	// get key-value pairs from query parameters
	queryParams := c.Request.URL.Query()

	// get pagination parameters
	size := 25
	from := 0
	var err error
	if sizeParam, ok := queryParams["size"]; ok {
		size, err = strconv.Atoi(sizeParam[0])
		if err != nil {
			HandleError(c, err)
			return
		}
	}
	if fromParam, ok := queryParams["from"]; ok {
		from, err = strconv.Atoi(fromParam[0])
		if err != nil {
			HandleError(c, err)
			return
		}
	}

	// check if user wants to search within a specific namespace or dataset
	var namespaceName string
	var datasetName string
	if namespaceParam, ok := queryParams["namespace"]; ok {
		namespaceName = namespaceParam[0]
	}
	if datasetParam, ok := queryParams["dataset"]; ok {
		if namespaceName == "" {
			HandleError(c, errors.New("must specify namespace if searching within a dataset"))
			return
		}
		datasetName = datasetParam[0]
	}

	// get user's dataset permissions
	userInfo := getUserInfo(c)
	datasetRoots := []DatasetTerm{}
	datasetNamespaces := map[string]string{}
	limit := 10
	offset := 0
	for {
		objStores, err := db.ListObjectStores(limit, offset)
		if err != nil {
			HandleError(c, err)
			return
		}
		if len(objStores) == 0 {
			break
		}
		for _, objStore := range objStores {
			perms, err := db.GetPermissionsByUser(objStore, userInfo.Username, false)
			if err != nil {
				HandleError(c, err)
				return
			}

			for _, perm := range perms {
				// if searching within a specific namespace or dataset, only add the requested datasets
				if namespaceName != "" {
					if perm.Dataset.Namespace.Name != namespaceName {
						continue
					}
					if datasetName != "" {
						if perm.Dataset.Name != datasetName {
							continue
						}
					}
				}

				datasetTerm := DatasetTerm{}
				datasetTerm.Term.DatasetExtended = getDatasetExtended(
					objStore.Name,
					perm.Dataset.Namespace.BucketName,
					perm.Dataset.RootDirectory,
				)
				datasetRoots = append(datasetRoots, datasetTerm)

				// update datasetNamespaces map
				datasetNamespaces[datasetTerm.Term.DatasetExtended] = perm.Dataset.Namespace.Name
			}
		}

		offset += limit
	}

	// create payload
	payload := MetadataSearchPayload{}
	payload.Size = size
	payload.From = from
	payload.Query.Bool.Filter.Bool.Should = datasetRoots
	payload.Query.Bool.Filter.Bool.MinimumShouldMatch = 1
	payload.Query.Bool.Filter.Bool.Must = []interface{}{}

	// add metadata key value pairs, or if none provided then return all objects
	if metadata, ok := queryParams["metadata"]; ok {
		metadataTerms := MetadataTerms{
			Metadata: strings.Split(metadata[0], ","),
		}
		payload.Query.Bool.Must.Terms = &metadataTerms
	} else {
		payload.Query.Bool.Must.MatchAll = map[string]int{"boost": 1.0}
	}

	// add core service filter
	coreServiceFilter := CoreServiceFilter{}
	coreServiceFilter.Term.CoreServiceEndpoint = getCoreServiceEndpoint()
	payload.Query.Bool.Filter.Bool.Must = append(
		payload.Query.Bool.Filter.Bool.Must,
		coreServiceFilter,
	)

	// add timerange filter if requested in query parameters
	layout := "2006-01-02T15:04:05.000Z"
	timeRangeQuery := TimeRangeFilter{}
	timeRangeQuery.Range.LastModifiedDate = make(map[string]string)
	var t1 time.Time
	var t2 time.Time
	if startTime, ok := queryParams["modified_after"]; ok {
		t1, err = time.Parse(layout, startTime[0])
		if err != nil {
			HandleError(c, err)
			return
		}
		timeRangeQuery.Range.LastModifiedDate["gte"] = startTime[0]
	}
	if endTime, ok := queryParams["modified_before"]; ok {
		t2, err = time.Parse(layout, endTime[0])
		if err != nil {
			HandleError(c, err)
			return
		} else if t1.After(t2) {
			HandleError(c, errors.New("the `modified_after` time must be before the `modified_before` time"))
			return
		}
		timeRangeQuery.Range.LastModifiedDate["lte"] = endTime[0]
	}

	if len(timeRangeQuery.Range.LastModifiedDate) > 0 {
		payload.Query.Bool.Filter.Bool.Must = append(
			payload.Query.Bool.Filter.Bool.Must,
			timeRangeQuery,
		)
	}

	// convert payload to bytes
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	// create request to search metadata index
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Server.ElasticsearchEndpoint+"/metadata-index/_search", bytes.NewBuffer(payloadBytes))
	if err != nil {
		HandleError(c, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	if resp.StatusCode != 200 {
		HandleError(c, fmt.Errorf("unable to query metadata index, status code = %d", resp.StatusCode))
		return
	}

	// unpack results
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		HandleError(c, err)
		return
	}
	response := MetadataSearchResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		HandleError(c, err)
		return
	}

	// reformat results to clean up and add URI
	results := []MetadataSearchResult{}
	for _, hit := range response.Hits.Hits {
		result := MetadataSearchResult{
			FilePath:         strings.Join(strings.Split(hit.Source.ObjectKey, "/")[1:], "/"),
			Dataset:          strings.Split(hit.Source.ObjectKey, "/")[0],
			Namespace:        datasetNamespaces[hit.Source.DatasetExtended],
			LastModifiedDate: hit.Source.LastModifiedDate,
			SizeBytes:        hit.Source.SizeBytes,
		}
		result.URI = fmt.Sprintf(
			"hoss+%s://%s:%s:%s/%s",
			strings.Split(os.Getenv("EXTERNAL_HOSTNAME"), "://")[0],
			strings.Split(os.Getenv("EXTERNAL_HOSTNAME"), "://")[1],
			result.Namespace,
			result.Dataset,
			result.FilePath,
		)
		metadata := []map[string]string{}
		for _, metadataPair := range hit.Source.Metadata {
			metadataPairSplit := strings.SplitN(metadataPair, ":", 2)
			metadata = append(metadata, map[string]string{metadataPairSplit[0]: metadataPairSplit[1]})
		}
		result.Metadata = metadata
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

type SuggestPayload struct {
	Source  string `json:"_source"`
	Suggest struct {
		TagSuggest struct {
			Prefix     string `json:"prefix"`
			Completion struct {
				Field    string            `json:"field"`
				Size     int               `json:"size"`
				Contexts map[string]string `json:"contexts"`
			} `json:"completion"`
		} `json:"tag_suggest"`
	} `json:"suggest"`
}

type SuggestResponse struct {
	Suggest struct {
		TagSuggest []struct {
			Options []struct {
				Text   string `json:"text"`
				Source struct {
					Metadata []string `json:"metadata"`
				} `json:"_source"`
			} `json:"options"`
		} `json:"tag_suggest"`
	} `json:"suggest"`
}

// make a suggest query to the metadata index
func makeSuggestQuery(prefix string, limit int, datasetExtended string, elasticUrl string) (SuggestResponse, error) {

	// construct query
	payload := SuggestPayload{}
	payload.Source = "metadata"
	payload.Suggest.TagSuggest.Prefix = prefix
	payload.Suggest.TagSuggest.Completion.Field = "metadata.autocomplete"
	payload.Suggest.TagSuggest.Completion.Size = limit
	payload.Suggest.TagSuggest.Completion.Contexts = map[string]string{"dataset": datasetExtended}

	// convert payload to bytes
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return SuggestResponse{}, err
	}

	// create request to search metadata index
	client := &http.Client{}
	req, err := http.NewRequest("GET", elasticUrl+"/metadata-index/_search", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return SuggestResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return SuggestResponse{}, err
	}
	if resp.StatusCode != 200 {
		return SuggestResponse{}, err
	}

	// unpack results
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SuggestResponse{}, err
	}
	response := SuggestResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return SuggestResponse{}, err
	}

	return response, nil
}

// SuggestKeys searches the elasticsearch metadata index for autocomplete suggestions from key prefixes
// @Summary Suggest metadata keys based on a prefix
// @Schemes
// @Description Search object metadata keys for a dataset, given a prefix. This can be used to find
// @Description keys in use or to auto-complete keys.
// @Tags Search
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName  	 path      string  true  "Dataset Name"
// @Param	limit  query  int  false  "Maximum number of keys to suggest" default(25)
// @Param	prefix  query  string  false  "The prefix of the desired key to auto-complete"
// @Success 200
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search/namespace/{namespaceName}/dataset/{datasetName}/key [get]
func SuggestKeys(c *gin.Context) {
	config, db := getAppConfig(c)

	// get dataset info
	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	datasetExtended := getDatasetExtended(namespace.ObjectStore.Name, namespace.BucketName, dataset.RootDirectory)

	// check that user has permissions for this dataset
	userInfo := getUserInfo(c)
	permUsers, err := db.GetUsersWithPermissionsToDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	permissionGranted := false
	for _, user := range permUsers {
		if user == userInfo.Username {
			permissionGranted = true
			break
		}
	}
	if !permissionGranted {
		HandleError(c, errors.New("user does not have permission to suggest tags within this dataset"))
	}

	// get query parameters
	queryParams := c.Request.URL.Query()
	prefix := ""
	limit := 25
	if prefixParam, ok := queryParams["prefix"]; ok {
		prefix = prefixParam[0]
	}
	if limitParam, ok := queryParams["limit"]; ok {
		limit, err = strconv.Atoi(limitParam[0])
		if err != nil {
			HandleError(c, err)
			return
		}
	}

	// query metadata index for key suggestions
	response, err := makeSuggestQuery(prefix, limit, datasetExtended, config.Server.ElasticsearchEndpoint)
	if err != nil {
		HandleError(c, err)
		return
	}

	// If no TagSuggest items then there were no results. This can happen when an index
	// has not been initialized on a server
	if len(response.Suggest.TagSuggest) == 0 {
		var empty []string
		c.JSON(http.StatusOK, gin.H{"keys": empty})
	}

	// collect and filter all tags from the returned docs in case any docs have multiple matches
	keyMap := make(map[string]bool)
	for _, option := range response.Suggest.TagSuggest[0].Options {
		for _, metadataStr := range option.Source.Metadata {
			metadataPair := strings.SplitN(metadataStr, ":", 2)
			key := metadataPair[0]
			if strings.HasPrefix(key, prefix) {
				keyMap[key] = true
			}
		}
	}

	// store unique keys
	var keys []string
	for key := range keyMap {
		keys = append(keys, key)
	}

	// return keys and key-value map
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// SuggestValues searches the elasticsearch metadata index for autocomplete suggestions from value prefixes
// @Summary Suggest a metadata key's values based on a prefix
// @Schemes
// @Description Search object metadata values for the provided key in a dataset, given a prefix. This can be used to find
// @Description values in use for that key, or to auto-complete values for a key.
// @Tags Search
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName  	 path      string  true  "Dataset Name"
// @Param	key  	 path      string  true  "Metadata key to use when searching for values"
// @Param	limit  query  int  false  "Maximum number of values to suggest" default(25)
// @Param	prefix  query  string  false  "The prefix of the desired value to auto-complete"
// @Success 200
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search/namespace/{namespaceName}/dataset/{datasetName}/key/{key}/value [get]
func SuggestValues(c *gin.Context) {
	config, db := getAppConfig(c)

	// get key
	key := c.Param("key")

	// get dataset info
	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	datasetExtended := getDatasetExtended(namespace.ObjectStore.Name, namespace.BucketName, dataset.RootDirectory)

	// check that user has permissions for this dataset
	userInfo := getUserInfo(c)
	permUsers, err := db.GetUsersWithPermissionsToDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	permissionGranted := false
	for _, user := range permUsers {
		if user == userInfo.Username {
			permissionGranted = true
			break
		}
	}
	if !permissionGranted {
		HandleError(c, errors.New("user does not have permission to suggest tags within this dataset"))
	}

	// get query parameters
	queryParams := c.Request.URL.Query()
	prefix := ""
	limit := 25
	if prefixParam, ok := queryParams["prefix"]; ok {
		prefix = prefixParam[0]
	}
	prefixCombined := fmt.Sprintf("%s:%s", key, prefix)
	if limitParam, ok := queryParams["limit"]; ok {
		limit, err = strconv.Atoi(limitParam[0])
		if err != nil {
			HandleError(c, err)
			return
		}
	}

	// query metadata index for key suggestions
	response, err := makeSuggestQuery(prefixCombined, limit, datasetExtended, config.Server.ElasticsearchEndpoint)
	if err != nil {
		HandleError(c, err)
		return
	}

	// If no TagSuggest items then there were no results. This can happen when an index
	// has not been initialized on a server
	if len(response.Suggest.TagSuggest) == 0 {
		var empty []string
		c.JSON(http.StatusOK, gin.H{"values": empty})
	}

	// collect and filter all tags from the returned docs in case any docs have multiple matches
	valueMap := make(map[string]bool)
	for _, option := range response.Suggest.TagSuggest[0].Options {
		for _, metadataStr := range option.Source.Metadata {
			metadataPair := strings.SplitN(metadataStr, ":", 2)
			val := metadataPair[1]
			if strings.HasPrefix(val, prefix) {
				valueMap[val] = true
			}
		}
	}

	// store unique keys
	var values []string
	for val := range valueMap {
		values = append(values, val)
	}

	// return keys and key-value map
	c.JSON(http.StatusOK, gin.H{"values": values})
}

// GetMetadata searches the elasticsearch metadata index for autocomplete suggestions from key prefixes
// @Summary Get metadata stored for an object in the search index
// @Schemes
// @Description Get the metadata key-value pairs that are stored in the search index for the specified object.
// @Description
// @Description The `objectKey` query arg is the "full" object key, url escaped. This mean the object key should
// @Description should include the prefix inside the bucket for the dataset, which typically is the name of the
// @Description dataset.
// @Description
// @Description Note, this endpoint exists due to CORS limitations on AWS S3 that prevent the file browser widget
// @Description from loading metadata. This endpoint is used in its place from the UI when running in S3. For
// @Description most other use cases, you shouldn't need this endpoint and can load the metadata directly from S3/minIO
// @Tags Search
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName  	 path      string  true  "Dataset Name"
// @Param	objectKey  	 query      string  true  "URL Escaped Full Object Key"
// @Success 200 {object} object{metadata=object{string}}}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search/namespace/{namespaceName}/dataset/{datasetName}/metadata [get]
func GetMetadata(c *gin.Context) {
	config, db := getAppConfig(c)

	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	key := c.Query("objectKey")

	decodedkey, err := url.QueryUnescape(key)
	if err != nil {
		HandleError(c, err)
		return
	}

	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	datasetExtended := getDatasetExtended(namespace.ObjectStore.Name, namespace.BucketName, dataset.RootDirectory)

	// check that user has permissions for this dataset
	userInfo := getUserInfo(c)
	permUsers, err := db.GetUsersWithPermissionsToDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}
	permissionGranted := false
	for _, user := range permUsers {
		if user == userInfo.Username {
			permissionGranted = true
			break
		}
	}
	if !permissionGranted {
		HandleError(c, errors.New("user does not have permission to access metadata in this dataset"))
	}

	payload := opensearch.MetadataIndexPayload{CoreServiceEndpoint: getCoreServiceEndpoint(),
		DatasetExtended: datasetExtended, ObjectKey: decodedkey}

	doc, err := opensearch.GetDocument(config.Server.ElasticsearchEndpoint, &payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	m := map[string]string{}
	for _, meta := range doc.Source.Metadata {
		parts := strings.SplitN(meta, ":", 2)
		m[parts[0]] = parts[1]
	}

	c.JSON(http.StatusOK, gin.H{"metadata": m})
}
