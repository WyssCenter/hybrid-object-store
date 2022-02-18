package database

import (
	"testing"

	"github.com/gigantum/hoss-core/pkg/test"
)

func TestCreateNamespace(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateNamespace("my_namespace", "My description", "default", "data")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	nsLoaded, err := db.GetNamespace("my_namespace")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}

	test.AssertEqual(t, nsLoaded.Name, "my_namespace")
	test.AssertEqual(t, nsLoaded.Description, "My description")
	test.AssertEqual(t, nsLoaded.BucketName, "data")
	test.AssertEqual(t, nsLoaded.ObjectStore.Name, "default")
	test.AssertEqual(t, nsLoaded.ObjectStore.Description, "Default object store")
	test.AssertEqual(t, nsLoaded.ObjectStore.Endpoint, "http://localhost")
	test.AssertEqual(t, nsLoaded.ObjectStore.ObjectStoreType, OBJECT_STORE_TYPE_MINIO)
	test.AssertEqual(t, nsLoaded.ObjectStore.Profile, "")
	test.AssertEqual(t, nsLoaded.ObjectStore.Region, "")
}

func TestCreateNamespaceExists(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateNamespace("my_namespace", "My description", "default", "data")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	err = db.CreateNamespace("my_namespace", "My description 2", "default", "data 2")
	if err == nil {
		t.Fatal("Created namespace again when an exception should have been raised")
	}

	test.AssertEqual(t, err.Error(), "record already exists")

}

func TestGetNamespaceDoesNotExist(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetNamespace("my_namespace_does_not_exist")
	if err == nil {
		t.Fatal("Expected error because namespace does not exist")
	}

	test.AssertEqual(t, err.Error(), "record not found")

}

func TestDeleteNamespace(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateNamespace("my_namespace", "My description", "default", "data")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	_, err = db.GetNamespace("my_namespace")
	if err != nil {
		t.Fatal("Expected namespace to exist")
	}

	err = db.DeleteNamespace("my_namespace")
	if err != nil {
		t.Fatal("Expected no error but delete namespace failed")
	}

	_, err = db.GetNamespace("my_namespace")
	if err == nil {
		t.Fatal("Expected error because namespace should not exist")
	}
	test.AssertEqual(t, err.Error(), "record not found")
}

func TestListNamespaces(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateNamespace("ns1", "My description", "default", "data1")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	err = db.CreateNamespace("ns2", "My description", "default", "data2")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	err = db.CreateNamespace("ns3", "My description", "default", "data3")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	err = db.CreateNamespace("ns4", "My description", "default", "data4")
	if err != nil {
		t.Fatal("Failed to create namespace")
	}

	namespaces, err := db.ListNamespaces(25, 0)
	if err != nil {
		t.Fatal("Failed to list namespaces")
	}
	if len(namespaces) != 5 {
		t.Fatalf("Length of namespaces is wrong, %v != 5", len(namespaces))
	}
	test.AssertEqual(t, namespaces[1].Name, "ns1")
	test.AssertEqual(t, namespaces[4].Name, "ns4")

	namespaces, err = db.ListNamespaces(25, 1)
	if err != nil {
		t.Fatal("Failed to list namespaces")
	}
	if len(namespaces) != 4 {
		t.Fatalf("Length of namespaces is wrong, %v != %v", len(namespaces), 4)
	}
	test.AssertEqual(t, namespaces[0].Name, "ns1")
	test.AssertEqual(t, namespaces[3].Name, "ns4")

	namespaces, err = db.ListNamespaces(2, 2)
	if err != nil {
		t.Fatal("Failed to list namespaces")
	}
	if len(namespaces) != 2 {
		t.Fatalf("Length of namespaces is wrong, %v != 2", len(namespaces))
	}
	test.AssertEqual(t, namespaces[0].Name, "ns2")
	test.AssertEqual(t, namespaces[1].Name, "ns3")

}
