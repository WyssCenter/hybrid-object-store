# Metadata Search Design
The metadata search system is based on Opensearch. Endpoints in the core service interact with Opensearch locally (it is not exposed via the external network) to get, add and remove documents from the search index, search metadata, and autocomplete keys and values. The core service also initializes the search index at boot if it does not exist.


## Schema

* `core_service_endpoint (keyword):` the uri of the core service that controls this object
* `dataset_extended (keyword):` the path to the dataset, extended in the format `"<object-store-name>|<bucket-name>|<dataset-name>"`
* `object_key (keyword):` the object key
* `last_modified_date (date):` the last modified date for the object, approximated using the `EventTime` field from the bucket notification record
* `size_bytes <double>:` the size of the object
* `metadata <keyword>:` a list of strings representing the key-value pairs associated with the object, in the format `"<key>:<value>"`. These are stored as one string for simpler searching and autocomplete
    * `autocomplete (completion):` [FIELD] a field of the `metadata` property, which can be queried for autocomplete suggestions. This field has a dataset context, meaning the `dataset_extended` property must be defined for any autocomplete searches and suggestions will only be provided from within the specified dataset. More info here: [https://www.elastic.co/guide/en/elasticsearch/reference/5.5/suggester-context.html](https://www.elastic.co/guide/en/elasticsearch/reference/5.5/suggester-context.html) 

    * NOTE: the storage format of metadata means that certain elasticsearch features (currently unused by Hoss) may not work as intended. For example, the term suggester provides suggestions based on edit distance to help users find terms despite misspellings, but if the key and value are stored in one string then the suggester would only be able to suggest key-value pairs, not individual keys or values. 
* NOTE: elasticsearch doesn't differentiate between a list-version of a property type and a single-item version, so any property can become a list of values if provided with a list of values. This is why the `metadata` property is defined as a `keyword` property and not a list of `keywords`.

* NOTE: elasticsearch treats `text` fields differently from `keyword` fields, but our usage currently matches `keyword` better. `text` fields go through additional indexing for each individual "word" (including sections of words delineated by punctuation) so a search on a partial form of a string might match many different larger strings


## Search endpoints

These endpoints can be used for internal testing or further development. The opensearch API is not exposed publicly, so all requests should be wrapped in the core service's API.

Create new index:


```
PUT http://opensearch:9200/<index-name>
Data: 
{
   "mappings": {
      "properties": {
	   <field>: {"type": <field-type>}
      }
   }
}
```


Add or update item in index:


```
PUT http://opensearch:9200/<index-name>/_doc/<item-id>
Data: 
{
   <field>: <value>
}
```


Delete an item in index:


```
DELETE http://opensearch:9200/<index-name>/_doc/<item-id>
```


Search an index:


```
GET http://opensearch:9200/<index-name>/_search 
Data:
{
   "size": <number of items to return>,
   "from": <number to offset start of search>,

   // queries work by assigning "scores" to items based on how well they
   // match the search terms, and returning items in order of their scores
   "query": {

	// a bool query is used to combine multiple sub-queries
      "bool": {

	   // a must query only returns items that match all of its terms
         "must": {

		// key-value pairs are stored together and must be formatted into 
    // a string for searching
            "terms": {
               "metadata": ["<key>:<value>"]
            },

		// this field is only set if no key-value pairs are provided so 
		// the search returns all items. To switch it to return no items
		// in this case, just remove this field.
            "match_all": {
               "boost": 1.0
		}
	   },

	   // a filter query further refines the returned items without 
   // affecting their overall score 
	   "filter": {
		"bool" {
		   "must": {
			"term" {
			   "core_service_endpoint": "<core-service-endpoint>"
			},
		   },

		   // a should query is used to return items if they fulfill at
       // least one of the terms. When paired with a must query, an 
       // item can fulfill none of the should terms and still be 
           // returned if it fulfills the must query, therefore an 
       // additional field "minimum_should_match" can be used to 
       // define the minimum number of terms a should query must 
       // match for the item to be returned.
		   "should": {
			[
			   {
				"term": {
				   "dataset_extended": "<dataset-extended-path>"
				}
			   }
			]
		   },
		   "minimum_should_match": 1
		}
	   }
	}
   }
}
```


Get autocomplete suggestions:


```
GET http://opensearch:9200/<index-name>/_search 
Data:
{
   "_source": <property to search for suggestions>,

   "suggest": {

	// a name for the suggest query, can be anything
      "tag_suggest": {

	   "prefix": <prefix to find completed suggestions for>,

         "completion": {


		// the full path to search, including the field of the property
		// if necessary, in the format <property>.<field>
		"field": <field to search for suggestions>,

		"size": <number of items to return>,

		// list of contexts must match the contexts defined in the
		// property mapping. These limit the suggestion search to a
		// subset of the index.
            "contexts": {
               "<context name>": "<contex definition>"
            }
	   }
	}
   }
}
```



## Integration with Core Service

Previously, a dataset only had bucket events enabled once it was sync enabled. To allow the sync service to add and update records in the metadata index, the core service was updated to set up bucket events for all datasets at the time they're created and only stop sending bucket events once the dataset is deleted.

Endpoints in the core service were added to interface with Opensearch locally (it is not exposed via the external network) to get, add and remove documents from the search index, search metadata, and autocomplete keys and values. The core service also initializes the search index at boot if it does not exist.

Details on how to use these endpoint can be found in the core service API documentation (`/core/v1/swagger/index.html`), but are also discussed below.

### Search

A new endpoint (`/search`) was added to the core service API for searching the metadata index. The endpoint takes query params described below and returns a result:

`metadata=<key>:<value>,<key>:<value>,...` - comma-delimited list of metadata key-value pairs to match on

`size=<int>` - number of items to return, defaults to 25

`from=<int>` - number to offset the start of the search (used with `size` when paging results)

`namespace=<str>` - namespace to search within. If omitted, the entire index is searched.

`dataset=<str>` - dataset to search within (if specified, a namespace must also be specified)

`modified_after=<iso-format-date>` - lower bound for the last modified date of the object

`modified_before=<iso-format-date>` - upper bound for the last modified date of the object

This endpoint searches only for objects controlled by the core service hosting it. It also filters results based on the datasets that the requesting user has at least read permissions for, so they do not see objects they don't have permission to access. Currently, if a request is made with no key-value pairs, the endpoint returns all objects the user has access to, however this could be updated to return no objects instead.

### Get a Document
This endpoint (`​/search​/namespace​/{namespaceName}​/dataset​/{datasetName}​/metadata`) was added to specificity address an issue when running with S3. Due to CORS issues, you can't get S3 to allow any header, so you can't fetch arbitrary metadata for an object via a HEAD request. This endpoint is used by the UI when running on an S3 object store to load the metadata in the side panel when a user clicked on an object.

### Suggest Metadata Keys
This endpoint (`​/search​/namespace​/{namespaceName}​/dataset​/{datasetName}​/key`) is used by the UI and client libraries to auto-complete keys available in a dataset. You provide via a query string with the prefix of the key and the service returns a list of available keys. In the UI, this is used to populate dropdown boxes where you are able to enter keys.


### Suggest Metadata Values
This endpoint (`​/search​/namespace​/{namespaceName}​/dataset​/{datasetName}​/key/{key}/value`) is used by the UI and client libraries to auto-complete values available for a given key, in a dataset. You provide via a query string the prefix of the value and the service returns a list of available values. In the UI, this is used to populate dropdown boxes where you are able to enter values.

### (Service Account Only) Create or update a document

A PUT to this endpoint (`​/search​/document/metadata`) is used by the sync service to create or update documents when objects are written. It is **only** available to the service account. All other authorized users will be rejected.

### (Service Account Only) Delete a document

A DELETE to this endpoint (`​/search​/document/metadata`) is used by the sync service to delete documents when objects are deleted. It is **only** available to the service account. All other authorized users will be rejected.



## Integration with Sync Service

The sync service performs the task of creating and updating the metadata-index in addition to dataset and API syncing. At start-up, the sync service waits for the core service to become available, to ensure opensearch is ready.

The majority of the logic for handling metadata in the sync service is located in `server/sync/pkg/message/bucket.go`. Here, a goroutine `handleMeta()` is used to handle the bucket event and make the required request to the core service. Object metadata is fetched once in the `Execute()` function to minimize HEAD requests during concurrent processing of events.


## Authentication

Currently no authentication is set up, and the service is protected by being only on the “internal” network so it can only be accessed from other service containers.

Authentication could be added to protect the search service once it’s upgraded to use multiple nodes. Authentication configuration can be found here: [https://opensearch.org/docs/latest/security-plugin/configuration/openid-connect/](https://opensearch.org/docs/latest/security-plugin/configuration/openid-connect/) 

The security config file should be updated to use an openid authenticator (in the “authc” section) using the wellknown url and the “role” and “username” claims from the JWT tokens. No extra step for authentication (the “authz” section) is necessary. Files can be found, edited, and mounted in these locations:


```
/usr/share/elasticsearch/plugins/opensearch/securityconfig/config.yml
/usr/share/elasticsearch/plugins/opensearch/securityconfig/roles-mapping.yml
```

With authentication, requests to the search service just need the `Authentication` header set with `Bearer <identity-token>`.
