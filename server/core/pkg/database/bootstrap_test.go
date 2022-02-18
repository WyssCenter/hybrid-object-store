package database

import (
	"testing"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/test"
)

func TestBootstrapDoesNotExist(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	defaultobj := &config.ObjectStore{
		Name:        "default",
		Description: "Default object store",
		Type:        "minio",
		Endpoint:    "http://localhost",
	}
	defaultns := &config.Namespace{
		Name:        "default",
		Description: "Default namespace",
		Bucket:      "test-data",
		ObjectStore: "default",
	}
	objectStores := []config.ObjectStore{}
	objectStores = append(objectStores, *defaultobj)

	namespaces := []config.Namespace{}
	namespaces = append(namespaces, *defaultns)

	server := &config.Server{Dev: true}
	testConfig := &config.Configuration{ObjectStores: objectStores, Namespaces: namespaces, Server: *server}

	err = BootstrapDefaults(testConfig, db)
	if err != nil {
		t.Fatalf("Failed to bootstrap namespace: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}
	test.AssertEqual(t, objLoaded.Name, "default")
	test.AssertEqual(t, objLoaded.Description, "Default object store")
	test.AssertEqual(t, objLoaded.Endpoint, "http://localhost")
	test.AssertEqual(t, objLoaded.ObjectStoreType, "minio")
	test.AssertEqual(t, objLoaded.Profile, "")
	test.AssertEqual(t, objLoaded.Region, "")

	nsLoaded, err := db.GetNamespace("default")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}
	test.AssertEqual(t, nsLoaded.Name, "default")
	test.AssertEqual(t, nsLoaded.Description, "Default namespace")
	test.AssertEqual(t, nsLoaded.BucketName, "test-data")
}

func TestBootstrapExist(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	defaultobj := &config.ObjectStore{
		Name:        "default",
		Description: "Default object store",
		Type:        "minio",
		Endpoint:    "http://localhost",
	}
	defaultns := &config.Namespace{
		Name:        "default",
		Description: "Default namespace",
		Bucket:      "test-data",
		ObjectStore: "default",
	}
	objectStores := []config.ObjectStore{}
	objectStores = append(objectStores, *defaultobj)

	namespaces := []config.Namespace{}
	namespaces = append(namespaces, *defaultns)

	server := &config.Server{Dev: true}
	testConfig := &config.Configuration{ObjectStores: objectStores, Namespaces: namespaces, Server: *server}

	err = db.CreateNamespace(defaultns.Name, defaultns.Description, "default", "a-different-bucket-to-verify-run")
	if err != nil {
		t.Fatal("Failed to create namespace record")
	}

	_, err = db.GetNamespace("default")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}

	err = BootstrapDefaults(testConfig, db)
	if err != nil {
		t.Fatal("Failed to bootstrap namespace")
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}
	test.AssertEqual(t, objLoaded.Name, "default")
	test.AssertEqual(t, objLoaded.Description, "Default object store")
	test.AssertEqual(t, objLoaded.Endpoint, "http://localhost")
	test.AssertEqual(t, objLoaded.ObjectStoreType, "minio")
	test.AssertEqual(t, objLoaded.Profile, "")
	test.AssertEqual(t, objLoaded.Region, "")

	nsLoaded, err := db.GetNamespace("default")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}
	test.AssertEqual(t, nsLoaded.Name, "default")
	test.AssertEqual(t, nsLoaded.Description, "Default namespace")
	test.AssertEqual(t, nsLoaded.BucketName, "a-different-bucket-to-verify-run")
}
