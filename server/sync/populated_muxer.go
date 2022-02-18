package main

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	service "github.com/gigantum/hoss-service"
	"github.com/gigantum/hoss-service/policy"

	"github.com/gigantum/hoss-sync/pkg/config"
	"github.com/gigantum/hoss-sync/pkg/credentials"
)

func newNamespace(tokens service.RenewingTokens, coreService *config.PopulatedCoreServiceConfiguration, namespaceName string) *config.PopulatedNamespaceConfiguration {
	idToken, err := tokens.GetIDToken()
	if err != nil {
		logrus.Fatalf("Could not get service ID token: %s", err.Error())
	}

	resp, err := config.GetNamespace(idToken, coreService.Endpoint, namespaceName)
	if err != nil {
		logrus.Fatalf("Could not get namespace information: %s", err.Error())
	}

	namespace := &config.PopulatedNamespaceConfiguration{
		CoreService: coreService,
		Name:        namespaceName,
		ObjectStore: coreService.ObjectStores[resp.ObjectStore.Name],
		BucketName:  resp.BucketName,

		SyncPolicies: map[string]string{},
		SyncFilters:  map[string]policy.PolicyFilter{},
		SyncTargets:  map[config.SyncKey]*config.SyncTarget{},
	}

	return namespace
}

// PopulatedCoreServiceConfigurations contains the PopulatedCoreServiceConfigurations derived from each of the Core Services being monitored
type PopulatedCoreServiceConfigurations struct {
	reload         chan struct{}
	reloadFinished *sync.Cond

	populatedConfigs map[string]*config.PopulatedCoreServiceConfiguration
}

// ForceReload requests the Monitors to repoll their Core Service for any configuration changes
//
func (pcs *PopulatedCoreServiceConfigurations) ForceReload() {
	// Start a fail safe, so that this function doesn't block too long
	// There is a rare possibility that the regular polling cycle could pickup the changes
	// that were the reason for calling ForceReload, right as ForceReload is being called.
	// This could cause this function to block until the next sync configuration change in
	// the system.
	//
	// This go routine means that ForceReload will pause for up to one minute, but no longer
	go func() {
		time.Sleep(1 * time.Second)
		pcs.reloadFinished.Broadcast()
	}()

	pcs.reloadFinished.L.Lock()
	pcs.reload <- struct{}{}
	pcs.reloadFinished.Wait()
	pcs.reloadFinished.L.Unlock()
}

// GetConfigs returns the latest version of the PopulatedCoreServiceConfigurations that can be used to process incoming messages
func (pcs *PopulatedCoreServiceConfigurations) GetConfigs() map[string]*config.PopulatedCoreServiceConfiguration {
	return pcs.populatedConfigs
}

