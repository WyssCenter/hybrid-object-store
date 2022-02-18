package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-sync/pkg/config"
)

// ApiSyncNotification is the struct defining a single API sync event
type ApiSyncNotification struct {
	EventType         string `json:"event_type"`
	SourceEndpoint    string `json:"source_endpoint"`
	Namespace         string `json:"namespace"`
	Dataset           string `json:"dataset"`
	Description       string `json:"description,omitempty"`
	Group             string `json:"group,omitempty"`
	Permission        string `json:"permission,omitempty"`
	ObjectStore       string `json:"object_store,omitempty"`
	TargetCoreService string `json:"target_core_service"`
	TargetNamespace   string `json:"target_namespace"`
	SyncPolicy        string `json:"sync_policy,omitempty"`

	HasReloaded bool `json:"-"` // Flag used so that RequireReload only returns true once
}

func (asn *ApiSyncNotification) String() string {
	return fmt.Sprintf("<ApiSyncNotification %s %s/%s>", asn.EventType, asn.Namespace, asn.Dataset)
}

// RequireReload returns true once if the API notification requires the latest sync configuration information
// It only returns true once so that the API notification message can be requeued without making an infinite loop
func (asn *ApiSyncNotification) RequireReload() bool {
	// Return true if dealing with a Namespace Duplex or Dataset create notification, both of which
	// require the latest sync configuration information so that Match() will match this message

	// Reload needs to happen when
	// - put-ds-sync: Required so that BucketNotifications are able to match against the
	//                SyncPrefixes field to determine if the notification should be synced
	//                Required so that Execute can find the target Namespace config, so that
	//                it is able to send the create dataset API call to the correct Core Service
	// - put-ns-duplex: Required so that Execute can find the target Namespace config, so
	//                  that it is able to send the enable duplex sync API call to the correct
	//                  Core Service
	if !asn.HasReloaded {
		asn.HasReloaded = true
		return asn.EventType == "put-ds-sync" || asn.EventType == "put-ns-duplex"
	} else {
		return false
	}
}

func (asn *ApiSyncNotification) findNamespace(populatedConfig *config.PopulatedCoreServiceConfiguration) *config.PopulatedNamespaceConfiguration {
	for _, namespace := range populatedConfig.Namespaces {
		if namespace.CoreService.Endpoint == asn.SourceEndpoint &&
			namespace.Name == asn.Namespace {
			return namespace
		}
	}
	return nil
}

// Match checks to see if the message should be passed on to the worker based on the provided config
func (asn *ApiSyncNotification) Match(populatedConfig *config.PopulatedCoreServiceConfiguration) (bool, bool) {
	populatedConfig.L.RLock()
	defer populatedConfig.L.RUnlock()

	if asn.EventType == "create-namespace" {
		// For the 'create-namespace' message, we simply match on the endpoints. This message indicates
		// a namespace was added and we should refresh the S3 client in the related object store.
		// The actual namespace may not be loaded into a populated config yet if sync has not been enabled
		// so searching for the namespace would fail.
		return populatedConfig.Endpoint == asn.SourceEndpoint, false
	} else {
		// NOTE: not matching on PopulatedSyncConfiguration.SourcePrefixes as API notifications are only
		//       emitted for Datasets that are syncing
		return asn.findNamespace(populatedConfig) != nil, false
	}
}

// Execute performs the required actions based on the message and populated configuration of the worker.
func (asn *ApiSyncNotification) Execute(populatedConfig *config.PopulatedCoreServiceConfiguration) {
	populatedConfig.L.RLock()
	defer populatedConfig.L.RUnlock()

	// Handle all sync targets in parallel but wait for the message to finish processing before returning

	// For the 'create-namespace' message, we for reloading the client for the related object store
	// and then return with no further processing.
	if asn.EventType == "create-namespace" {
		for _, objStore := range populatedConfig.ObjectStores {
			if objStore.Name == asn.ObjectStore {
				logrus.Infof("Due to namespace creation, forcing reload of the client for '%s'", objStore.Name)
				objStore.Client.ForceRefresh()
				break
			}
		}
		return
	}

	var wg sync.WaitGroup
	namespace := asn.findNamespace(populatedConfig)
	for _, target := range namespace.SyncTargets {
		// We want to match the Namespace Duplex event to the specific config,
		//not just any originating from the source of the API events
		if asn.EventType == "put-ns-duplex" {
			if target.Target.Name != asn.TargetNamespace ||
				target.Target.CoreService.Endpoint != asn.TargetCoreService {
				continue
			}
		}

		wg.Add(1)
		go func(t *config.SyncTarget) {
			defer wg.Done()
			err := asn.handleSync(namespace, t.Target, populatedConfig.SyncObjectQueue)
			if err != nil {
				logrus.Error(err)
			}
		}(target)
	}

	wg.Wait()
}

