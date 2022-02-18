package database

import (
	"testing"
	"time"

	"github.com/gigantum/hoss-core/pkg/test"
	"github.com/gigantum/hoss-service/policy"
)

func TestCreateDatasetExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateDataset(ns, "test_dataset", "my dataset", "/test_dataset", "test_user")
	if err == nil {
		t.Fatal("Expected error but create dataset succeeded")
	}

	if err != ErrExists {
		t.Fatalf("Expected already exists error but got a different type: %v", err)
	}
}

func TestCreateDatasetNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateDataset(ns, "test_dataset1", "my dataset", "/test_dataset1", "test_user")
	if err != nil {
		t.Fatalf("Expected no error but create dataset failed: %v", err)
	}

	ds, err := db.GetDataset(ns, "test_dataset1")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}

	test.AssertEqual(t, ds.Name, "test_dataset1")
	test.AssertEqual(t, ds.Description, "my dataset")
	test.AssertEqual(t, ds.RootDirectory, "/test_dataset1")
}

func TestCreateDatasetNotExistingOwner(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateDataset(ns, "test_dataset1", "test dataset description", "/test_dataset1", "test_user1")
	if err != nil {
		t.Fatal("Expected no error but create dataset failed: ", err.Error())
	}
}

func TestGetDatasetExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.Namespace.Name, "test_namespace")
	test.AssertEqual(t, ds.Name, "test_dataset")
	test.AssertEqual(t, ds.DeleteStatus, "NOT_SCHEDULED")
}

func TestGetDatasetNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	_, err = db.GetDataset(ns, "test_dataset1")
	if err == nil {
		t.Fatal("Expected error but get dataset succeeded")
	}

	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestSetDatasetSyncNotChanged(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}

	ts := time.Now()
	lastUpdate, err := db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ds.SyncEnabled, false)
	test.AssertEqual(t, ts.After(lastUpdate), true)

	err = db.SetDatasetSync(ds, false, SYNC_TYPE_SIMPLEX, "")
	if err != nil {
		t.Fatal("Expected no error but set dataset sync failed: ", err.Error())
	}

	lastUpdate, err = db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	// If the SyncEnabled flag didn't change then the database should not be updated
	test.AssertEqual(t, ts.After(lastUpdate), true)
}

func TestSetDatasetSyncChanged(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}

	ts := time.Now()
	lastUpdate, err := db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ds.SyncEnabled, false)
	test.AssertEqual(t, ts.After(lastUpdate), true)

	err = db.SetDatasetSync(ds, true, SYNC_TYPE_SIMPLEX, policy.DefaultOpenPolicy)
	if err != nil {
		t.Fatal("Expected no error but set dataset sync failed: ", err.Error())
	}

	// Ensure the model object was updated without re-query
	test.AssertEqual(t, ds.SyncEnabled, true)

	// Make sure the trigger was run and the last update timestamp is updated
	time.Sleep(1 * time.Second)
	lastUpdate, err = db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ts.Before(lastUpdate), true)

	// Make sure the database was updated
	ds, err = db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}

	test.AssertEqual(t, ds.SyncEnabled, true)
}

func TestDeleteDatasetExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Failed to get dataset")
	}
	test.AssertEqual(t, ds.Namespace.Name, "test_namespace")
	test.AssertEqual(t, ds.Name, "test_dataset")

	userPerms, err := db.GetPermissionsByUser(&ns.ObjectStore, "test_user", false)
	if err != nil {
		t.Fatal("Failed to get user permissions")
	}
	test.AssertEqual(t, len(userPerms), 1)

	err = db.DeleteDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but delete dataset failed")
	}

	_, err = db.GetDataset(ns, "test_dataset")
	if err == nil {
		t.Fatal("Expected error getting deleted dataset")
	}

	userPerms, err = db.GetPermissionsByUser(&ns.ObjectStore, "test_user", false)
	if err != nil {
		t.Fatal("Failed to get user permissions")
	}
	test.AssertEqual(t, len(userPerms), 0)
}

func TestDeleteDatasetNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.DeleteDataset(ns, "test_dataset1")
	if err != nil {
		t.Fatal("Expected no error but delete dataset failed")
	}
}

func TestUpdateDatasetPermissionsExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.UpdateDatasetPermissions(ns, "test_dataset", "test_group", "rw")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed")
	}
}

