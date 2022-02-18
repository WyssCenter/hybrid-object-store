package sync

import (
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gin-gonic/gin"
)

const EVENT_PUT_NAMESPACE_DUPLEX = "put-ns-duplex"
const EVENT_PUT_DATASET_PERMS = "put-ds-perm"
const EVENT_DELETE_DATASET_PERMS = "delete-ds-perm"
const EVENT_PUT_DATASET_SYNC = "put-ds-sync"
const EVENT_PUT_DATASET_DUPLEX = "put-ds-duplex"
const EVENT_CREATE_NAMESPACE = "create-namespace"

type ApiEventMsg struct {
	EventType      string `json:"event_type"`
	SourceEndpoint string `json:"source_endpoint"`
	Namespace      string `json:"namespace"`
	Dataset        string `json:"dataset"`

	// Dataset Sync
	Description string `json:"description,omitempty"`

	// Dataset Perms
	Group      string `json:"group,omitempty"`
	Permission string `json:"permission,omitempty"`

	// Namespace Duplex
	TargetCoreService string `json:"target_core_service,omitempty"`
	TargetNamespace   string `json:"target_namespace,omitempty"`

	// Create Namespace
	// We must explicitly send the object store name because
	// where this information is needed in the sync service
	// you cannot lookup the object store because the
	// namespace config info might not be loaded
	ObjectStore string `json:"object_store,omitempty"`

	// Sync Policy is used when enabling duplex sync
	SyncPolicy string `json:"sync_policy,omitempty"`
}

func getApiSyncExchange(c *gin.Context, objStore store.ObjectStore) (ApiSyncExchange, error) {
	val, ok := c.Get("apiSyncExchange")
	if !ok {
		return nil, errors.New("Failed to load apiSyncExchange in context")
	}

	exchanges := val.(map[string]ApiSyncExchange)
	exchange, ok := exchanges[objStore.GetName()]
	if !ok {
		return nil, errors.New("No ApiSyncExchange for namespace")
	}

	return exchange, nil
}

// SyncPermissionsHandler is a function that will emit a message to sync
// the mutation of dataset permissions. The sync service will then reconcile
// if the message should trigger an actual sync operation and map to
// the correct target namespace & core service.
func SyncPermissionsHandler(c *gin.Context,
	objStore store.ObjectStore, namespace *database.Namespace,
	dataset *database.Dataset, group string, permission string) error {

	if dataset.SyncEnabled {
		var msg ApiEventMsg
		if c.Request.Method == "PUT" {
			msg = ApiEventMsg{
				EventType:      EVENT_PUT_DATASET_PERMS,
				SourceEndpoint: msgSourceEndpoint(),
				Namespace:      namespace.Name,
				Dataset:        dataset.Name,
				Group:          group,
				Permission:     permission,
			}

		} else if c.Request.Method == "DELETE" {
			msg = ApiEventMsg{
				EventType:      EVENT_DELETE_DATASET_PERMS,
				SourceEndpoint: msgSourceEndpoint(),
				Namespace:      namespace.Name,
				Dataset:        dataset.Name,
				Group:          group,
			}
		} else {
			// Not a request that is synced
			return nil
		}

		ase, err := getApiSyncExchange(c, objStore)
		if err != nil {
			return errors.Wrap(err, "Failed to get API sync exchange")
		}
		err = ase.SendMessage(&msg)
		if err != nil {
			return errors.Wrap(err, "Failed to publish api sync message")
		}
	}

	return nil
}

