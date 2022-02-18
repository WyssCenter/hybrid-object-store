package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gigantum/hoss-core/pkg/test"
	"github.com/pkg/errors"
)

// SetupWorkerTest sets the client and config for testing
func SetupWorkerTest(t *testing.T) (*config.Configuration, store.ObjectStore, *database.Database, error) {
	testConfig, currentStore, db, err := store.SetupMinioTest(t)
	if err != nil {
		t.Fatalf("failed to run setup minio test fixture: %v", err)
	}

	// Set cleanup
	t.Cleanup(func() {
		TeardownWorkerTest(t, testConfig, db)
		store.TeardownMinioTest(t, testConfig)
		database.TeardownDatabaseTest(t, db)
	})

	// Fully create a dataset as if it's being created from the API and populate it manually
	// with some data
	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace: %v", err)
	}
	err = db.CreateDataset(ns, "delete-test-1", "my test dataset", "delete-test-1/", "testuser")
	if err != nil {
		t.Fatalf("failed to create test dataset in db: %v", err)
	}
	err = db.UpdateDatasetPermissions(ns, "delete-test-1", db.GetUserDefaultGroup("testuser"), database.PERM_READ_WRITE)
	if err != nil {
		t.Fatalf("failed to create test dataset permissions db: %v", err)
	}
	err = currentStore.CreateDataset("delete-test-1", ns)
	if err != nil {
		t.Fatalf("failed to create test dataset in object store: %v", err)
	}
	ds, err := db.GetDataset(ns, "delete-test-1")
	if err != nil {
		t.Fatalf("failed to load test dataset: %v", err)
	}
	err = currentStore.EnableEvents(ns, ds)
	if err != nil {
		t.Fatalf("failed to enable test dataset events: %v", err)
	}
	perms, err := db.GetPermissionsByUser(&ns.ObjectStore, "testuser", false)
	if err != nil {
		t.Fatalf("failed to get user perms: %v", err)
	}
	err = currentStore.SetUserPolicy("testuser", perms)
	if err != nil {
		t.Fatalf("failed to set user policy: %v", err)
	}

	// Create dummy files
	for i := 0; i < 200; i++ {
		filename := filepath.Join(test.DefaultBucketDir(t, testConfig.Namespaces[0].Bucket),
			"delete-test-1", fmt.Sprintf("%v.txt", i))
		err := ioutil.WriteFile(filename, []byte("DUMMY FILE"), 0777)
		if err != nil {
			errors.Wrap(err, "error creating test data in test dataset")
		}
	}

	return testConfig, currentStore, db, nil
}

// TeardownMinioTest gracefully tries to remove all data created by a test
func TeardownWorkerTest(t *testing.T, c *config.Configuration, db *database.Database) {
	possibleDatasets := [...]string{"delete-test-1"}
	for _, name := range possibleDatasets {
		p := filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket), name)
		os.RemoveAll(p)
	}

	// Try to delete the dataset if it exists
	ns, err := db.GetNamespace("default")
	if err != nil {
		t.Fatalf("failed to load namespace during cleanup: %v", err)
	}
	_, err = db.GetDataset(ns, "delete-test-1")
	if err == nil {
		db.DeleteDataset(ns, "delete-test-1")
	}
}