func TestUpdateDatasetPermissionsNotExistingDataset(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.UpdateDatasetPermissions(ns, "test_dataset1", "test_group", "rw")
	if err == nil {
		t.Fatal("Expected error but update dataset permissions succeeded")
	}

	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestUpdateDatasetPermissionsNotExistingUser(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.UpdateDatasetPermissions(ns, "test_dataset", "test_group1", "rw")
	if err != nil {
		t.Fatalf("Expected no error but update dataset permissions failed: %v", err)
	}
}

func TestUpdateDatasetPermissionsInvalidPermission(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.UpdateDatasetPermissions(ns, "test_dataset", "test_group", "unknown")
	if err == nil {
		t.Fatal("Expected error but update dataset permissions succeeded")
	}

	if err != ErrInvalidInput {
		t.Fatal(err)
	}
}

func TestRemoveDatasetPermissionsExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.RemoveDatasetPermissions(ns, "test_dataset", "test_group")
	if err != nil {
		t.Fatalf("Expected no error but remove dataset permissions failed: %v", err)
	}
}

func TestRemoveDatasetPermissionsNotExistingUser(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.RemoveDatasetPermissions(ns, "test_dataset", "test_group1")
	if err != nil {
		t.Fatalf("Expected no error but remove dataset permissions failed: %v\n", err)
	}
}

func TestRemoveDatasetPermissionsNotExistingDataset(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.RemoveDatasetPermissions(ns, "test_dataset1", "test_user")
	if err == nil {
		t.Fatal("Expected error but remove dataset permissions succeeded")
	}

	if err != ErrNotFound {
		t.Fatal(err)
	}
}

func TestSetDatasetDeleteMarker(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	// Set to scheduled
	err = db.SetDatasetDeleteMarker(ns, "test_dataset", SCHEDULED, 60)
	if err != nil {
		t.Fatalf("Error while setting delete marker: %s", err.Error())
	}
	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "SCHEDULED")
	delete_on_after := time.Now().Add(time.Duration(61) * time.Minute)
	delete_on_before := time.Now().Add(time.Duration(59) * time.Minute)

	if ds.DeleteOn.Before(delete_on_before) {
		t.Fatal("Delete On is earlier than expected")
	}
	if ds.DeleteOn.After(delete_on_after) {
		t.Fatal("Delete On is later than expected")
	}

	// Set to in progress
	err = db.SetDatasetDeleteMarker(ns, "test_dataset", IN_PROGRESS, 60)
	if err != nil {
		t.Fatalf("Error while setting delete marker: %s", err.Error())
	}
	ds, err = db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "IN_PROGRESS")

	// Set to not scheduled
	err = db.SetDatasetDeleteMarker(ns, "test_dataset", NOT_SCHEDULED, 60)
	if err != nil {
		t.Fatalf("Error while setting delete marker: %s", err.Error())
	}
	ds, err = db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "NOT_SCHEDULED")

	// Set to error
	err = db.SetDatasetDeleteMarker(ns, "test_dataset", ERROR, 60)
	if err != nil {
		t.Fatalf("Error while setting delete marker: %s", err.Error())
	}
	ds, err = db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "ERROR")

}

func TestGetDatasetsToDelete(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	datasets, err := db.GetDatasetsByDeleteStatus(SCHEDULED)
	if err != nil {
		t.Fatalf("Error while getting datasets to delete: %s", err.Error())
	}
	test.AssertEqual(t, len(datasets), 0)

	// Set to scheduled
	err = db.SetDatasetDeleteMarker(ns, "test_dataset", SCHEDULED, 60)
	if err != nil {
		t.Fatalf("Error while setting delete marker: %s", err.Error())
	}
	ds, err := db.GetDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get dataset failed: ", err.Error())
	}
	test.AssertEqual(t, ds.DeleteStatus, "SCHEDULED")
	delete_on_after := time.Now().Add(time.Duration(61) * time.Minute)
	delete_on_before := time.Now().Add(time.Duration(59) * time.Minute)

	if ds.DeleteOn.Before(delete_on_before) {
		t.Fatal("Delete On is earlier than expected")
	}
	if ds.DeleteOn.After(delete_on_after) {
		t.Fatal("Delete On is later than expected")
	}

	datasets, err = db.GetDatasetsByDeleteStatus(SCHEDULED)
	if err != nil {
		t.Fatalf("Error while getting datasets to delete: %s", err.Error())
	}
	test.AssertEqual(t, len(datasets), 1)
	test.AssertEqual(t, datasets[0].Name, "test_dataset")
	test.AssertEqual(t, datasets[0].DeleteStatus, "SCHEDULED")
}
