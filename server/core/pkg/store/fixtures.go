package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/test"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// SetupMinioTest sets the client and config for testing
func SetupMinioTest(t *testing.T) (*config.Configuration, ObjectStore, *database.Database, error) {
	defaultobj := &config.ObjectStore{
		Name:        "default",
		Description: "Default object store",
		Type:        "minio",
		Endpoint:    "http://localhost",
	}
	defaultns := &config.Namespace{
		Name:        "default",
		Description: "Default namespace",
		Bucket:      "data",
		ObjectStore: "default",
	}
	objectStores := []config.ObjectStore{}
	objectStores = append(objectStores, *defaultobj)

	nss := []config.Namespace{}
	nss = append(nss, *defaultns)

	server := &config.Server{Dev: true}
	testConfig := &config.Configuration{ObjectStores: objectStores, Namespaces: nss, Server: *server}

	// Set env vars
	err := test.LoadEnvFile("~/.hoss/.env")
	if err != nil {
		return nil, nil, nil, err
	}

	if err := os.Setenv("POSTGRES_DB", "hoss_core"); err != nil {
		return nil, nil, nil, err
	}
	if err := os.Setenv("POSTGRES_HOST", "localhost:5432"); err != nil {
		return nil, nil, nil, err
	}

	// Setting the TESTING env var will prevent minio endpoints from being converted
	// from localhost to minio:9000
	if err := os.Setenv("TESTING", "1"); err != nil {
		return nil, nil, nil, err
	}

	db := database.Load()

	// Set cleanup
	t.Cleanup(func() {
		TeardownMinioTest(t, testConfig)
		database.TeardownDatabaseTest(t, db)
	})

	database.BootstrapDefaults(testConfig, db)
	availableStores, err := db.ListObjectStores(25, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	stores := LoadObjectStores(testConfig, availableStores)

	// Create metadata files for test datasets
	mockDataset(t, testConfig, "test-ds-1")
	mockDataset(t, testConfig, "test-ds-2")

	return testConfig, stores[availableStores[0].Name], db, nil
}

// TeardownMinioTest gracefully tries to remove all data created by a test
func TeardownMinioTest(t *testing.T, c *config.Configuration) {
	possibleDatasets := [...]string{"test-ds-1", "test-ds-2", "test-ds-3"}
	for _, name := range possibleDatasets {
		p := filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket), name)
		os.RemoveAll(p)
	}

	// Remove the testuser policy if it had been created
	cmd := exec.Command("mc", "admin", "policy", "remove", "hoss-default", "testuser")
	cmd.Run()

	// Remove bucket events if they have been created
	cmd = exec.Command("mc", "event", "remove", "hoss-default/data", "arn:minio:sqs::_:amqp")
	cmd.Run()

	// List the aliases
	cmd = exec.Command("mc", "alias", "list")
	out, err := cmd.CombinedOutput()
	t.Logf("output from list mc alias: %v", out)
	if err != nil {
		t.Fatalf("failed to list mc alias: %v", err)
	}

	// Remove the default alias
	cmd = exec.Command("mc", "alias", "remove", "hoss-default")
	out, err = cmd.CombinedOutput()
	t.Logf("output from remove mc alias: %v", out)
	if err != nil {
		t.Fatalf("failed to remove mc alias: %v", err)
	}
}

//writeTestDatasetMetadata
func mockDataset(t *testing.T, c *config.Configuration, name string) error {
	os.Mkdir(filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket), name), os.FileMode(0755))

	metadata := NewMetadataFile(name)
	raw, err := yaml.Marshal(metadata)
	if err != nil {
		return errors.Wrapf(err, "error creating test dataset %s", name)
	}

	if err := ioutil.WriteFile(filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket),
		metadata.Key()), raw, 0777); err != nil {
		return errors.Wrapf(err, "error creating test dataset %s", name)
	}

	// Create dummy files
	for i := 0; i < 50; i++ {
		filename := filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket),
			name, fmt.Sprintf("%v.txt", i))
		err := ioutil.WriteFile(filename, []byte("DUMMY FILE"), 0777)
		if err != nil {
			errors.Wrapf(err, "error creating test data in dataset %s", name)
		}
	}

	os.Mkdir(filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket),
		name, "subdir"), os.FileMode(0522))

	for i := 0; i < 50; i++ {
		filename := filepath.Join(test.DefaultBucketDir(t, c.Namespaces[0].Bucket), name, "subdir", fmt.Sprintf("%v.txt", i))
		err := ioutil.WriteFile(filename, []byte("DUMMY FILE"), 0755)
		if err != nil {
			errors.Wrapf(err, "error creating test data in dataset %s", name)
		}
	}

	return nil
}

func mockPermissionQuery(t *testing.T, c *config.Configuration) ([]*database.Permission, error) {
	db := database.Load()

	err := db.CreateNamespace("test_namespace", "My test namespace", c.ObjectStores[0].Name, "data")
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't create namespace record")
	}
	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't get namespace record")
	}

	user := &database.User{Id: 1, Username: "testuser"}
	group := &database.Group{Id: 1, GroupName: "testgroup"}
	dataset1 := &database.Dataset{Id: 1, Namespace: ns, Name: "test-ds-1", Owner: user}
	dataset2 := &database.Dataset{Id: 1, Namespace: ns, Name: "test-ds-2", Owner: user}

	// Mock database data
	var permissions []*database.Permission

	permissions = append(permissions, &database.Permission{Group: group, Dataset: dataset1, Permission: database.PERM_READ})
	permissions = append(permissions, &database.Permission{Group: group, Dataset: dataset2, Permission: database.PERM_READ_WRITE})

	return permissions, nil
}
