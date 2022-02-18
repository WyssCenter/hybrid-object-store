package database

import (
	"testing"

	"github.com/gigantum/hoss-auth/pkg/test"
)

func TestCreateUser(t *testing.T) {
	db, user_new, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateUser(*user_new)
	if err != nil {
		t.Fatal("Expected no error but create user failed: ", err.Error())
	}

	user, err := db.GetUser(user_new.Username)
	if err != nil {
		t.Fatal("Unable to get newly created user: ", err.Error())
	}

	if len(user.Memberships) != 1 {
		t.Fatalf("Expected one membership for this user but there were %d", len(user.Memberships))
	}

	if user.Memberships[0].Group.GroupName != "test_user1-hoss-default-group" {
		t.Fatalf("Expected user's group name to match the default name, but the group name is %s", user.Memberships[0].Group.GroupName)
	}
}

func TestCreateOrUpdateUser(t *testing.T) {
	db, user_new, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateOrUpdateUser(*user_new)
	if err != nil {
		t.Fatal("Expected no error but create user failed: ", err.Error())
	}

	user, err := db.GetUser(user_new.Username)
	if err != nil {
		t.Fatal("Unable to get newly created user: ", err.Error())
	}

	test.AssertEqual(t, user.Email, "test_user1@gmail.com")
	test.AssertEqual(t, user.GivenName, "Fake")
	test.AssertEqual(t, user.FamilyName, "Person")
	test.AssertEqual(t, user.Username, "test_user1")
	test.AssertEqual(t, user.EmailVerified, true)
	test.AssertEqual(t, user.Role, "user")
	test.AssertEqual(t, user.Subject, "sub")

	if len(user.Memberships) != 1 {
		t.Fatalf("Expected one membership for this user but there were %d", len(user.Memberships))
	}

	if user.Memberships[0].Group.GroupName != "test_user1-hoss-default-group" {
		t.Fatalf("Expected user's group name to match the default name, but the group name is %s", user.Memberships[0].Group.GroupName)
	}

	user_new.Role = "admin"
	user_new.GivenName = "Really Fake"

	err = db.CreateOrUpdateUser(*user_new)
	if err != nil {
		t.Fatal("Expected no error but create user failed: ", err.Error())
	}

	user, err = db.GetUser(user_new.Username)
	if err != nil {
		t.Fatal("Unable to get newly created user: ", err.Error())
	}
	test.AssertEqual(t, user.Email, "test_user1@gmail.com")
	test.AssertEqual(t, user.GivenName, "Really Fake")
	test.AssertEqual(t, user.FamilyName, "Person")
	test.AssertEqual(t, user.Username, "test_user1")
	test.AssertEqual(t, user.EmailVerified, true)
	test.AssertEqual(t, user.Role, "admin")
	test.AssertEqual(t, user.Subject, "sub")

}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, _, user_exists, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateUser(*user_exists)
	if err == nil {
		t.Fatal("Expected an error but create user succeeded")
	}
}

func TestGetUserExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	user, err := db.GetUser("test_user")
	if err != nil {
		t.Fatal("Expected no error but get user failed: ", err.Error())
	}

	if user.FullName != "Test User" {
		t.Fatalf("User returned from database does not have same full name as input: %s", user.FullName)
	}
}

func TestGetUserNotExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetUser("test_user1")
	if err == nil {
		t.Fatal("Expected an error but get user succeeded: ")
	}
}

func TestCreateGroupAlreadyExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.CreateGroup("test_user-hoss-default-group", "test description")
	if err == nil {
		t.Fatal("Expected create group to fail, but it succeeded")
	}
}

func TestCreateGroupNotExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	group, err := db.CreateGroup("test_group2", "fake new group")
	if err != nil {
		t.Fatal("Expected create group to succeed, but it failed", err.Error())
	}

	if group.Description != "fake new group" {
		t.Fatal("Description does not match: ", group.Description)
	}
}

func TestGetGroupExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	group, err := db.GetGroup("test_user-hoss-default-group")
	if err != nil {
		t.Fatal("Expected get group to succeed, but it failed", err.Error())
	}

	if group.Description != "test description" {
		t.Fatal("Description does not match: ", group.Description)
	}
}

func TestGetGroupNotExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetGroup("test_group2")
	if err == nil {
		t.Fatal("Expected create group to fail, but it succeeded")
	}
}

func TestDeleteGroupExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.DeleteGroup("test_group")
	if err != nil {
		t.Fatal("Expected delete group to succeed, but it failed: ", err.Error())
	}

	_, err = db.GetGroup("test_group")
	if err == nil {
		t.Fatal("Expected deleted group to no longer exist, but it does")
	}
}

func TestDeleteGroupNotExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.DeleteGroup("test_group2")
	if err != nil {
		t.Fatal("Expected delete group to succeed, but it failed: ", err.Error())
	}
}

func TestUpdateGroupMembershipExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.UpdateGroupMembership("test_user", "test_empty_group")
	if err != nil {
		t.Fatal("Expected add user to group to succeed, but it failed: ", err.Error())
	}

	group, err := db.GetGroup("test_empty_group")
	if err != nil {
		t.Fatal("Error getting group after adding a user: ", err.Error())
	}
	if len(group.Memberships) != 1 {
		t.Fatal("Expected group to have one user, but it had: ", len(group.Memberships))
	}
}

func TestUpdateGroupMembershipNotExistsUser(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.UpdateGroupMembership("test_user1", "test_empty_group")
	if err == nil {
		t.Fatal("Expected add user to group to fail, but it succeeded")
	}

	group, err := db.GetGroup("test_empty_group")
	if err != nil {
		t.Fatal("Error getting group after adding a user: ", err.Error())
	}
	if len(group.Memberships) != 0 {
		t.Fatal("Expected group to have no users, but it had: ", len(group.Memberships))
	}
}

func TestUpdateGroupMembershipNotExistsGroup(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.UpdateGroupMembership("test_user", "test_group2")
	if err == nil {
		t.Fatal("Expected add user to group to fail, but it succeeded")
	}
}

func TestRemoveGroupMembershipExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user", "test_user-hoss-default-group")
	if err != nil {
		t.Fatal("Expected remove user to group to succeed, but it failed: ", err.Error())
	}

	group, err := db.GetGroup("test_user-hoss-default-group")
	if err != nil {
		t.Fatal("Error getting group after adding a user: ", err.Error())
	}
	if len(group.Memberships) != 0 {
		t.Fatal("Expected group to have no users, but it had: ", len(group.Memberships))
	}
}

func TestRemoveGroupMembershipNotExistsUser(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user1", "test_user-hoss-default-group")
	if err == nil {
		t.Fatal("Expected remove user to group to fail, but it succeeded")
	}

	group, err := db.GetGroup("test_user-hoss-default-group")
	if err != nil {
		t.Fatal("Error getting group after adding a user: ", err.Error())
	}
	if len(group.Memberships) != 1 {
		t.Fatal("Expected group to have one user, but it had: ", len(group.Memberships))
	}
}

func TestRemoveGroupMembershipNotExistsGroup(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.RemoveGroupMembership("test_user", "test_group2")
	if err == nil {
		t.Fatal("Expected remove user to group to fail, but it succeeded")
	}
}

func TestCreatePATUserExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	pat, err := db.CreatePAT("test_user", "test description")
	if err != nil {
		t.Fatal("Expected no error but create PAT failed: ", err.Error())
	}

	if pat.Description != "test description" {
		t.Fatal("PAT description is incorrect: ", pat.Description)
	}
}

func TestCreatePATUserDoesNotExist(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.CreatePAT("test_user1", "test description")
	if err == nil {
		t.Fatal("Expected an error but create PAT succeeded")
	}
}

