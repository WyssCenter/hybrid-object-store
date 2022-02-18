package database

// GetPermissionsByUser gets the permissions for the given user in the specified object store
// if includeDeletedDatasets is true, datasets that have been marked for delete will be included
func (db *Database) GetPermissionsByUser(objStore *ObjectStore, username string, includeDeletedDatasets bool) ([]*Permission, error) {
	user, err := db.GetOrCreateUser(username)
	if err != nil {
		return nil, err
	}

	var perms []*Permission
	for _, membership := range user.Memberships {
		groupPerms, err := db.GetPermissionsByGroup(objStore, membership.Group.GroupName, includeDeletedDatasets)
		if err != nil {
			err = ConvertError(err)
			if err == ErrNotFound {
				continue
			} else {
				return nil, err
			}
		}
		perms = append(perms, groupPerms...)
	}

	return perms, nil
}

// GetPermissionsByGroup gets the permissions for the given group in the specified object store
// if includeDeletedDatasets is true, datasets that have been marked for delete will be included
func (db *Database) GetPermissionsByGroup(objStore *ObjectStore, groupName string, includeDeletedDatasets bool) ([]*Permission, error) {
	var group Group

	err := db.conn.Model(&group).
		Where(`"group"."group_name" = ?`, groupName).
		Relation("Permissions.Dataset.Namespace").
		Relation("Permissions.Dataset.Owner").
		Relation("Permissions.Dataset.Permissions.Group").
		Select()
	if err != nil {
		err = ConvertError(err)
		if err == ErrNotFound {
			return []*Permission{}, nil
		} else {
			return nil, err
		}
	}
	var perms []*Permission
	for _, perm := range group.Permissions {
		if !includeDeletedDatasets {
			// Only include datasets in the object store provided and are "NOT_SCHEDULED" for delete. If the
			// delete_state is "SCHEDULED", "ERROR", or "IN_PROGRESS" users should not be interacting
			// with the underlying storage.
			if perm.Dataset.Namespace.ObjectStoreId == objStore.Id && perm.Dataset.DeleteStatus == "NOT_SCHEDULED" {
				perms = append(perms, perm)
			}
		} else {
			// Include any dataset in the object store provided. This path should ONLY be used for listing all
			// available datasets to a user. Policies and permissions should typically not include datasets
			// that are marked for delete
			if perm.Dataset.Namespace.ObjectStoreId == objStore.Id {
				perms = append(perms, perm)
			}
		}
	}

	return perms, nil
}

// GetUsersWithPermissionsToDataset gets all of the users who have access to a dataset
func (db *Database) GetUsersWithPermissionsToDataset(namespace *Namespace, datasetName string) ([]string, error) {
	dataset := &Dataset{}
	err := db.conn.Model(dataset).
		Where(`"dataset"."namespace_id" = ? AND "dataset"."name" = ?`, namespace.Id, datasetName).
		Relation("Permissions.Group.Memberships.User").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	// Use a map to implement a set-like data structure
	usernameSet := make(map[string]bool)
	for _, p := range dataset.Permissions {
		for _, m := range p.Group.Memberships {
			usernameSet[m.User.Username] = true
		}
	}

	usernameList := make([]string, len(usernameSet))
	i := 0
	for k := range usernameSet {
		usernameList[i] = k
		i++
	}

	return usernameList, nil
}

// UpdateDatasetPermissions adds/updates the permissions the group has on the dataset
func (db *Database) UpdateDatasetPermissions(namespace *Namespace, datasetName, groupName, permission string) error {
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		return err
	}

	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		return err
	}

	entry := Permission{
		GroupId:    group.Id,
		DatasetId:  dataset.Id,
		Permission: permission,
	}

	// Perform an insert or update
	_, err = db.conn.Model(&entry).
		OnConflict("(group_id, dataset_id) DO UPDATE").
		Set("permission = ?permission").
		Insert()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// RemoveDatasetPermissions removes any permissions the group has on the dataset
// Note: there is no error if you remove a non-existent set of permissions
func (db *Database) RemoveDatasetPermissions(namespace *Namespace, datasetName, groupName string) error {
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		return err
	}

	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		return err
	}

	entry := Permission{
		GroupId:   group.Id,
		DatasetId: dataset.Id,
	}

	_, err = db.conn.Model(&entry).Where(`group_id = ?group_id AND dataset_id = ?dataset_id`).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// GetDatasetPermissions returns the permission the group has on the dataset
func (db *Database) GetDatasetPermissions(namespace *Namespace, datasetName, groupName string) (string, error) {
	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		return "", err
	}

	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		return "", err
	}

	entry := Permission{}

	err = db.conn.Model(&entry).Where(`group_id = ? AND dataset_id = ?`, group.Id, dataset.Id).Select()
	if err != nil {
		err = ConvertError(err)
		if err == ErrNotFound {
			return "", nil
		} else {
			return "", err
		}
	}

	return entry.Permission, nil
}
