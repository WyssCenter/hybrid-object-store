package message

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-service/policy"
	"github.com/gigantum/hoss-sync/pkg/config"
)

// LookupPrefix determines if the given string starts with any of the given prefixes
func LookupPrefix(s string, prefixes map[string]string) string {
	for prefix, _ := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return prefix
		}
	}

	return ""
}

// BucketNotification is the struct defining the S3 Bucket Event messages being received
type BucketNotification struct {
	Records []BucketNotificationRecord `json:"Records"`
}

// BucketNotificationRecord is the struct defining a single S3 Bucket Event
type BucketNotificationRecord struct {
	EventName string `json:"eventName"`
	EventTime string `json:"eventTime"`
	S3        struct {
		Bucket struct {
			Name string `json:"name"`
		} `json:"bucket"`
		Object struct {
			Key  string `json:"key"`
			Size int    `json:"size"`
		} `json:"object"`
	} `json:"s3"`
	Source struct {
		Host      string `json:"host"`
		UserAgent string `json:"userAgent"`
	} `json:"source"`

	Endpoint string `json:"-"`
}

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

// RequireReload returns true if this message requires the latest sync configuration information to Match() and Execute()
func (bnr *BucketNotificationRecord) RequireReload() bool {
	return false
}

// FileBucket returns the name of the originating bucket name
func (bnr *BucketNotificationRecord) FileBucket() string {
	return bnr.S3.Bucket.Name
}

// FileKey returns the key of the originating file
func (bnr *BucketNotificationRecord) FileKey() string {
	key, err := url.QueryUnescape(bnr.S3.Object.Key)
	if err != nil {
		return bnr.S3.Object.Key
	} else {
		return key
	}
}

// FileDataset returns the root directory of the originating file
func (bnr *BucketNotificationRecord) FileDataset() string {
	return strings.Split(bnr.FileKey(), "/")[0]
}

// DatasetExtended returns the object key extended by the object store and bucket
func (bnr *BucketNotificationRecord) DatasetExtended(objStoreName string) string {
	return fmt.Sprintf("%s|%s|%s", objStoreName, bnr.FileBucket(), bnr.FileDataset())
}

// FileSize returns the file size of the originating file
func (bnr *BucketNotificationRecord) FileSize() int {
	return bnr.S3.Object.Size
}

// FileModifiedTime returns the time the event occurred
func (bnr *BucketNotificationRecord) FileModifiedTime() string {
	return bnr.EventTime
}

// FileOperation returns the event name
func (bnr *BucketNotificationRecord) FileOperation() string {
	return bnr.EventName
}

// FileID returns an ID composed of identifying elements of the file
func (bnr *BucketNotificationRecord) FileID(populatedConfig *config.PopulatedObjectStoreConfiguration) string {
	strID := fmt.Sprintf("%s|%s|%s", populatedConfig.CoreService.Endpoint, bnr.DatasetExtended(populatedConfig.Name), bnr.FileKey())
	return b64.StdEncoding.EncodeToString([]byte(strID))
}

// String returns the string representation of the notification
func (bnr *BucketNotificationRecord) String() string {
	return fmt.Sprintf("<BucketNotification %s %s/%s>", bnr.FileOperation(), bnr.FileBucket(), bnr.FileKey())
}

func (bnr *BucketNotificationRecord) findObjectStore(populatedConfig *config.PopulatedCoreServiceConfiguration) *config.PopulatedObjectStoreConfiguration {
	for _, objStore := range populatedConfig.ObjectStores {
		if objStore.Endpoint == bnr.Endpoint {
			return objStore
		}
	}

	return nil
}

func (bnr *BucketNotificationRecord) findNamespace(populatedConfig *config.PopulatedCoreServiceConfiguration) *config.PopulatedNamespaceConfiguration {
	for _, namespace := range populatedConfig.Namespaces {
		if namespace.BucketName == bnr.FileBucket() &&
			namespace.ObjectStore.Endpoint == bnr.Endpoint &&
			LookupPrefix(bnr.FileKey(), namespace.SyncPolicies) != "" {
			return namespace
		}
	}
	return nil
}

