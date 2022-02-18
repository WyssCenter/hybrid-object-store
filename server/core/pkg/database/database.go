package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	hossMigrations "github.com/gigantum/hoss-core/pkg/database/migrations"
)

// Database holds any database related data needed for interacting with the database
type Database struct {
	conn *pg.DB
}

// Load creates a connection to the database and applies the migrations
func Load() *Database {
	// Only load the migrations once, as running this multiple times will break the migrations
	if len(migrations.DefaultCollection.Migrations()) == 0 {
		// Initial database setup
		hossMigrations.Register0001()
		// Background Dataset Delete support: https://github.com/gigantum/hybrid-object-store/issues/227
		hossMigrations.Register0002()
		// Sync Policy support
		hossMigrations.Register0003()
	}

	db := &Database{}

	db.conn = pg.Connect(&pg.Options{
		Addr:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	})

	isDatabaseReady := false
	ctx := context.Background()
	for i := 0; i < 12; i++ {
		if err := db.conn.Ping(ctx); err == nil {
			isDatabaseReady = true
			break
		}

		logrus.Info("Database is not ready, sleeping")
		time.Sleep(5 * time.Second)
	}
	if !isDatabaseReady {
		logrus.Fatal("Couldn't connect to database after 60 seconds")
	}

	_, _, err := migrations.Run(db.conn, "init") // create migration schema
	if err != nil {
		db.conn.Close()
		logrus.WithError(err).Fatal("Couldn't initialize migrations schema")
	}

	oldVersion, newVersion, err := migrations.Run(db.conn, "up")
	if err != nil {
		db.conn.Close()
		logrus.WithError(err).Fatal("Couldn't apply database migrations")
	}

	if newVersion != oldVersion {
		fmt.Printf("Database migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("Database version is %d\n", oldVersion)
	}

	return db
}

// ValidateMembership checks whether a user is a member of a group
func (db *Database) ValidateMembership(username, groupName string) bool {
	user, err := db.GetOrCreateUser(username)
	if err != nil {
		return false
	}

	for _, membership := range user.Memberships {
		if groupName == membership.Group.GroupName {
			return true
		}
	}

	return false
}

// GetUserDefaultGroup returns the default group name for a user
func (db *Database) GetUserDefaultGroup(username string) string {
	return fmt.Sprintf("%s-%s", username, DEFAULT_GROUP_SUFFIX)
}

// GetOrCreateUser will get the User from the database or create a new one if it doesn't exist
func (db *Database) GetOrCreateUser(username string) (*User, error) {
	user := &User{}

	err := db.conn.Model(user).
		Where(`"user"."username" = ?`, username).
		Relation("Memberships.Group").
		Select()
	if err != nil {
		user.Username = username
		groupName := db.GetUserDefaultGroup(username)

		_, err := db.conn.Model(user).Insert()
		if err != nil {
			return nil, ConvertError(err)
		}

		_, err = db.GetOrCreateGroup(groupName)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create new user's default group")
		}

		err = db.UpdateGroupMembership(username, groupName)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't add new user to default group")
		}

		// re-select user so their membership is included
		err = db.conn.Model(user).
			Where(`"user"."username" = ?`, username).
			Relation("Memberships.Group").
			Select()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get new user")
		}
	}
	return user, nil
}

// ListUsers will list all users with limit/offset for pagination
func (db *Database) ListUsers(limit int, offset int) ([]*User, error) {
	var users []*User

	err := db.conn.Model(&users).Limit(limit).Offset(offset).Relation("Memberships.Group").Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return users, nil
}

// GetGroup will get the Group from the database
func (db *Database) GetOrCreateGroup(groupName string) (*Group, error) {
	group := &Group{}

	err := db.conn.Model(group).
		Where(`"group"."group_name" = ?`, groupName).
		Relation("Permissions.Dataset").
		Relation("Memberships.User").
		Select()
	if err != nil {
		group.GroupName = groupName
		_, err := db.conn.Model(group).Insert()
		if err != nil {
			return nil, ConvertError(err)
		}
	}

	return group, nil
}

// DeleteGroup deletes a group from the database
func (db *Database) DeleteGroup(groupName string) error {
	entry := Group{
		GroupName: groupName,
	}

	_, err := db.conn.Model(&entry).Where(`group_name = ?group_name`).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// UpdateGroupMembership adds/updates a user's membership to a group
func (db *Database) UpdateGroupMembership(username, groupName string) error {
	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		return err
	}

	user, err := db.GetOrCreateUser(username)
	if err != nil {
		return err
	}

	entry := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}

	// Perform an insert or update
	_, err = db.conn.Model(&entry).
		OnConflict("(group_id, user_id) DO UPDATE").
		Insert()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// RemoveGroupMembership removes a user's membership to a group