func (asn *ApiSyncNotification) handleSync(sourceNamespace,
	targetNamespace *config.PopulatedNamespaceConfiguration, objChan chan config.Message) error {
	var err error

	switch asn.EventType {
	case "put-ds-perm":
		// PUT a group permission to a dataset in the target
		path := fmt.Sprintf("/namespace/%s/dataset/%s/group/%s/access/%s", targetNamespace.Name, asn.Dataset, asn.Group, asn.Permission)
		err = asn.makeSyncApiRequest("PUT", path, nil, targetNamespace)
	case "delete-ds-perm":
		// DELETE a group permission from a dataset in the target
		path := fmt.Sprintf("/namespace/%s/dataset/%s/group/%s", targetNamespace.Name, asn.Dataset, asn.Group)
		err = asn.makeSyncApiRequest("DELETE", path, nil, targetNamespace)
	case "put-ds-sync":
		// Create a dataset in target
		var jsonBytes = []byte(fmt.Sprintf(`{"name":"%s", "description":"%s"}`, asn.Dataset, asn.Description))
		path := fmt.Sprintf("/namespace/%s/dataset/", targetNamespace.Name)
		err = asn.makeSyncApiRequest("POST", path, jsonBytes, targetNamespace)

		// Start goroutine to sync any data that already exists in the dataset
		go func(n *config.PopulatedNamespaceConfiguration, c chan config.Message) {
			logrus.Infof("Starting to populate sync queue with existing objects for '%s'", asn.Dataset)

			// Make sure the target has time to create the dataset before syncing begins
			time.Sleep(3 * time.Second)

			sourceClient, err := n.ObjectStore.Client.GetClient()
			if err != nil {
				logrus.Errorf("Failed to load S3 client while populating existing objects for sync: %v", err)
				return
			}

			// The prefix is the RootDirectory value from the Dataset. This is not directly available here
			// so we manually create it from the dataset name. If the RootDirectory ever changes from the
			// default value, changes would be required to make this work properly.
			prefix := asn.Dataset + "/"
			input := &s3.ListObjectsV2Input{
				Bucket: &n.BucketName,
				Prefix: &prefix,
			}

			for {
				response, err := sourceClient.ListObjectsV2(context.TODO(), input)
				if err != nil {
					logrus.Errorf("Failed to fetch objects while populating existing objects for sync: %v", err)
					return
				}

				for _, item := range response.Contents {
					var msg BucketNotificationRecord
					msg.EventName = "s3:ObjectCreated:Put"
					msg.EventTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")
					msg.S3.Bucket.Name = n.BucketName
					msg.S3.Object.Key = *item.Key
					msg.S3.Object.Size = int(item.Size)
					msg.Source.Host = asn.SourceEndpoint
					msg.Source.UserAgent = "sync/1"
					msg.Endpoint = n.ObjectStore.Endpoint
					c <- &msg
				}

				if !response.IsTruncated {
					logrus.Infof("Finished populating sync queue with existing objects for '%s'", asn.Dataset)
					return
				}

				input.ContinuationToken = response.NextContinuationToken
			}
		}(sourceNamespace, objChan)

	case "put-ds-duplex":
		// Enable duplex sync in the target dataset
		type syncInput struct {
			SyncType   string `json:"sync_type"`
			SyncPolicy string `json:"sync_policy"`
		}

		policy := syncInput{SyncType: config.DuplexSyncType, SyncPolicy: asn.SyncPolicy}
		jsonBytes, err := json.Marshal(policy)
		if err != nil {
			return err
		}

		path := fmt.Sprintf("/namespace/%s/dataset/%s/sync", targetNamespace.Name, asn.Dataset)
		err = asn.makeSyncApiRequest("PUT", path, jsonBytes, targetNamespace)
	case "put-ns-duplex":
		// Enable duplex sync in the target namespace
		var jsonBytes = []byte(fmt.Sprintf(`{"target_core_service":"%s","target_namespace":"%s","sync_type":"%s"}`, asn.SourceEndpoint, asn.Namespace, config.DuplexSyncType))
		path := fmt.Sprintf("/namespace/%s/sync", asn.TargetNamespace)
		err = asn.makeSyncApiRequest("PUT", path, jsonBytes, targetNamespace)
	default:
		return errors.New("Unhandled API Sync event type " + asn.String())
	}

	if err != nil {
		return err
	}

	return nil
}

// makeSyncApiRequest is a helper function to make a REST request to the target core service
func (asn *ApiSyncNotification) makeSyncApiRequest(verb string, path string, jsonBytes []byte,
	targetNamespace *config.PopulatedNamespaceConfiguration) error {
	client := &http.Client{}

	targetCoreService := targetNamespace.CoreService.Endpoint

	// Hack to support running on localhost
	targetCoreService = strings.Replace(targetCoreService, "localhost/core", "core:8080", 1)

	var req *http.Request
	var err error
	var expectedStatus []int
	switch verb {
	case "PUT":
		expectedStatus = []int{204}
		if jsonBytes == nil {
			req, err = http.NewRequest(http.MethodPut, targetCoreService+path, nil)
		} else {
			req, err = http.NewRequest(http.MethodPut, targetCoreService+path, bytes.NewBuffer(jsonBytes))
		}
		if err != nil {
			return err
		}

	case "POST":
		expectedStatus = []int{201, 403} // Created or Forbidden (Already Exists)
		req, err = http.NewRequest(http.MethodPost, targetCoreService+path, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return err
		}

	case "DELETE":
		expectedStatus = []int{204}
		req, err = http.NewRequest(http.MethodDelete, targetCoreService+path, nil)
		if err != nil {
			return err
		}
	default:
		return errors.New("Unsupported request type: " + verb)
	}

	idToken, err := targetNamespace.CoreService.Tokens.GetIDToken()
	if err != nil {
		return errors.Wrap(err, "could not get service ID Token for authentication")
	}

	req.Header.Set("Authorization", "Bearer "+idToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "exec-env/hoss-sync-service")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	for _, code := range expectedStatus {
		if code == resp.StatusCode {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Failed to make API sync request to target `%s`, Status Code %v",
		targetCoreService+path, resp.Status))
}