// UpdateMuxer waits for the CoreServiceConfigurations Monitors to notify it of a change in configuration and then works to reconcile the current statue with the new state
func (pcs *PopulatedCoreServiceConfigurations) UpdateMuxer(ctx context.Context, configuration *config.Configuration, tokens service.RenewingTokens) {
	logrus.Info("Starting core service configuration update muxer")

	pcs.populatedConfigs = make(map[string]*config.PopulatedCoreServiceConfiguration)
	pcs.reload = make(chan struct{})
	pcs.reloadFinished = sync.NewCond(&sync.Mutex{})
	notify := make(chan struct{})
	var configMonitors []*SyncConfigurations

	// Start the different monitors watching each core service for changes in the sync configurations
	idToken, err := tokens.GetIDToken()
	if err != nil {
		logrus.Fatalf("Could not get service ID token: %s", err.Error())
	}

	for _, coreService := range configuration.CoreServices {
		populatedCoreService := &config.PopulatedCoreServiceConfiguration{
			Tokens:   tokens,
			Endpoint: coreService,

			ObjectStores: map[string]*config.PopulatedObjectStoreConfiguration{},
			Namespaces:   map[string]*config.PopulatedNamespaceConfiguration{},

			WorkerQueue:     make(chan config.Message, configuration.WorkerBufferSize),
			SyncObjectQueue: make(chan config.Message, configuration.WorkerBufferSize),
		}
		pcs.populatedConfigs[coreService] = populatedCoreService

		// Populate Object Store information and
		// start STS Credential / S3 Client renewal routine
		objectStores, err := config.GetObjectStores(idToken, coreService)
		if err != nil {
			logrus.Fatalf("Could not get object store data: %s", err.Error())
		}

		for _, objStore := range objectStores {
			populatedObjectStore := &config.PopulatedObjectStoreConfiguration{
				CoreService: populatedCoreService,
				Name:        objStore.Name,
				Endpoint:    objStore.Endpoint,
				Client:      credentials.GetRenewingClient(tokens, coreService, objStore.Name, configuration.RefreshIntervals.StsCredentials),
			}
			populatedCoreService.ObjectStores[objStore.Name] = populatedObjectStore

			go populatedObjectStore.Client.RefreshRoutine(ctx)
		}

		// Start Sync Configuration change monitor
		configMonitor := &SyncConfigurations{}
		configMonitors = append(configMonitors, configMonitor)
		go configMonitor.Monitor(tokens, coreService, configuration.RefreshIntervals.CoreService, notify)

		// Create worker routines
		for i := 0; i < configuration.WorkerInstanceCount; i++ {
			go populatedCoreService.Worker(ctx)
		}
	}

	for {
		select {
		case <-pcs.reload:
			for _, monitor := range configMonitors {
				monitor.ForceReload()
			}
		case <-notify:
			logrus.Info("Received notice of sync configuration change")

			// Collect all of the changed configs
			toCreate := map[string]config.SyncConfiguration{}
			for _, monitor := range configMonitors {
				for _, syncConfig := range monitor.GetConfigs() {
					toCreate[syncConfig.Hash()] = syncConfig
				}
			}

			// Collect all of the current configs
			toDelete := map[string]config.SyncConfiguration{}
			for _, populatedCoreService := range pcs.populatedConfigs {
				for _, populatedNamespace := range populatedCoreService.Namespaces {
					for syncKey, syncTarget := range populatedNamespace.SyncTargets {
						syncConfig := config.SyncConfiguration{
							SyncType:          syncTarget.SyncType,
							SourceCoreService: populatedCoreService.Endpoint,
							SourceNamespace:   populatedNamespace.Name,
							SourcePolicies:    populatedNamespace.SyncPolicies,
							TargetCoreService: syncKey.CoreService,
							TargetNamespace:   syncKey.Namespace,
						}

						toDelete[syncConfig.Hash()] = syncConfig
					}
				}
			}

			// Remove the intersection of the two list, as those don't need to be changed
			for k := range toCreate {
				if _, ok := toDelete[k]; ok {
					delete(toCreate, k)
					delete(toDelete, k)
				}
			}

			// Update the config reference for the Demuxer
			// Figure out which Core Services we need to lock for updating the sync data
			toLock := map[string]*config.PopulatedCoreServiceConfiguration{}
			for _, syncConfig := range toDelete {
				coreService := pcs.populatedConfigs[syncConfig.SourceCoreService]
				toLock[syncConfig.SourceCoreService] = coreService
			}
			for _, syncConfig := range toCreate {
				coreService := pcs.populatedConfigs[syncConfig.SourceCoreService]
				toLock[syncConfig.SourceCoreService] = coreService
			}

			// Acquire write locks as we will now mutate the data structure
			for _, coreService := range toLock {
				coreService.L.Lock()
			}

			// Remove Sync Target information from the Namespace Configuration
			// Note: done first so that new / updated Sync Targets are not accidentally removed
			for _, syncConfig := range toDelete {
				logrus.Debugf("Deleting: %+v", syncConfig)
				coreService := pcs.populatedConfigs[syncConfig.SourceCoreService]

				namespace := coreService.Namespaces[syncConfig.SourceNamespace]

				syncKey := config.SyncKey{
					CoreService: syncConfig.TargetCoreService,
					Namespace:   syncConfig.TargetNamespace,
				}

				delete(namespace.SyncTargets, syncKey)
				if len(namespace.SyncTargets) == 0 {
					// If there are no SyncTargets remove the SyncPolicies
					// Not really needed but keeps the data structure clean
					namespace.SyncPolicies = map[string]string{}
					namespace.SyncFilters = map[string]policy.PolicyFilter{}
				}
			}

			// Add Sync Target information to the Namespace Configuration
			// If the PopulatedNamespaceConfiguration doesn't exist create it
			//    by querying the Core Service for the namespace's information
			for _, syncConfig := range toCreate {
				logrus.Debugf("Adding: %+v", syncConfig)
				coreService := pcs.populatedConfigs[syncConfig.SourceCoreService]

				namespace, ok := coreService.Namespaces[syncConfig.SourceNamespace]
				syncKey := config.SyncKey{
					CoreService: syncConfig.TargetCoreService,
					Namespace:   syncConfig.TargetNamespace,
				}
				if !ok {
					namespace = newNamespace(tokens, coreService, syncConfig.SourceNamespace)
					coreService.Namespaces[namespace.Name] = namespace
				}

				// Updates the Policies and Filters if they changed
				namespace.SyncPolicies = syncConfig.SourcePolicies
				for k, v := range namespace.SyncPolicies {
					f, err := policy.Parse(v)
					if err != nil {
						logrus.Errorf("problem parsing policy %s / %s / %s, failing open: %v", coreService.Endpoint, namespace.Name, k, err)
						namespace.SyncFilters[k] = policy.DefaultOpenPolicyFilter
					} else {
						namespace.SyncFilters[k] = f
					}
				}

				if syncTarget, ok := namespace.SyncTargets[syncKey]; ok {
					// Updating Sync Type field, no need to update the target
					syncTarget.SyncType = syncConfig.SyncType
				} else {
					// Adding a new Sync Target
					namespace.SyncTargets[syncKey] = &config.SyncTarget{
						SyncType: syncConfig.SyncType,
						Target:   nil, // There will be a second pass to set this link
					}
				}
			}

			// Link the SyncTargets together
			for _, coreService := range pcs.populatedConfigs {
				for _, namespace := range coreService.Namespaces {
					for key, val := range namespace.SyncTargets {
						if val.Target == nil {
							val.Target = pcs.populatedConfigs[key.CoreService].Namespaces[key.Namespace]
							if val.Target == nil { // The target isn't the source of sync configurations
								targetNamespace := newNamespace(tokens, pcs.populatedConfigs[key.CoreService], key.Namespace)
								pcs.populatedConfigs[key.CoreService].Namespaces[key.Namespace] = targetNamespace
								val.Target = targetNamespace
							}
						}
					}
				}
			}

			// Release the write locks as we have finished mutating the data structure
			for _, coreService := range toLock {
				coreService.L.Unlock()
			}

			// Notify ForceReload() that the update has finished
			pcs.reloadFinished.Broadcast()

			logrus.Info("Finished with sync configuration update")
		case <-ctx.Done():
			logrus.Info("Update Muxer is stopping")
			return
		}
	}
}