// SyncDatasetHandler is a function that will emit a message to sync
// the a dataset when syncing is enabled. The sync service will reconcile if the message should
// trigger an actual sync operation and map to the correct target
// namespace & core service, and lookup the required data via the source core service.
func SyncDatasetHandler(c *gin.Context,
	objStore store.ObjectStore, namespace *database.Namespace,
	dataset *database.Dataset, syncPolicy string, isUpdate bool) error {

	if dataset.SyncEnabled {
		var msg ApiEventMsg
		ase, err := getApiSyncExchange(c, objStore)
		if err != nil {
			return errors.Wrap(err, "Failed to get API sync exchange")
		}

		// Only create dataset in the target if setting sync for the first time
		if !isUpdate {
			if c.Request.Method == "PUT" {
				msg = ApiEventMsg{
					EventType:      EVENT_PUT_DATASET_SYNC,
					SourceEndpoint: msgSourceEndpoint(),
					Namespace:      namespace.Name,
					Dataset:        dataset.Name,
					Description:    dataset.Description,
				}
			} else {
				// Not a request that is synced
				return nil
			}

			err = ase.SendMessage(&msg)
			if err != nil {
				return errors.Wrap(err, "Failed to publish api sync message (dataset create)")
			}
		}

		if dataset.SyncType == database.SYNC_TYPE_DUPLEX {
			// Hack...should process within 2 seconds to ensure ordering
			time.Sleep(2 * time.Second)

			msg = ApiEventMsg{
				EventType:      EVENT_PUT_DATASET_DUPLEX,
				SourceEndpoint: msgSourceEndpoint(),
				Namespace:      namespace.Name,
				Dataset:        dataset.Name,
				SyncPolicy:     syncPolicy,
			}

			err = ase.SendMessage(&msg)
			if err != nil {
				return errors.Wrap(err, "Failed to publish api sync message (dataset enable duplex)")
			}
		}

		// Only set permissions if setting sync for the first time
		if !isUpdate {
			// Hack...should process within 2 seconds to ensure ordering
			time.Sleep(2 * time.Second)
			for _, perm := range dataset.Permissions {
				msg = ApiEventMsg{
					EventType:      EVENT_PUT_DATASET_PERMS,
					SourceEndpoint: msgSourceEndpoint(),
					Namespace:      namespace.Name,
					Dataset:        dataset.Name,
					Group:          perm.Group.GroupName,
					Permission:     perm.Permission,
				}

				err = ase.SendMessage(&msg)
				if err != nil {
					return errors.Wrap(err, "Failed to publish api sync message (dataset create permission)")
				}
			}
		}
	}

	return nil
}

// SyncNamespaceHandler is a function that will emit a message to duplex sync
// the a namespace when duplex syncing is enabled.
func SyncNamespaceHandler(c *gin.Context,
	objStore store.ObjectStore, namespace *database.Namespace, targetCoreService, targetNamespace string) error {
	var msg ApiEventMsg
	if c.Request.Method == "PUT" {
		msg = ApiEventMsg{
			EventType:         EVENT_PUT_NAMESPACE_DUPLEX,
			SourceEndpoint:    msgSourceEndpoint(),
			Namespace:         namespace.Name,
			TargetCoreService: targetCoreService,
			TargetNamespace:   targetNamespace,
		}
	} else {
		// Not a request that is synced
		return nil
	}

	ase, err := getApiSyncExchange(c, objStore)
	if err != nil {
		return errors.Wrap(err, "Failed to get API sync exchange")
	}
	err = ase.SendMessage(&msg)
	if err != nil {
		return errors.Wrap(err, "Failed to publish api sync message (namespace duplex)")
	}

	return nil
}

// CreateNamespaceHandler is a function that will emit a message to indicate a namespace has been created
// This will be used in the sync service to reload credentials so the service account has access to datasets
// created in the new namespace.
func CreateNamespaceHandler(c *gin.Context, objStore store.ObjectStore, namespace *database.Namespace) error {
	msg := ApiEventMsg{
		EventType:      EVENT_CREATE_NAMESPACE,
		SourceEndpoint: msgSourceEndpoint(),
		Namespace:      namespace.Name,
		ObjectStore:    namespace.ObjectStore.Name,
	}

	ase, err := getApiSyncExchange(c, objStore)
	if err != nil {
		return errors.Wrap(err, "Failed to get API sync exchange")
	}
	err = ase.SendMessage(&msg)
	if err != nil {
		return errors.Wrap(err, "Failed to publish api sync message (namespace creation)")
	}

	return nil
}

func msgSourceEndpoint() string {
	return os.Getenv("EXTERNAL_HOSTNAME") + "/core/v1"
}
