package worker

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gigantum/hoss-core/pkg/test"
)

func TestDeleteDataset(t *testing.T) {
	config, currentStore, db, err := SetupWorkerTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// The default delete delay is 0, which is fine for testing,
	// but we need to change the period to 1 second so tests happen quickly
	config.Server.DatasetDeletePeriodSeconds = 1

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	ds, err := db.GetDataset(ns, "delete-test-1")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "NOT_SCHEDULED")

	// To test the delete worker we:
	// 1) Check your permissions and show that the dataset is there
	// 2) Mark the dataset "delete-test-1" as "scheduled" 1 minute into the future
	// 3) Check your permissions and show that the dataset is not there
	// 4) Check your permissions again, including deleted datasets, and show that the dataset is there
	// 5) Start the dataset delete goroutine
	// 6) Wait 55 seconds, then make sure the dataset exists and is still "scheduled"
	// 7) Wait 10 seconds for delete to happen
	// 8) Make sure the dataset is completely gone now
	perms, err := db.GetPermissionsByUser(&ns.ObjectStore, "testuser", false)
	if err != nil {
		t.Fatalf("Expected no error but get datasets by user failed: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("Expected access to 1 dataset but got %v", len(perms))
	}
	test.AssertEqual(t, perms[0].Dataset.Name, "delete-test-1")

	err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.SCHEDULED, 1)
	if err != nil {
		t.Fatalf("Failed to set delete_status to SCHEDULED for %s: %s", ds, err.Error())
	}

	perms, err = db.GetPermissionsByUser(&ns.ObjectStore, "testuser", false)
	if err != nil {
		t.Fatalf("Expected no error but get dataset perms by user failed: %v", err)
	}
	if len(perms) != 0 {
		t.Fatalf("Expected access to 0 datasets but got %v", len(perms))
	}

	perms, err = db.GetPermissionsByUser(&ns.ObjectStore, "testuser", true)
	if err != nil {
		t.Fatalf("Expected no error but get dataset perms by user failed: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("Expected access to 1 dataset but got %v", len(perms))
	}
	test.AssertEqual(t, perms[0].Dataset.Name, "delete-test-1")

	objMap := map[string]store.ObjectStore{}
	objMap[currentStore.GetName()] = currentStore
	exitCh := make(chan bool)
	go DeleteDatasetWorker(config, db, objMap, exitCh)

	time.Sleep(55 * time.Second)

	ds, err = db.GetDataset(ns, "delete-test-1")
	if err != nil {
		t.Fatalf("Expected no error but get dataset failed: %v", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "SCHEDULED")

	time.Sleep(10 * time.Second)

	_, err = db.GetDataset(ns, "delete-test-1")
	if err == nil {
		t.Fatalf("Expected an error when fetching the dataset: %v", err.Error())
	}

	if _, err := os.Stat(filepath.Join(test.DefaultBucketDir(t, config.Namespaces[0].Bucket), "dataset-test-1")); !os.IsNotExist(err) {
		t.Fatal("Expected dataset contents to be removed but they still exist.")
	}

	// Shutdown the worker before teardown
	exitCh <- true
	time.Sleep(5 * time.Second)
}

// TestDeleteDatasetInErrorState tests to make sure that when the dataset
// delete worker starts up, it does a check to "reset" any errored
// datasets. This gives admins an easy path to re-run the delete process
// after they have addressed errors.
func TestDeleteDatasetInErrorState(t *testing.T) {
	config, currentStore, db, err := SetupWorkerTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// The default delete delay is 0, which is fine for testing,
	// but we need to change the period to 1 second so tests happen quickly
	config.Server.DatasetDeletePeriodSeconds = 1

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	ds, err := db.GetDataset(ns, "delete-test-1")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "NOT_SCHEDULED")

	// To test the delete worker we:
	// 1) Check your permissions and show that the dataset is there
	// 2) Mark the dataset "delete-test-1" as "ERROR" 1 minute into the future
	// 3) Check your permissions and show that the dataset is not there
	// 4) Start the dataset delete goroutine
	// 5) Wait 5 seconds for delete to happen, when moved from ERROR it schedules immediately
	// 6) Make sure the dataset is completely gone now
	perms, err := db.GetPermissionsByUser(&ns.ObjectStore, "testuser", false)
	if err != nil {
		t.Fatalf("Expected no error but get datasets by user failed: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("Expected access to 1 dataset but got %v", len(perms))
	}
	test.AssertEqual(t, perms[0].Dataset.Name, "delete-test-1")

	err = db.SetDatasetDeleteMarker(ds.Namespace, ds.Name, database.ERROR, 1)
	if err != nil {
		t.Fatalf("Failed to set delete_status to ERROR for %s: %s", ds, err.Error())
	}

	perms, err = db.GetPermissionsByUser(&ns.ObjectStore, "testuser", false)
	if err != nil {
		t.Fatalf("Expected no error but get dataset perms by user failed: %v", err)
	}
	if len(perms) != 0 {
		t.Fatalf("Expected access to 0 datasets but got %v", len(perms))
	}

	objMap := map[string]store.ObjectStore{}
	objMap[currentStore.GetName()] = currentStore
	exitCh := make(chan bool)
	go DeleteDatasetWorker(config, db, objMap, exitCh)

	time.Sleep(5 * time.Second)

	_, err = db.GetDataset(ns, "delete-test-1")
	if err == nil {
		t.Fatalf("Expected an error when fetching the dataset: %v", err.Error())
	}

	if _, err := os.Stat(filepath.Join(test.DefaultBucketDir(t, config.Namespaces[0].Bucket), "dataset-test-1")); !os.IsNotExist(err) {
		t.Fatal("Expected dataset contents to be removed but they still exist.")
	}

	// Shutdown the worker before teardown
	exitCh <- true
	time.Sleep(5 * time.Second)
}
