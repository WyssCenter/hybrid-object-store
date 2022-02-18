package database

import (
	"os"
	"testing"
	"time"

	"github.com/gigantum/hoss-core/pkg/test"
	"github.com/go-pg/migrations/v8"
	"github.com/pkg/errors"
)

// SetupDatabaseTest sets the database for testing
func SetupDatabaseTest(t *testing.T) (*Database, error) {
	// Set env vars
	err := test.LoadEnvFile("~/.hoss/.env")
	if err != nil {
		return nil, err
	}

	if err := os.Setenv("POSTGRES_DB", "hoss_core"); err != nil {
		return nil, err
	}
	if err := os.Setenv("POSTGRES_HOST", "localhost:5432"); err != nil {
		return nil, err
	}

	db := Load()

	// Set cleanup
	t.Cleanup(func() {
		TeardownDatabaseTest(t, db)
	})

	// create test data
	user := User{Username: "test_user"}
	group := Group{GroupName: "test_user-hoss-default-group"}
	if _, err := db.conn.Model(&user).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create test user")
	}
	if _, err := db.conn.Model(&group).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create test user's default group")
	}

	membership := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}
	if _, err := db.conn.Model(&membership).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create test membership for user and default group")
	}

	group1 := Group{GroupName: "test_group1"}
	if _, err := db.conn.Model(&group1).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create empty test group")
	}

	// Create Test Object Store
	obj := &ObjectStore{Name: "default",
		Description:     "Default object store",
		Endpoint:        "http://localhost",
		ObjectStoreType: OBJECT_STORE_TYPE_MINIO}
	_, err = db.conn.Model(obj).Insert()
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't create objectstore record")
	}

	// Create Test Namespace
	ns := &Namespace{Name: "test_namespace",
		Description:   "My test namespace",
		ObjectStoreId: obj.Id,
		BucketName:    "data"}
	_, err = db.conn.Model(ns).Insert()
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't create namespace record")
	}

	dataset := Dataset{
		NamespaceId:   ns.Id,
		Name:          "test_dataset",
		Created:       time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		RootDirectory: "/test_dataset",
		OwnerId:       user.Id,
		SyncEnabled:   false,
		SyncType:      "simplex",
	}
	if _, err := db.conn.Model(&dataset).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create test dataset")
	}

	perm := Permission{
		GroupId:    group.Id,
		DatasetId:  dataset.Id,
		Permission: PERM_READ_WRITE,
	}
	if _, err := db.conn.Model(&perm).Insert(); err != nil {
		return nil, errors.Wrap(err, "Could not create test permission for default group")
	}

	return db, nil
}

// TeardownDatabaseTest gracefully tries to remove all data created by a test
func TeardownDatabaseTest(t *testing.T, c *Database) {
	// Reset the migrations by reverting all migrations that have been applied
	_, _, err := migrations.Run(c.conn, "reset")
	if err != nil {
		panic(err)
	}
}