func TestGetPATExists(t *testing.T) {
	db, _, _, pat_exists, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetPAT(pat_exists.PAT)
	if err != nil {
		t.Fatal("Expected no error but get PAT failed: ", err.Error())
	}
}

func TestGetPATDoesNotExist(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetPAT("fakepatid")
	if err == nil {
		t.Fatal("Expected error but get PAT succeeded")
	}
}

func TestListPATUserExists(t *testing.T) {
	db, _, _, pat_exists, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	pat_list, err := db.ListPAT("test_user")
	if err != nil {
		t.Fatal("Expected no error but list PAT failed: ", err.Error())
	}

	if pat_list[0].PAT != pat_exists.PAT {
		t.Fatal("Listed PAT does not match existing PAT")
	}
}

func TestListPATUserDoesNotExist(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.ListPAT("test_user1")
	if err == nil {
		t.Fatal("Expected an error but list PAT succeeded")
	}
}

func TestListPATEmpty(t *testing.T) {
	db, _, _, pat_exists, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.DeletePAT(pat_exists.Id, "test_user")
	if err != nil {
		t.Fatal("Expected no error but delete PAT failed: ", err.Error())
	}

	pat_list, err := db.ListPAT("test_user")
	if err != nil {
		t.Fatal("Expected no error but list PAT failed: ", err.Error())
	}

	if pat_list != nil {
		t.Fatal("Listed PAT expected to be nil")
	}
}

func TestDeletePATExists(t *testing.T) {
	db, _, _, pat_exists, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.DeletePAT(pat_exists.Id, "test_user")
	if err != nil {
		t.Fatal("Expected no error but delete PAT failed: ", err.Error())
	}
}

func TestDeletePATNotExists(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.DeletePAT(int64(1234), "test_user")
	if err == nil {
		t.Fatal("Expected an error but delete PAT succeeded")
	}
	if err.Error() != "record not found" {
		t.Fatal("Expected a not found error")

	}
}

func TestDeletePATExistsWrongUser(t *testing.T) {
	// test deleting a PAT that isn't yours
	db, user_new, _, pat_exists, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateOrUpdateUser(*user_new)
	if err != nil {
		t.Fatal("Expected no error but create user failed: ", err.Error())
	}

	err = db.DeletePAT(pat_exists.Id, "test_user1")
	if err == nil {
		t.Fatal("Expected an error but delete PAT succeeded")
	}
	if err.Error() != "record not found" {
		t.Fatal("Expected a not found error")

	}
}

func TestGetUsersByEmail(t *testing.T) {
	db, _, _, _, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// Add another user with same email
	user := User{
		Username:      "other_account",
		FullName:      "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         "test_user@gmail.com",
		EmailVerified: true,
		Subject:       "sub",
		Role:          "user",
	}
	_, err = db.conn.Model(&user).Insert()
	if err != nil {
		t.Fatal("failed to create test user: ", err.Error())
	}

	// Add another user with different email
	user = User{
		Username:      "other_user",
		FullName:      "Other User",
		GivenName:     "Other",
		FamilyName:    "User",
		Email:         "other@gmail.com",
		EmailVerified: true,
		Subject:       "sub",
		Role:          "user",
	}
	_, err = db.conn.Model(&user).Insert()
	if err != nil {
		t.Fatal("failed to create test user: ", err.Error())
	}

	users, err := db.GetUsersByEmail("test_user@gmail.com")
	if err != nil {
		t.Fatalf("Expected no error but get user by email failed: %s", err.Error())
	}

	if len(*users) != 2 {
		t.Fatalf("Expected two users but got %v", len(*users))
	}

	users, err = db.GetUsersByEmail("other@gmail.com")
	if err != nil {
		t.Fatalf("Expected no error but get user by email failed: %s", err.Error())
	}

	if len(*users) != 1 {
		t.Fatalf("Expected one users but got %v", len(*users))
	}

}
