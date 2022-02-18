package store

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/test"
)

const testUsername = "testuser"

func TestVerifyMinioConnection(t *testing.T) {
	_, _, _, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed to run minio fixture: %v", err)
	}

	cmd := exec.Command("mc", "alias", "ls", "hoss-default")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to load mc alias: %v", err)
	}

	cmd = exec.Command("mc", "admin", "info", "hoss-default")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to connect to minio: %v", err)
	}
}

func TestLoad(t *testing.T) {
	c, _, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed to run minio fixture: %v", err)
	}

	// Delete hoss-default from mc config
	cmd := exec.Command("mc", "alias", "rm", "hoss-default")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to remove mc alias: %v", err)
	}

	cmd = exec.Command("mc", "alias", "ls", "hoss-default")
	if err := cmd.Run(); err == nil {
		t.Fatal("expected default alias to not exist but it does")
	}

	obj, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatalf("failed to get object store: %v", err)
	}

	newStore := &MinioStore{}
	newStore.Load(c, obj)

	cmd = exec.Command("mc", "alias", "ls", "hoss-default")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to load mc alias: %v", err)
	}

	cmd = exec.Command("mc", "admin", "info", "hoss-default")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to connect to minio: %v", err)
	}

}

func TestCreateDatasetExisting(t *testing.T) {
	_, currentStore, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	err = currentStore.CreateDataset("test-ds-1", ns)
	if err == nil {
		t.Fatal("Expected error but create dataset succeeded")
	}

	if err.Error() != "Failed to create dataset. Dataset already exists" {
		t.Fatal(err)
	}
}

func TestCreateDataset(t *testing.T) {
	coreConfig, currentStore, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	err = currentStore.CreateDataset("test-ds-3", ns)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	m := NewMetadataFile("test-ds-3")

	mPath := filepath.Join(test.DefaultBucketDir(t, coreConfig.Namespaces[0].Bucket), m.Key())
	_, err = os.Stat(mPath)
	if os.IsNotExist(err) {
		t.Fatalf("Dataset metadata file doesn't exist after create: %v", err)
	}
}

func TestDeleteDataset(t *testing.T) {
	coreConfig, currentStore, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	m := NewMetadataFile("test-ds-1")
	mPath := filepath.Join(test.DefaultBucketDir(t, coreConfig.Namespaces[0].Bucket), m.Key())
	dPath := filepath.Join(test.DefaultBucketDir(t, coreConfig.Namespaces[0].Bucket), "test-ds-1")
	_, err = os.Stat(mPath)
	if os.IsNotExist(err) {
		t.Fatalf("Dataset metadata file expected to exist at test start: %v", err)
	}
	_, err = os.Stat(dPath)
	if os.IsNotExist(err) {
		t.Fatalf("Dataset directory expected to exist at test start: %v", err)
	}

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	err = currentStore.DeleteDataset("test-ds-1", ns)
	if err != nil {
		t.Fatalf("Failed to delete dataset: %v", err)
	}

	_, err = os.Stat(mPath)
	if !os.IsNotExist(err) {
		t.Fatalf("Dataset metadata file still exists after create: %v", err)
	}
	_, err = os.Stat(dPath)
	if !os.IsNotExist(err) {
		t.Fatalf("Dataset directory file still exists after create: %v", err)
	}
}

func TestUpdatePolicy(t *testing.T) {
	coreConfig, currentStore, _, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy exists but expect it not to yet")
	}

	permData, err := mockPermissionQuery(t, coreConfig)
	if err != nil {
		t.Fatalf("Failed to mock permissions: %v", err)
	}

	err = currentStore.SetUserPolicy(testUsername, permData)
	if err != nil {
		t.Fatalf("Failed to update policy: %v", err)
	}

	if !test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy doesn't exist but expect it to")
	}
}