// Match checks to see if the message should be passed on to the worker based on the provided config
func (bnr *BucketNotificationRecord) Match(populatedConfig *config.PopulatedCoreServiceConfiguration) (bool, bool) {
	switch bnr.FileOperation() {
	case "s3:ObjectAccessed:Get",
		"s3:ObjectAccessed:Head",
		"ObjectAccessed:Get",
		"ObjectAccessed:Head":
		// We currently do not do anything for HEAD or GET operations in the sync service, so just return
		// isMatch=true and ignore=true so the event is immediately ignored. If in the future this changes and this check is removed so
		// events proceed to the Execute() function, you must be sure to properly protect against "echoing"
		// HEAD operations when duplex syncing is enabled.
		return true, true
	}

	populatedConfig.L.RLock()
	defer populatedConfig.L.RUnlock()

	// ignore dataset yaml files
	should_ignore := false
	if strings.HasSuffix(bnr.FileKey(), ".dataset.yaml") {
		logrus.Debug("Skipping notification for dataset yaml")
		should_ignore = true
	}

	return bnr.findObjectStore(populatedConfig) != nil, should_ignore
}

// Execute performs the required actions based on the message and populated configuration of the worker.
func (bnr *BucketNotificationRecord) Execute(populatedConfig *config.PopulatedCoreServiceConfiguration) {
	populatedConfig.L.RLock()
	defer populatedConfig.L.RUnlock()

	objStore := bnr.findObjectStore(populatedConfig)

	if bnr.FileOperation() == "ObjectRemoved:Delete" && bnr.FileSize() == 0 {
		// This is likely a Delete Marker being removed, moving a real object
		// into the latest version. This is a bit rare, so the solution here is
		// simple but inefficient. If fetching metadata succeeds, then we can
		// assume a file now exists here. Switch the type to `s3:ObjectCreated:Put`
		// to complete the restore process.
		client, err := objStore.Client.GetClient()
		if err != nil {
			logrus.Error(errors.Wrap(err, "unable to get objectstore client"))
			return
		}
		_, err = bnr.getObjectMetadata(client)
		if err == nil {
			bnr.EventName = "s3:ObjectCreated:Put"
			logrus.Infof("Object restore detected %s - %s", bnr.FileBucket(), bnr.FileKey())
		}
	}

	// Only fetch metadata if it is an event that writes data. All other events do not need metadata.
	// We fetch metadata here to minimize the number of duplicate HEAD requests required to process the event.
	var metadata map[string]string
	switch bnr.FileOperation() {
	case "s3:ObjectCreated:Put",
		"s3:ObjectCreated:Copy",
		"s3:ObjectCreated:CompleteMultipartUpload",
		"ObjectCreated:Put", // AWS doesn't include the 's3:' prefix
		"ObjectCreated:Copy",
		"ObjectCreated:CompleteMultipartUpload":

		// get metadata for the object from minio or s3
		client, err := objStore.Client.GetClient()
		if err != nil {
			logrus.Error(errors.Wrap(err, "unable to get objectstore client"))
			return
		}
		metadata, err = bnr.getObjectMetadata(client)
		if err != nil {
			logrus.Error(errors.Wrap(err, "unable to get metadata"))
			return
		}
	}

	// Handle all meta and sync targets in parallel but wait for the message to finish processing before returning
	var wg sync.WaitGroup

	wg.Add(1)
	go func(o *config.PopulatedObjectStoreConfiguration, m map[string]string) {
		defer wg.Done()
		err := bnr.handleMeta(o, m)
		if err != nil {
			logrus.Error(err)
		}
	}(objStore, metadata)

	// filter messages caused by the sync service. We do this for the sync handler, but we send
	// all messages through to the metadata handler to support multi-search index updating.
	// Uses AWS_EXECUTION_ENV (https://docs.aws.amazon.com/sdk-for-go/api/aws/corehandlers/)
	//   to add a custom suffix to the user agent
	if strings.Contains(bnr.Source.UserAgent, "exec-env/hoss-sync-service") {
		logrus.Debugf("Skipping notification caused by the sync service: %s", bnr)
	} else {
		namespace := bnr.findNamespace(populatedConfig)
		if namespace != nil {
			key := LookupPrefix(bnr.FileKey(), namespace.SyncPolicies)
			filter := namespace.SyncFilters[key]

			msgInfo := &policy.MessageInformation{
				EventOperation: bnr.FileOperation(),
				ObjectKey:      bnr.FileKey(),
				ObjectSize:     bnr.S3.Object.Size,
				ObjectMetadata: metadata,
			}
			passed, err := filter(msgInfo)
			if err != nil {
				logrus.Errorf("Cannot apply policy filter to message %s: %v", bnr.String(), err)
				// ??? should this fail open?
			} else if passed {
				for _, target := range namespace.SyncTargets {
					wg.Add(1)
					go func(t *config.SyncTarget, m map[string]string) {
						defer wg.Done()
						err := bnr.handleSync(namespace, t.Target, m)
						if err != nil {
							logrus.Error(err)
						}
					}(target, metadata)
				}
			}
		}
	}

	wg.Wait()
}