func (db *Database) RemoveGroupMembership(username, groupName string) error {
	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		return err
	}

	user, err := db.GetOrCreateUser(username)
	if err != nil {
		return err
	}

	entry := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}

	// Perform an insert or update
	_, err = db.conn.Model(&entry).
		Where(`group_id = ?group_id AND user_id = ?user_id`).
		Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// GetGroupsByDataset gets the Groups and their permissions for the given dataset
func (db *Database) GetGroupsByDataset(namespace *Namespace, name string) ([]*Permission, error) {
	var dataset Dataset
	err := db.conn.Model(&dataset).
		Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, name).
		Relation("Permissions.Group").
		Select()
	if err != nil {
		// DP ???: Should this return an empty list or an error if the dataset doesn't exist?
		err = ConvertError(err)
		if err == ErrNotFound {
			return []*Permission{}, nil
		} else {
			return nil, err
		}
	}

	return dataset.Permissions, nil
}

// CreateSyncConfiguration creates a new SyncConfiguration entry in the database
// Note: triggers an update to the LastModified SyncConfigurationMeta timestamp
func (db *Database) CreateSyncConfiguration(namespace *Namespace, targetCoreService, targetNamespace, syncType string) error {
	sc := SyncConfiguration{
		SourceNamespaceId: namespace.Id,
		TargetCoreService: targetCoreService,
		TargetNamespace:   targetNamespace,
		SyncType:          syncType,
	}
	_, err := db.conn.Model(&sc).Insert()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// DeleteSyncConfiguration deletes a SyncConfiguration from the database
// Note: there is no error if you delete a non-existent sync configuration
// Note: triggers an update to the LastModified SyncConfigurationMeta timestamp if there is a change in the database
func (db *Database) DeleteSyncConfiguration(namespace *Namespace, targetCoreService, targetNamespace string) error {
	sc := SyncConfiguration{
		SourceNamespaceId: namespace.Id,
		TargetCoreService: targetCoreService,
		TargetNamespace:   targetNamespace,
	}
	_, err := db.conn.Model(&sc).
		Where("source_namespace_id = ?source_namespace_id").
		Where("target_core_service = ?target_core_service").
		Where("target_namespace = ?target_namespace").
		Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// GetNamespaceSyncTargets returns the SyncConfigurations that the given Namespace is configured to send mnotifications to
func (db *Database) GetNamespaceSyncTargets(namespace *Namespace) ([]*SyncConfiguration, error) {
	syncConfigurations := []*SyncConfiguration{}

	err := db.conn.Model(&syncConfigurations).
		Where("source_namespace_id = ?", namespace.Id).Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return syncConfigurations, nil
}

// GetSyncConfigurations returns all of the SyncConfigurations managed by the Core Service
func (db *Database) GetSyncConfigurations() ([]*SyncConfiguration, error) {
	syncConfigurations := []*SyncConfiguration{}

	err := db.conn.Model(&syncConfigurations).
		Relation("SourceNamespace").Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return syncConfigurations, nil
}

// GetSyncEnabledPolicies returns the Dataset.RootDirectory and Dataset.SyncPolicy for datasets in the namespace that are sync enabled
func (db *Database) GetSyncEnabledPolicies(namespace *Namespace) (map[string]string, error) {
	datasets := []*Dataset{}

	err := db.conn.Model(&datasets).
		Column("root_directory", "sync_policy").
		Where("sync_enabled").
		Where("namespace_id = ?", namespace.Id).
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	prefixPolicy := map[string]string{}
	for _, ds := range datasets {
		prefixPolicy[ds.RootDirectory] = ds.SyncPolicy
	}

	return prefixPolicy, nil
}

// GetLastSyncUpdated returns the last updated timestamp for any change to the sync configuration managed by the Core Service
func (db *Database) GetLastSyncUpdated() (time.Time, error) {
	var lastUpdated time.Time
	err := db.conn.Model((*SyncConfigurationMeta)(nil)).
		Column("last_updated").Limit(1).Select(&lastUpdated)
	if err != nil {
		return time.Time{}, ConvertError(err)
	}

	return lastUpdated, nil
}

// ListDatasetsInNamespace will list all datasets within a namespace. Since permissions will not be
// applied, this should be infrequently used. The main reason it was created is to support patching events
// for all the datasets in a minio object store
func (db *Database) ListDatasetsInNamespace(namespace *Namespace) ([]*Dataset, error) {
	var datasets []*Dataset

	err := db.conn.Model(&datasets).
		Where(`"dataset"."namespace_id" = ?`, namespace.Id).
		Order("id ASC").
		Relation("Namespace.ObjectStore").Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return datasets, nil
}
