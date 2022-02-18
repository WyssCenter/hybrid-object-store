package database

import (
	"errors"
	"strings"
)

// GetObjectStore will get an object store from the database by its name
func (db *Database) GetObjectStore(name string) (*ObjectStore, error) {
	obj := &ObjectStore{}

	err := db.conn.Model(obj).
		Where(`"object_store"."name" = ?`, name).
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return obj, nil
}

// CreateObjectStore will create a objectstore in the database
func (db *Database) CreateObjectStore(name, description, objectStoreType, endpoint, region, profile, roleArn, notificationArn string) error {

	// Object store names are not allowed to contain pipe character to ensure
	// objects have unique id's when indexing their metadata
	if strings.Contains(name, "|") {
		return errors.New("object store names cannot contain `|` character")
	}

	obj := &ObjectStore{Name: name,
		Description:     description,
		Endpoint:        endpoint,
		ObjectStoreType: objectStoreType,
		Region:          region,
		Profile:         profile,
		RoleArn:         roleArn,
		NotificationArn: notificationArn}

	_, err := db.conn.Model(obj).Insert()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// DeleteObjectStore will delete an object store from the database by its name
func (db *Database) DeleteObjectStore(name string) error {
	obj := &ObjectStore{}

	_, err := db.conn.Model(obj).Where(`"object_store"."name" = ?`, name).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// ListObjectStores will list all object stores with limit/offset for pagination
func (db *Database) ListObjectStores(limit int, offset int) ([]*ObjectStore, error) {
	var objs []*ObjectStore

	err := db.conn.Model(&objs).Order("id ASC").Limit(limit).Offset(offset).Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return objs, nil
}