func (bnr *BucketNotificationRecord) handleSync(sourceNamespace,
	targetNamespace *config.PopulatedNamespaceConfiguration, metadata map[string]string) error {

	sourceClient, err := sourceNamespace.ObjectStore.Client.GetClient()
	if err != nil {
		return err
	}

	targetClient, err := targetNamespace.ObjectStore.Client.GetClient()
	if err != nil {
		return err
	}

	targetBucket := targetNamespace.BucketName

	switch bnr.FileOperation() {
	case "s3:ObjectCreated:Put",
		"s3:ObjectCreated:Copy",
		"s3:ObjectCreated:CompleteMultipartUpload",
		"ObjectCreated:Put", // AWS doesn't include the 's3:' prefix
		"ObjectCreated:Copy",
		"ObjectCreated:CompleteMultipartUpload":
		logrus.Infof("Processing Sync Event: %s", bnr)

		tmpfile, err := ioutil.TempFile("", "sync-file-data-*")
		if err != nil {
			return errors.Wrap(err, "Couldn't create temp file to sync "+bnr.String())
		}

		if err := bnr.Download(sourceClient, bnr.FileBucket(), bnr.FileKey(), tmpfile); err != nil {
			// cleanup the temp file is there is an error
			tmpfile.Close()
			os.Remove(tmpfile.Name())
			return errors.Wrap(err, "Couldn't download file "+bnr.String())
		} else {
			if err := bnr.Upload(targetClient, targetBucket, bnr.FileKey(), tmpfile, metadata); err != nil {
				// cleanup the temp file is there is an error
				tmpfile.Close()
				os.Remove(tmpfile.Name())
				return errors.Wrap(err, "Couldn't upload file "+bnr.String())
			}
		}

		if err := tmpfile.Close(); err != nil {
			return errors.Wrap(err, "Couldn't close temp file used to sync "+bnr.String())
		}

		if err := os.Remove(tmpfile.Name()); err != nil {
			return errors.Wrap(err, "Couldn't remove temp file used to sync "+bnr.String())
		}
	case "s3:ObjectRemoved:Delete",
		"ObjectRemoved:Delete",
		"ObjectRemoved:DeleteMarkerCreated",
		"s3:ObjectRemoved:DeleteMarkerCreated":
		logrus.Infof("Processing Sync Event: %s", bnr)

		if err := bnr.Delete(targetClient, targetBucket, bnr.FileKey()); err != nil {
			return errors.Wrap(err, "Couldn't delete file "+bnr.String())
		}
	case "s3:ObjectAccessed:Get",
		"s3:ObjectAccessed:Head",
		"ObjectAccessed:Get",
		"ObjectAccessed:Head":
		// Operations that don't have an action to take. These currently are dropped in the Match method.
		// If you need to do something with them, you'll have to allow them through the Match method.
	default:
		return errors.New("Unhandled " + bnr.String())
	}

	return nil
}

