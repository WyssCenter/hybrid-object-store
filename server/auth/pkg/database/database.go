package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	hossMigrations "github.com/gigantum/hoss-auth/pkg/database/migrations"
	userinfo "github.com/gigantum/hoss-auth/pkg/userinfo"
)

// Database holds any database related data needed for interacting with the database
type Database struct {
	conn *pg.DB
}

// Load creates a connection to the database and applies the migrations
func Load() *Database {
	// Only load the migrations once, as running this multiple times will break the migrations
	if len(migrations.DefaultCollection.Migrations()) == 0 {
		hossMigrations.Register0001()
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

		logrus.Info("Database is not ready, auth svc sleeping")
		time.Sleep(5 * time.Second)
	}
	if !isDatabaseReady {
		logrus.Fatal("Couldn't connect to database from auth svc after 60 seconds")
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
	user, err := db.GetUser(username)
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

// GetUser will get the User from the database
func (db *Database) GetUser(username string) (*User, error) {
	var user User

	err := db.conn.Model(&user).
		Where(`"user"."username" = ?`, username).
		Relation("Memberships.Group").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &user, nil
}

// GetUsersByEmail will get the Users from the database by an email address
// since email addresses are not required to be unique, this may return multiple
// users
func (db *Database) GetUsersByEmail(email string) (*[]User, error) {
	var users []User

	err := db.conn.Model(&users).
		Where(`"user"."email" = ?`, email).
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &users, nil
}

// CreateOrUpdateUser updates a User if it exists or creates a new User if not
func (db *Database) CreateOrUpdateUser(userInfo userinfo.UserInfo) error {
	if _, err := db.GetUser(userInfo.Username); err == nil {
		// user exists, so update
		err = db.UpdateUser(userInfo)
		if err != nil {
			return ConvertError(err)
		}
	} else {
		// user doesn't exist, so create
		err = db.CreateUser(userInfo)
		if err != nil {
			return ConvertError(err)
		}
	}

	return nil
}

// UpdateUser updates an existing User record with its current user info
func (db *Database) UpdateUser(userInfo userinfo.UserInfo) error {
	user := User{
		Username:      userInfo.Username,
		FullName:      userInfo.FullName,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Email:         userInfo.Email,
		EmailVerified: *userInfo.EmailVerified,
		Subject:       userInfo.Subject,
		Role:          userInfo.Role,
	}

	_, err := db.conn.Model(&user).
		Column("username", "full_name", "given_name", "family_name", "email", "email_verified", "subject", "role").
		Where(`"user"."username" = ?`, userInfo.Username).
		Update()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// CreateUser creates a new User in the database
func (db *Database) CreateUser(userInfo userinfo.UserInfo) error {
	user := User{
		Username:      userInfo.Username,
		FullName:      userInfo.FullName,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Email:         userInfo.Email,
		EmailVerified: *userInfo.EmailVerified,
		Subject:       userInfo.Subject,
		Role:          userInfo.Role,
	}

	groupName := fmt.Sprintf("%s-%s", userInfo.Username, DEFAULT_GROUP_SUFFIX)

	// create new user and default group
	_, err := db.conn.Model(&user).Insert()
	if err != nil {
		return ConvertError(err)
	}

	group := Group{}
	group.GroupName = groupName
	group.Description = fmt.Sprintf("Personal group for user %s", userInfo.Username)
	_, err = db.conn.Model(&group).Insert()
	if err != nil {
		return ConvertError(err)
	}

	// Perform an insert or update
	_, err = db.UpdateGroupMembership(user.Username, groupName)
	if err != nil {
		return errors.Wrap(err, "couldn't create user membership in personal group")
	}

	return nil
}

// CreateGroup will create the Group in the database
func (db *Database) CreateGroup(groupName string, description string) (*Group, error) {
	if strings.HasSuffix(groupName, DEFAULT_GROUP_SUFFIX) || len(groupName) > 64 {
		return nil, ErrInvalidGroupName
	}

	group := Group{}
	group.GroupName = groupName
	group.Description = description

	_, err := db.conn.Model(&group).Insert()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &group, nil
}

// GetGroup will get the Group from the database
func (db *Database) GetGroup(groupName string) (*Group, error) {
	group := Group{}

	err := db.conn.Model(&group).
		Where(`"group"."group_name" = ?`, groupName).
		Relation("Memberships.User").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &group, nil
}

// UpdateGroupMembership adds/updates a user's membership to a group
func (db *Database) UpdateGroupMembership(username, groupName string) (*Membership, error) {
	group, err := db.GetGroup(groupName)
	if err != nil {
		return nil, err
	}

	user, err := db.GetUser(username)
	if err != nil {
		return nil, err
	}

	entry := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}

	// Perform an insert or update
	_, err = db.conn.Model(&entry).
		OnConflict("(group_id, user_id) DO NOTHING").
		Insert()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &entry, nil
}

// RemoveGroupMembership removes a user's membership to a group
func (db *Database) RemoveGroupMembership(username, groupName string) error {
	group, err := db.GetGroup(groupName)
	if err != nil {
		return err
	}

	user, err := db.GetUser(username)
	if err != nil {
		return err
	}

	entry := Membership{
		GroupId: group.Id,
		UserId:  user.Id,
	}

	_, err = db.conn.Model(&entry).
		Where(`group_id = ?group_id AND user_id = ?user_id`).
		Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
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

// ListUserGroupNames returns a list of groups the user is a member of
func (db *Database) ListUserGroupNames(username string) ([]string, error) {
	var groupNames []string

	user, err := db.GetUser(username)
	if err != nil {
		return nil, err
	}

	for _, membership := range user.Memberships {
		groupNames = append(groupNames, membership.Group.GroupName)
	}

	return groupNames, nil
}

// CreatePAT creates the personal access token in the database
func (db *Database) CreatePAT(username, description string) (*PersonalAccessToken, error) {
	owner, err := db.GetUser(username)
	if err != nil {
		return nil, err
	}

	token := PersonalAccessToken{
		OwnerId:     owner.Id,
		Description: description,
	}

	_, err = db.conn.Model(&token).Insert()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &token, nil
}

// GetPAT gets the personal access token from the database
func (db *Database) GetPAT(pat string) (*PersonalAccessToken, error) {
	token := PersonalAccessToken{}
	err := db.conn.Model(&token).
		Where(`"personal_access_token"."pat" = ?`, pat).
		Relation("Owner").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return &token, nil
}

// GetPAT gets the personal access token from the database
func (db *Database) ListPAT(username string) ([]*PersonalAccessToken, error) {
	var tokens []*PersonalAccessToken

	user, err := db.GetUser(username)
	if err != nil {
		return tokens, err
	}

	err = db.conn.Model(&tokens).
		Where(`"personal_access_token"."owner_id" = ?`, user.Id).
		Select()
	if err != nil {
		return tokens, ConvertError(err)
	}

	return tokens, nil
}

// DeletePAT deletes the personal access token from the database
func (db *Database) DeletePAT(patId int64, username string) error {
	user, err := db.GetUser(username)
	if err != nil {
		return ConvertError(err)
	}

	token := PersonalAccessToken{
		Id:      patId,
		OwnerId: user.Id,
	}

	// Verify token exists to delete and the user owns it
	err = db.conn.Model(&token).
		Where(`id = ?id AND owner_id = ?owner_id`).
		Select()
	if err != nil {
		return ConvertError(err)
	}

	_, err = db.conn.Model(&token).Where(`id = ?id AND owner_id = ?owner_id`).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}
