package config

import (
	"context"
	"sync"

	service "github.com/gigantum/hoss-service"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-sync/pkg/credentials"

	"github.com/gigantum/hoss-service/policy"
)

// SyncKey is the key for the PopulatedNamespaceConfiguration.SyncTargets map
// and defines the target Namespace for the sync
type SyncKey struct {
	CoreService string
	Namespace   string
}

// SyncTarget is the value for the PopulatedNamespaceConfiguration.SyncTargets map
// and defines the type of sync and contains a pointer to the target Namespace
type SyncTarget struct {
	SyncType string

	Target *PopulatedNamespaceConfiguration
}

// PopulatedObjectStoreConfiguration holds all of the information about an Object Store that this service needs to process messages
// Note: Object Store information is only queried once, at service startup
type PopulatedObjectStoreConfiguration struct {
	CoreService *PopulatedCoreServiceConfiguration

	Name     string
	Endpoint string
	Client   credentials.RenewingClient
}

// PopulatedNamespaceConfiguration holds all of the information about a Namespace that this service needs to process messages
// Note: Namespace sync information is periodically polled and updated
type PopulatedNamespaceConfiguration struct {
	CoreService *PopulatedCoreServiceConfiguration

	Name        string
	ObjectStore *PopulatedObjectStoreConfiguration
	BucketName  string

	SyncPolicies map[string]string
	SyncFilters  map[string]policy.PolicyFilter
	SyncTargets  map[SyncKey]*SyncTarget
}

// PopulatedCoreServiceConfiguration holds all of the information about a Core Service that this service needs to process messages
// Note: All Object Stores in the Core Service exist in the ObjectStores field
// Note: Only Namespaces that are configured as the source or target of a sync exist in the Namespaces field
type PopulatedCoreServiceConfiguration struct {
	L sync.RWMutex

	Tokens service.RenewingTokens

	Endpoint string // previously CoreService

	ObjectStores map[string]*PopulatedObjectStoreConfiguration
	Namespaces   map[string]*PopulatedNamespaceConfiguration

	WorkerQueue chan Message

	// SyncObjectQueue is a channel that the *worker* users to send messages to the demuxer.
	// When syncing is enabled on a dataset, the worker will list the dataset and push
	// any existing objects onto this channel for processing by the demuxer.
	SyncObjectQueue chan Message
}

// Worker is the go routine that will receive messages from the given Core Service
// and execute them
func (pcs *PopulatedCoreServiceConfiguration) Worker(ctx context.Context) {
	logrus.Infof("Worker starting for %s", pcs.Endpoint)
	for {
		select {
		case msg := <-pcs.WorkerQueue:
			msg.Execute(pcs)
		case <-ctx.Done():
			logrus.Infof("Stopping worker for %s", pcs.Endpoint)
			return
		}
	}
}