func TestDeletePolicy(t *testing.T) {
	coreConfig, currentStore, _, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy exists but expect it not to yet")
	}

	permData, err := mockPermissionQuery(t, coreConfig)
	if err != nil {
		t.Fatalf("Failed to mock permissions: %v", err)
	}

	err = currentStore.SetUserPolicy(testUsername, permData)
	if err != nil {
		t.Fatalf("Failed to update policy: %v", err)
	}

	if !test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy doesn't exist but expect it to")
	}

	err = currentStore.DeleteUserPolicy(testUsername)
	if err != nil {
		t.Fatalf("Failed to delete policy: %v", err)
	}

	if test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy exists but expect it not to yet")
	}
}

func TestUpdateDenyPolicy(t *testing.T) {
	_, currentStore, _, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy exists but expect it not to yet")
	}

	var permData []*database.Permission

	err = currentStore.SetUserPolicy(testUsername, permData)
	if err != nil {
		t.Fatalf("Failed to update policy: %v", err)
	}

	if !test.MinioPolicyExists(t, currentStore.UserPolicyName(testUsername)) {
		t.Fatalf("policy doesn't exist but expect it to")
	}

}

func TestPatchEventsNotEnabled(t *testing.T) {
	_, currentStore, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	err = db.CreateDataset(ns, "test-ds-3", "my dataset", "/test-ds-3", "test_user")
	if err != nil {
		t.Fatalf("Expected no error but create dataset failed: %v", err)
	}
	ds, err := db.GetDataset(ns, "test-ds-3")
	if err != nil {
		t.Fatalf("failed to load dataset: %v", err)
	}
	err = currentStore.CreateDataset("test-ds-3", ns)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	// If events are enabled, disable them
	isEnabled, err := currentStore.EventsEnabled(ns, ds)
	if err != nil {
		t.Fatalf("failed to check if dataset %s sync is enabled: %s", ds.Name, err.Error())
	}
	if isEnabled {
		err = currentStore.DisableEvents(ns, ds)
		if err != nil {
			t.Fatalf("failed to enable events before patch: %v", err)
		}
	}

	objMap := map[string]ObjectStore{}
	objMap[currentStore.GetName()] = currentStore
	err = PatchMinioEvents(db, objMap)
	if err != nil {
		t.Fatalf("Failed to patch events: %v", err)
	}

	isEnabled, err = currentStore.EventsEnabled(ns, ds)
	if err != nil {
		t.Fatalf("failed to check if events were enabled: %v", err)
	}
	if !isEnabled {
		t.Fatalf("Expected events to be enabled after patch.")
	}
}

func TestPatchEventsEnabled(t *testing.T) {
	_, currentStore, db, err := SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}

	err = db.CreateDataset(ns, "test-ds-3", "my dataset", "/test-ds-3", "test_user")
	if err != nil {
		t.Fatalf("Expected no error but create dataset failed: %v", err)
	}

	err = currentStore.CreateDataset("test-ds-3", ns)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}
	ds, err := db.GetDataset(ns, "test-ds-3")
	if err != nil {
		t.Fatalf("failed to load dataset: %v", err)
	}

	// If events are disabled, enable them
	isEnabled, err := currentStore.EventsEnabled(ns, ds)
	if err != nil {
		t.Fatalf("failed to check if dataset %s sync is enabled: %s", ds.Name, err.Error())
	}
	if !isEnabled {
		err = currentStore.EnableEvents(ns, ds)
		if err != nil {
			t.Fatalf("failed to enable events before patch: %v", err)
		}
	}

	objMap := map[string]ObjectStore{}
	objMap[currentStore.GetName()] = currentStore
	err = PatchMinioEvents(db, objMap)
	if err != nil {
		t.Fatalf("Failed to patch events: %v", err)
	}

	isEnabled, err = currentStore.EventsEnabled(ns, ds)
	if err != nil {
		t.Fatalf("failed to check if events were enabled: %v", err)
	}
	if !isEnabled {
		t.Fatalf("Expected events to be enabled after patch.")
	}
}
