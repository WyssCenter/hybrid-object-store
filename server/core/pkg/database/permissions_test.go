package database

import (
	"fmt"
	"testing"
)

func TestGetDatasetsByGroupExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	perms, err := db.GetPermissionsByGroup(objLoaded, "test_user-hoss-default-group", false)
	if err != nil {
		t.Fatal("Expected no error but get datasets by group failed")
	}

	if perms == nil {
		t.Fatal("Result is nil")
	}

	if len(perms) != 1 {
		t.Fatalf("Expected one result from get datasets by group which returned %d", len(perms))
	}

	if perms[0].Dataset.Name != "test_dataset" {
		t.Fatal("Unexpected dataset returned by get datasets by group")
	}
}

func TestGetUsersWithPermissionsToDataset(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	usernames, err := db.GetUsersWithPermissionsToDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get datasets by group failed")
	}

	if usernames == nil {
		t.Fatal("usernames is nil")
	}

	if len(usernames) != 1 {
		t.Fatalf("Expected one result but got %d", len(usernames))
	}

	if usernames[0] != "test_user" {
		t.Fatal("Unexpected username")
	}
}

func TestGetDatasetsByGroupNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	perms, err := db.GetPermissionsByGroup(objLoaded, "test_group2", false)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal("Expected no error but get datasets by group failed")
	}

	if len(perms) != 0 {
		t.Fatalf("Expected no results from get datasets by group which returned %d", len(perms))
	}
}

func TestGetDatasetsByUserExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	perms, err := db.GetPermissionsByUser(objLoaded, "test_user", false)
	if err != nil {
		t.Fatalf("Expected no error but get datasets by user failed: %v", err)
	}

	if len(perms) != 1 {
		t.Fatalf("Expected one result from get datasets by user which returned %d", len(perms))
	}

	if perms[0].Dataset.Name != "test_dataset" {
		t.Fatal("Unexpected dataset returned by get datasets by user")
	}

	if perms[0].Dataset.Namespace.Name != "test_namespace" {
		t.Fatal("Unexpected namespace returned by get datasets by user")
	}
}

func TestGetDatasetsByUserNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	_, err = db.GetPermissionsByUser(objLoaded, "test_user1", false)
	if err != nil {
		t.Fatal("Expected no error but get datasets by user failed: ", err.Error())
	}
}

func TestGetGroupPermissionsExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	perms, err := db.GetPermissionsByUser(objLoaded, "test_user", false)
	if err != nil {
		t.Fatal("Expected no error but get datasets by group failed")
	}

	if len(perms) != 1 {
		t.Fatalf("Expected one result from get datasets by group which returned %d", len(perms))
	}

	if perms[0].Dataset.Name != "test_dataset" {
		t.Fatal("Unexpected permission returned by get datasets by group")
	}
}

func TestGetGroupPermissionsNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	objLoaded, err := db.GetObjectStore("default")
	if err != nil {
		t.Fatal("Expected no error but get object store failed")
	}

	perms, err := db.GetPermissionsByUser(objLoaded, "test_group2", false)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal("Expected no error but get datasets by group failed")
	}

	if len(perms) != 0 {
		t.Fatalf("Expected no group permissions but there are %d", len(perms))
	}
}

func TestGetDatasetPermissions(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	perm, err := db.GetDatasetPermissions(ns, "test_dataset", "test_user-hoss-default-group")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fatal("Expected no error but get dataset permissions failed")
	}

	if perm != PERM_READ_WRITE {
		t.Fatalf("Unexpected get dataset permissions value: %s", perm)
	}
}

func TestGetDatasetPermissionsNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	perm, err := db.GetDatasetPermissions(ns, "test_dataset", "test_group1")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fatal("Expected no error but get dataset permissions failed")
	}

	if perm != "" {
		t.Fatalf("Unexpected get dataset permissions value: %s", perm)
	}
}

func TestGetDatasetPermissionsNotExistingGroup(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	_, err = db.GetDatasetPermissions(ns, "test_dataset", "test_group2")
	if err != nil {
		t.Fatal("Expected no error but get dataset permissions failed: ", err.Error())
	}
}
