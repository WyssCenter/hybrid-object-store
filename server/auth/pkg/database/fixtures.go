package database

import (
	"os"
	"testing"

	"github.com/gigantum/hoss-auth/pkg/test"
	userinfo "github.com/gigantum/hoss-auth/pkg/userinfo"
	"github.com/go-pg/migrations/v8"
	"github.com/pkg/errors"
)

// SetupDatabaseTest sets the database for testing
func SetupDatabaseTest(t *testing.T) (*Database, *userinfo.UserInfo, *userinfo.UserInfo, *PersonalAccessToken, error) {
	// Set env vars
	err := test.LoadEnvFile("~/.hoss/.env")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err := os.Setenv("POSTGRES_DB", "hoss_auth"); err != nil {
		return nil, nil, nil, nil, err
	}
	if err := os.Setenv("POSTGRES_HOST", "localhost:5432"); err != nil {
		return nil, nil, nil, nil, err
	}

	db := Load()

	// create test userinfo
	emailVerified := true
	testUserNew := userinfo.UserInfo{
		Subject:       "sub",
		FullName:      "Fake Person",
		GivenName:     "Fake",
		FamilyName:    "Person",
		Username:      "test_user1",
		Email:         "test_user1@gmail.com",
		EmailVerified: &emailVerified,
		Role:          "user",
	}

	testUserExists := userinfo.UserInfo{
		Subject:       "sub",
		FullName:      "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Username:      "test_user",
		Email:         "test_user@gmail.com",
		EmailVerified: &emailVerified,
		Role:          "user",
	}

	// add existing test user
	user := User{
		Username:      testUserExists.Username,
		FullName:      testUserExists.FullName,
		GivenName:     testUserExists.GivenName,
		FamilyName:    testUserExists.FamilyName,
		Email:         testUserExists.Email,
		EmailVerified: *testUserExists.EmailVerified,
		Subject:       testUserExists.Subject,
		Role:          testUserExists.Role,
	}
	_, err = db.conn.Model(&user).Insert()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Could not create test user")
	}

	// test user's default group
	group := Group{GroupName: "test_user-hoss-default-group", Description: "test description"}
	_, err = db.conn.Model(&group).Insert()
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Could not create test user's default group")
	}

	// test user's membership in default group
	membership := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}
	if _, err := db.conn.Model(&membership).Insert(); err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Could not create test membership for user and default group")
	}

	// add empty group
	group1 := Group{GroupName: "test_empty_group", Description: "fake description"}
	if _, err := db.conn.Model(&group1).Insert(); err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Could not create empty group")
	}

	// personal access token for test user
	pat := PersonalAccessToken{
		Description: "test description",
		OwnerId:     user.Id,
	}
	if _, err := db.conn.Model(&pat).Insert(); err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Could not create test pat")
	}

	// Set cleanup
	t.Cleanup(func() {
		TeardownDatabaseTest(t, db)
	})

	return db, &testUserNew, &testUserExists, &pat, nil
}

// TeardownDatabaseTest gracefully tries to remove all data created by a test
func TeardownDatabaseTest(t *testing.T, c *Database) {
	// Reset the migrations by reverting all migrations that have been applied
	_, _, err := migrations.Run(c.conn, "reset")
	if err != nil {
		panic(err)
	}
}
