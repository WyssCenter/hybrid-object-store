package database

import (
	"time"

	"github.com/pkg/errors"
)

type DatasetDeleteStatus string

const (
	NOT_SCHEDULED DatasetDeleteStatus = "NOT_SCHEDULED"
	SCHEDULED     DatasetDeleteStatus = "SCHEDULED"
	IN_PROGRESS   DatasetDeleteStatus = "IN_PROGRESS"
	ERROR         DatasetDeleteStatus = "ERROR"
)

// CreateDataset creates the Dataset in the database
func (db *Database) CreateDataset(namespace *Namespace, name, description, rootDir, username string) error {
	owner, err := db.GetOrCreateUser(username)
	if err != nil {
		return ConvertError(errors.Wrap(err, "failed to get or create user"))
	}

	entry := Dataset{
		NamespaceId:   namespace.Id,
		Name:          name,
		Description:   description,
		Created:       time.Now(),
		RootDirectory: rootDir,
		OwnerId:       owner.Id,
		SyncEnabled:   false,
		SyncType:      "simplex",
	}

	_, err = db.conn.Model(&entry).Insert()
	if err != nil {
		return ConvertError(errors.Wrap(err, "failed to create dataset"))
	}

	return nil
}

// GetDataset gets the Dataset from the database
func (db *Database) GetDataset(namespace *Namespace, name string) (*Dataset, error) {
	dataset := &Dataset{}
	err := db.conn.Model(dataset).
		Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, name).
		Relation("Owner.username").
		Relation("Permissions.Group").
		Relation("Namespace").
		Relation("Namespace.ObjectStore").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return dataset, nil
}

// SetDatasetSync sets the SyncEnabled flag for a dataset
// Note: triggers an update to the LastModified SyncConfigurationMeta timestamp if there is a change in the database
func (db *Database) SetDatasetSync(dataset *Dataset, syncEnabled bool, syncType, syncPolicy string) error {
	_, err := db.conn.Model(dataset).
		Set("sync_enabled = ?", syncEnabled).
		Set("sync_type = ?", syncType).
		Set("sync_policy = ?", syncPolicy).
		Where("id = ?id").Update()
	if err != nil {
		return ConvertError(err)
	}

	dataset.SyncEnabled = syncEnabled
	dataset.SyncType = syncType
	dataset.SyncPolicy = syncPolicy

	return nil
}

// DeleteDataset deletes a dataset from the database
// Note: there is no error if you delete a non-existent dataset
// Note: triggers an update to the LastModified SyncConfigurationMeta timestamp if there is a change in the database
func (db *Database) DeleteDataset(namespace *Namespace, name string) error {
	var dataset Dataset
	_, err := db.conn.Model(&dataset).
		Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, name).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// SetDatasetDeleteMarker sets the dataset delete_status and delete_on fields
func (db *Database) SetDatasetDeleteMarker(namespace *Namespace, datasetName string,
	status DatasetDeleteStatus, deleteDelayMinutes int) error {

	var dataset Dataset
	if status == SCHEDULED {
		// If we are scheduling a delete, you must also set the delete_on field based on the service configuration
		delete_on := time.Now().Add(time.Duration(deleteDelayMinutes) * time.Minute)
		_, err := db.conn.Model(&dataset).
			Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, datasetName).
			Set("delete_status = ?", status).
			Set("delete_on = ?", delete_on).Update()
		if err != nil {
			return ConvertError(err)
		}
	} else {
		_, err := db.conn.Model(&dataset).
			Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, datasetName).
			Set("delete_status = ?", status).Update()
		if err != nil {
			return ConvertError(err)
		}
	}

	return nil
}

// GetDatasetsByDeleteStatus returns a list of datasets that have `dataset_status`=`status DatasetDeleteStatus`
func (db *Database) GetDatasetsByDeleteStatus(status DatasetDeleteStatus) ([]*Dataset, error) {
	var datasets []*Dataset

	err := db.conn.Model(&datasets).
		Where(`"dataset"."delete_status" = ?`, status).
		Order("id ASC").
		Relation("Namespace").
		Relation("Namespace.ObjectStore").
		Relation("Permissions.Group").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return datasets, nil
}
