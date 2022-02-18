package worker

import (
	"time"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/sirupsen/logrus"
)

func DeleteDatasetWorker(c *config.Configuration, db *database.Database, stores map[string]store.ObjectStore, exit <-chan bool) {
	// On first boot, if any datasets are in an ERROR state we reset them to SCHEDULED. This gives an
	// easy path for admins to attempt to fix a failed delete and then trigger the delete again.
	// If the same error occurs it will just go back into ERROR state and continue to be skipped until
	// another service restart.
	logrus.Info("[DATASET DELETE WORKER] STARTING WORKER")
	datasets, err := db.GetDatasetsByDeleteStatus(database.ERROR)
	if err != nil {
		logrus.Errorf("[DATASET DELETE WORKER] Failed to list datasets scheduled for delete on start up: %s", err.Error())
	}
	for _, ds := range datasets {
		if ds.DeleteStatus == string(database.ERROR) {
			// Set the datasets as SCHEDULED again with the delay set to 0. This will
			// immediately re-try the delete.
			err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.SCHEDULED, 0)
			if err != nil {
				logrus.Errorf("[DATASET DELETE WORKER] Failed to set delete_status to SCHEDULED for %s: %s", ds, err.Error())
				continue
			}
			logrus.Infof("[DATASET DELETE WORKER] Reset delete_status from ERROR to SCHEDULED for: %s", ds)
		}
	}

	for {
		select {
		case <-exit:
			logrus.Info("[DATASET DELETE WORKER] Shutting down.")
			return
		default:
			datasets, err := db.GetDatasetsByDeleteStatus(database.SCHEDULED)
			if err != nil {
				logrus.Errorf("[DATASET DELETE WORKER] Failed to list datasets scheduled for delete: %s", err.Error())
				time.Sleep(time.Duration(c.Server.DatasetDeletePeriodSeconds * int(time.Second)))
				continue
			}

			for _, ds := range datasets {
				if ds.DeleteOn.Before(time.Now().UTC()) {
					if ds.DeleteStatus == string(database.ERROR) {
						logrus.Infof("[DATASET DELETE WORKER] Skipping dataset %s in ERROR state. Please review logs.", ds)
						continue
					}

					logrus.Infof("[DATASET DELETE WORKER] Starting delete process for %s", ds)
					err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.IN_PROGRESS, c.Server.DatasetDeleteDelayMinutes)
					if err != nil {
						logrus.Errorf("[DATASET DELETE WORKER] Failed to set delete_status to IN_PROGRESS for %s: %s", ds, err.Error())
						continue
					}

					// Disable sync before deleting the dataset in the object store
					// When the object store delete op runs, it will emit delete object events that
					// the sync service will see. This is important because the search index will then
					// update and remove the data related to these objects. We don't, however, want to
					// remove the target data. Future work will provide a proper way to configure "cascading"
					// the delete operation to a target dataset via a single API notification
					if ds.SyncEnabled {
						// Note: This will mutate the sync_configuration_meta.last_updated cell with the current timestamp, if the query succeeds
						// This will trigger the configuration reload delay inside the sync service.
						logrus.Infof("[DATASET DELETE WORKER] Disabling sync for dataset %s before delete. Waiting 25 seconds before continuing.", ds)
						err = db.SetDatasetSync(ds, false, "", "")
						if err != nil {
							// This should be a reliable operation. If error occurs mark error state.
							logrus.Errorf("[DATASET DELETE WORKER] Error when attempting to disable sync for %s: %s", ds, err.Error())
							err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.ERROR, c.Server.DatasetDeleteDelayMinutes)
							if err != nil {
								logrus.Errorf("[DATASET DELETE WORKER] Failed to set delete_status to ERROR for %s: %s", ds, err.Error())
							}
							continue
						}
						time.Sleep(25 * time.Second)
					}

					// Delete data in the object store.
					// This can take a while depending on the number of objects and in theory could have transitory failures
					// so the object delete operation is done with up to 3 retries
					store, ok := stores[ds.Namespace.ObjectStore.Name]
					if !ok {
						// This should be a reliable operation. If error occurs mark error state.
						logrus.Errorf("[DATASET DELETE WORKER] Object store not found for dataset %s. Try restarting the server.", ds)
						err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.ERROR, c.Server.DatasetDeleteDelayMinutes)
						if err != nil {
							logrus.Errorf("[DATASET DELETE WORKER] Failed to set delete_status to ERROR for %s: %s", ds, err.Error())
						}
						continue
					}

					try := 0
					succeeded := false
					for try < 3 {
						err = store.DeleteDataset(ds.Name, ds.Namespace)
						if err != nil {
							logrus.Warnf("[DATASET DELETE WORKER] (Try %v of 3) Failed to delete dataset %s data: %s", try, ds.Name, err.Error())
							try = try + 1
							continue
						}
						succeeded = true
						break
					}

					if !succeeded {
						logrus.Errorf("[DATASET DELETE WORKER] Failed to delete dataset %s data: %s", ds.Name, err.Error())
						err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.ERROR, c.Server.DatasetDeleteDelayMinutes)
						if err != nil {
							logrus.Errorf("[DATASET DELETE WORKER] Failed to set delete_status to ERROR for %s: %s", ds, err.Error())
						}
						continue
					}

					// Disable bucket events for dataset prefix since the dataset prefix should no longer exist
					isEnabled, err := store.EventsEnabled(ds.Namespace, ds)
					if err != nil {
						logrus.Errorf("[DATASET DELETE WORKER] Couldn't check if dataset %s events is enabled, continuing with delete operation: %s", ds.Name, err.Error())
						isEnabled = false
					}
					if isEnabled {
						err = store.DisableEvents(ds.Namespace, ds)
						if err != nil {
							logrus.Errorf("[DATASET DELETE WORKER] Failed to disable dataset %s sync, continuing with delete operation: %s", ds.Name, err.Error())
						}
					}

					// Finally, delete the database entry
					err = db.DeleteDataset(ds.Namespace, ds.Name)
					if err != nil {
						logrus.Errorf("[DATASET DELETE WORKER] Failed to delete dataset %s database entry: %s", ds.Name, err.Error())
					}

					logrus.Infof("[DATASET DELETE WORKER] Successfully deleted %s", ds)
				}
			}

			// Wait the specified delay before trying to delete again
			time.Sleep(time.Duration(c.Server.DatasetDeletePeriodSeconds * int(time.Second)))
		}
	}
}
