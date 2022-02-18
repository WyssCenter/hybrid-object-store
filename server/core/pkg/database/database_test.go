package database

import (
	"testing"
	"time"

	"github.com/gigantum/hoss-core/pkg/test"
)

func TestGetOrCreateUserExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	user, err := db.GetOrCreateUser("test_user")
	if err != nil {
		t.Fatalf("Expected no error but get or create user failed: %v", err)
	}

	if len(user.Memberships) != 1 {
		t.Fatal("Expected user to have one group but it has: ", len(user.Memberships))
	}

	test.AssertEqual(t, user.Memberships[0].Group.GroupName, "test_user-hoss-default-group")
}

func TestGetOrCreateUserNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	user, err := db.GetOrCreateUser("test_user1")
	if err != nil {
		t.Fatalf("Expected no error but get or create user failed: %v", err)
	}

	if len(user.Memberships) != 1 {
		t.Fatal("Expected user to have one group but it has: ", len(user.Memberships))
	}

	test.AssertEqual(t, user.Memberships[0].Group.GroupName, "test_user1-hoss-default-group")
}

func TestGetOrCreateGroupExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	group, err := db.GetOrCreateGroup("test_user-hoss-default-group")
	if err != nil {
		t.Fatalf("Expected no error but get group failed: %v", err)
	}

	if len(group.Memberships) != 1 {
		t.Fatal("Expected group to have one member but it has: ", len(group.Memberships))
	}
}

func TestGetOrCreateGroupNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	group, err := db.GetOrCreateGroup("test_group2")
	if err != nil {
		t.Fatalf("Expected no error but get group failed: %v", err)
	}

	if len(group.Memberships) != 0 {
		t.Fatal("Expected group to have no members but it has: ", len(group.Memberships))
	}
}

func TestUpdateGroupMembershipExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.UpdateGroupMembership("test_user", "test_group1")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed: ", err.Error())
	}

	group, err := db.GetOrCreateGroup("test_group1")
	if err != nil {
		t.Fatal("Error getting updated group: ", err.Error())
	}

	if len(group.Memberships) != 1 {
		t.Fatal("Expected group to have one member but it has: ", len(group.Memberships))
	}
}

func TestUpdateGroupMembershipNotExistingGroup(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.UpdateGroupMembership("test_user", "test_group2")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed: ", err.Error())
	}
}

func TestUpdateGroupMembershipNotExistingUser(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.UpdateGroupMembership("test_user2", "test_group1")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed: ", err.Error())
	}

	group, err := db.GetOrCreateGroup("test_group1")
	if err != nil {
		t.Fatal("Error getting updated group: ", err.Error())
	}

	if len(group.Memberships) != 1 {
		t.Fatal("Expected group to have one member but it has: ", len(group.Memberships))
	}
}

func TestRemoveGroupMembershipExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user", "test_group")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed: ", err.Error())
	}

	group, err := db.GetOrCreateGroup("test_group")
	if err != nil {
		t.Fatal("Error getting updated group: ", err.Error())
	}
	if len(group.Memberships) != 0 {
		t.Fatal("Expected group to have no members but it has: ", len(group.Memberships))
	}
}

func TestRemoveGroupMembershipNotExistingGroup(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user", "test_group2")
	if err != nil {
		t.Fatal("Expected no error but remove dataset permissions failed: ", err.Error())
	}
}

func TestRemoveGroupMembershipNotExistingUser(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user2", "test_group1")
	if err != nil {
		t.Fatal("Expected no error but update dataset permissions failed: ", err.Error())
	}
}

func TestGetGroupsByDatasetExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	perms, err := db.GetGroupsByDataset(ns, "test_dataset")
	if err != nil {
		t.Fatal("Expected no error but get groups by dataset failed: ", err.Error())
	}

	if len(perms) != 1 {
		t.Fatalf("Expected one result from get groups by dataset which returned %d", len(perms))
	}

	if perms[0].Group.GroupName != "test_user-hoss-default-group" {
		t.Fatal("Unexpected group returned by get groups by dataset")
	}
}

func TestGetGroupsByDatasetNotExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	perms, err := db.GetGroupsByDataset(ns, "test_dataset1")
	if err != nil {
		t.Fatal("Expected no error but get groups by dataset failed: ", err.Error())
	}

	if len(perms) != 0 {
		t.Fatalf("Expected no results from get users by dataset which returned %d", len(perms))
	}
}

func TestCreateSyncConfiguration(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	ts := time.Now()
	lastUpdate, err := db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ts.After(lastUpdate), true)

	err = db.CreateSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace", "simplex")
	if err != nil {
		t.Fatalf("Expected no error but create sync configuration failed: %v", err)
	}

	lastUpdate, err = db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ts.Before(lastUpdate), true)

	syncConfigs, err := db.GetSyncConfigurations()
	if err != nil {
		t.Fatalf("Expected no error but get sync configuration failed: %v", err)
	}

	test.AssertEqual(t, len(syncConfigs), 1)
	test.AssertEqual(t, syncConfigs[0].SourceNamespaceId, ns.Id)
	test.AssertEqual(t, syncConfigs[0].TargetCoreService, "http://localhost/core/v1")
	test.AssertEqual(t, syncConfigs[0].TargetNamespace, "target_namespace")
}

func TestCreateSyncConfigurationExisting(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace", "simplex")
	if err != nil {
		t.Fatalf("Expected no error but create sync configuration failed: %v", err)
	}

	err = db.CreateSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace", "simplex")
	if err == nil {
		t.Fatal("Expected error but create sync configuration succeeded")
	}

	if err != ErrExists {
		t.Fatalf("Expected already exists error but got a different type: %v", err)
	}
}

func TestDeleteSyncConfiguration(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace", "simplex")
	if err != nil {
		t.Fatalf("Expected no error but create sync configuration failed: %v", err)
	}

	ts := time.Now()

	err = db.DeleteSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace")
	if err != nil {
		t.Fatalf("Expected no error but delete sync configuration failed: %v", err)
	}

	lastUpdate, err := db.GetLastSyncUpdated()
	if err != nil {
		t.Fatal("Expected no error but get last sync updated failed: ", err.Error())
	}

	test.AssertEqual(t, ts.Before(lastUpdate), true)

	syncConfigs, err := db.GetSyncConfigurations()
	if err != nil {
		t.Fatalf("Expected no error but get sync configuration failed: %v", err)
	}

	test.AssertEqual(t, len(syncConfigs), 0)
}

func TestDeleteSyncConfigurationNotExists(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	ns, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.DeleteSyncConfiguration(ns, "http://localhost/core/v1", "target_namespace")
	if err != nil {
		t.Fatalf("Expected no error but delete sync configuration failed: %v", err)
	}
}

func TestDatasetsInANamespace(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	nsTest, err := db.GetNamespace("test_namespace")
	if err != nil {
		t.Fatal("Failed to get namespace")
	}

	err = db.CreateDataset(nsTest, "user-ds", "my dataset", "/user-ds", "other_user")
	if err != nil {
		t.Fatalf("failed to create other dataset: %v", err)
	}

	err = db.CreateDataset(nsTest, "other-ds", "my dataset", "/other-ds", "test_user")
	if err != nil {
		t.Fatalf("failed to create other dataset: %v", err)
	}

	dsList, err := db.ListDatasetsInNamespace(nsTest)
	if err != nil {
		t.Fatal("Failed to list datasets in namespace")
	}
	if len(dsList) != 3 {
		t.Fatal("incorrect number of items in dataset list")
	}

}