// Download uses the AWS Download Manager to download a file locally
func (bnr *BucketNotificationRecord) Download(client *s3.Client, bucket, key string, tmpfile *os.File) error {
	downloader := manager.NewDownloader(client)
	_, err := downloader.Download(context.TODO(), tmpfile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

// Upload uses the AWS Upload Manager to upload a local file
func (bnr *BucketNotificationRecord) Upload(client *s3.Client, bucket, key string, tmpfile *os.File, metadata map[string]string) error {
	uploader := manager.NewUploader(client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     tmpfile,
		Metadata: metadata,
	})

	return err
}

// Delete removes the given file from the object store
func (bnr *BucketNotificationRecord) Delete(client *s3.Client, bucket, key string) error {
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func (bnr *BucketNotificationRecord) getObjectMetadata(client *s3.Client) (map[string]string, error) {
	metadata := map[string]string{}

	// query metadata
	bucket := bnr.FileBucket()
	key := bnr.FileKey()
	metadataOutput, err := client.HeadObject(
		context.TODO(),
		&s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		return metadata, errors.Wrap(err, "error requesting object metadata")
	}

	return metadataOutput.Metadata, nil
}

func (bnr *BucketNotificationRecord) handleMeta(populatedConfig *config.PopulatedObjectStoreConfiguration, metadata map[string]string) error {
	// update metadata index depending on which operation the file has experienced
	switch bnr.FileOperation() {
	case "s3:ObjectCreated:Put",
		"s3:ObjectCreated:Copy",
		"s3:ObjectCreated:CompleteMultipartUpload",
		"ObjectCreated:Put", // AWS doesn't include the 's3:' prefix
		"ObjectCreated:Copy",
		"ObjectCreated:CompleteMultipartUpload":
		logrus.Infof("Processing Metadata Event: %s", bnr)

		formattedMetadata := []string{}
		for key, value := range metadata {
			formattedMetadata = append(formattedMetadata, fmt.Sprintf("%s:%s", key, value))
		}

		// fill payload with object metadata, size, metadata
		payload := MetadataIndexPayload{
			ObjectKey:           bnr.FileKey(),
			DatasetExtended:     bnr.DatasetExtended(populatedConfig.Name),
			CoreServiceEndpoint: populatedConfig.CoreService.Endpoint,
			LastModifiedDate:    bnr.FileModifiedTime(),
			SizeBytes:           bnr.FileSize(),
			Metadata:            formattedMetadata,
		}

		err := bnr.makeMetadataIndexRequest(populatedConfig, "PUT", &payload)
		if err != nil {
			return errors.New("could not add or update object in metadata index: " + err.Error())
		}

	case "s3:ObjectRemoved:Delete",
		"ObjectRemoved:Delete",
		"ObjectRemoved:DeleteMarkerCreated",
		"s3:ObjectRemoved:DeleteMarkerCreated":
		logrus.Infof("Processing Metadata Event: %s", bnr)
		// fill payload with object metadata, size, metadata
		emptyMeta := []string{}
		payload := MetadataIndexPayload{
			ObjectKey:           bnr.FileKey(),
			DatasetExtended:     bnr.DatasetExtended(populatedConfig.Name),
			CoreServiceEndpoint: populatedConfig.CoreService.Endpoint,
			LastModifiedDate:    "",
			SizeBytes:           0,
			Metadata:            emptyMeta,
		}

		err := bnr.makeMetadataIndexRequest(populatedConfig, "DELETE", &payload)
		if err != nil {
			return errors.New("could not remove object from metadata index: " + err.Error())
		}
	case "s3:ObjectAccessed:Get",
		"s3:ObjectAccessed:Head",
		"ObjectAccessed:Get",
		"ObjectAccessed:Head":
		// Operations that don't have an action to take. These currently are dropped in the Match method.
		// If you need to do something with them, you'll have to allow them through the Match method.
	default:
		// ignore get, head, other operations that don't change the file
	}

	return nil
}

// makeMetadataIndexRequest is a helper function to make a REST request to the elasticsearch service
func (bnr *BucketNotificationRecord) makeMetadataIndexRequest(populatedConfig *config.PopulatedObjectStoreConfiguration,
	verb string, payload *MetadataIndexPayload) error {

	client := &http.Client{}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "unable to marshal payload JSON")
	}

	// Hack to support running on localhost
	coreService := strings.Replace(populatedConfig.CoreService.Endpoint, "localhost/core", "core:8080", 1)
	path := coreService + "/search/document/metadata"

	var req *http.Request
	switch verb {
	case "PUT":
		req, err = http.NewRequest(http.MethodPut, path, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return err
		}

	case "DELETE":
		req, err = http.NewRequest(http.MethodDelete, path, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return err
		}
	default:
		return errors.New("Unsupported request type: " + verb)
	}

	token, err := populatedConfig.CoreService.Tokens.GetIDToken()
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "exec-env/hoss-sync-service")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 204 {
		return nil
	}

	d, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return errors.New("problem with metadata index response: " + err.Error())
	}
	logrus.Error(string(d))

	return errors.New(fmt.Sprintf("Failed to make metadata update request to target `%s`, Status Code %v",
		path, resp.Status))
}
