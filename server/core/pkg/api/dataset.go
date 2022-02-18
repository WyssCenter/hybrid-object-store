package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-core/pkg/database"
)

type datasetInput struct {
	// Name is the unique name of the dataset that is used to identify it in the API
	Name string `json:"name" binding:"required"`
	// Description is a short description of the dataset to create
	Description *string `json:"description" binding:"required"`
}

// CreateDataset creates a new dataset by creating a root folder in the object store and
// updating the user's policy
// @Summary Create a dataset
// @Schemes
// @Description Create a new dataset in the specified namespace and update the user's access policy
// @Description The authorized user must have the admin or privileged role.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	datasetInput		body	api.datasetInput	true	"Dataset Input"
// @Param        namespaceName   path      string  true  "Namespace Name"
// @Success 200 {object} database.Dataset
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/ [post]
func CreateDataset(c *gin.Context) {
	_, db := getAppConfig(c)

	dsInput := datasetInput{}
	err := c.Bind(&dsInput)
	if err != nil {
		HandleError(c, err)
		return
	}

	namespaceName := c.Param("namespace")
	userInfo := getUserInfo(c)
	username := userInfo.Username

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		HandleError(c, ErrUnauthorized)
		return
	}

	// Dataset names not allowed to have a slash as this can cause collisions
	// ex: file subdir/file.txt in dataset dir has path /dir/subdir/file.txt
	//	   file file.txt in dataset dir/subdir has path /dir/subdir/file.txt
	if strings.Contains(dsInput.Name, "/") {
		HandleError(c, errors.New("dataset names cannot contain `/` character"))
	}

	// Dataset names are also not allowed to contain pipe character to ensure
	// objects have unique id's when indexing their metadata
	if strings.Contains(dsInput.Name, "|") {
		HandleError(c, errors.New("dataset names cannot contain `|` character"))
	}

	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Load the Object Store from the request context
	currentStore, err := getStoreByName(getStores(c), namespace.ObjectStore.Name)
	if err != nil {
		HandleError(c, err)
		return
	}

	tx := &Transaction{}

	// DP NOTE: The service account is used by the sync service to create datasets when needed
	//          This code skips adding the service account as a permitted user and rendering their
	//          policy, as the GetSTSCredentials function handles generating special credentials
	//          for the service account. The rest of the API has special checks to allow the service
	//          account to make needed API calls.

	// Create database entry
	tx.AddFunction(func() error {
		// Note, we set the RootDirectory value to the dataset name followed by a '/'. This value is used
		// throughout the system as the prefix where objects are written for the dataset. In the sync
		// service, in the api message handleSync() function, we manually craft this value again since
		// RootDirectory for the dataset is not available there. If this default behavior is ever changed
		// the handleSync() function must also be modified.
		return db.CreateDataset(namespace, dsInput.Name, *dsInput.Description, dsInput.Name+"/", username)
	})
	tx.AddRollback(func() error {
		return db.DeleteDataset(namespace, dsInput.Name)
	})

	// Add initial permissions
	if !userInfo.IsService {
		// Grant the user access via their default group (group they are the sole member of)
		tx.AddFunction(func() error {
			return db.UpdateDatasetPermissions(namespace, dsInput.Name, db.GetUserDefaultGroup(username), database.PERM_READ_WRITE)
		})
		// Add the automated "admin" group to ensure admins can see and mutate all resources
		tx.AddFunction(func() error {
			return db.UpdateDatasetPermissions(namespace, dsInput.Name, "admin", database.PERM_READ_WRITE)
		})
	}
	// NOTE: no rollback as deleting the dataset will remove the permissions

	// Create datastore entry
	tx.AddFunction(func() error {
		return currentStore.CreateDataset(dsInput.Name, namespace)
	})
	tx.AddRollback(func() error {
		// Get the dataset's root directory
		ds, err := db.GetDataset(namespace, dsInput.Name)
		if err != nil {
			return err
		}
		return currentStore.DeleteDataset(ds.RootDirectory, namespace)
	})

	// Enable bucket notifications for the given dataset
	tx.AddFunction(func() error {
		ds, err := db.GetDataset(namespace, dsInput.Name)
		if err != nil {
			return err
		}
		return currentStore.EnableEvents(namespace, ds)
	})
	tx.AddRollback(func() error {
		ds, err := db.GetDataset(namespace, dsInput.Name)
		if err != nil {
			return err
		}
		return currentStore.DisableEvents(namespace, ds)
	})

	// Render user's policy
	if !userInfo.IsService {
		tx.AddFunction(func() error {
			perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, username, false)
			if err != nil {
				return err
			}

			return currentStore.SetUserPolicy(username, perms)
		})
	}
	// NOTE: no rollback as this is the final step

	err = tx.Execute()
	if err != nil {
		HandleError(c, err)
		return
	}

	// Load dataset you just created
	ds, err := db.GetDataset(namespace, dsInput.Name)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, ds)
}

// DeleteDataset deletes marks a dataset for deletion and removes user permissions to the underlying data
// @Summary Delete a dataset
// @Schemes
// @Description Delete an existing dataset in the specified namespace and update access policies
// @Description for users who had access to the dataset.
// @Description The authorized user must have the admin or privileged role.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName  	 path      string  true  "Dataset Name"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName} [delete]
func DeleteDataset(c *gin.Context) {
	config, db := getAppConfig(c)

	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	userInfo := getUserInfo(c)

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		HandleError(c, ErrUnauthorized)
		return
	}

	// Load the specified namespace
	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Get all users that have access to this dataset so we can re-render their policies
	// for the object store containing the dataset. They will lose access to the dataset
	// while it is in the "SCHEDULED" state
	usernames, err := db.GetUsersWithPermissionsToDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Mark for delete. This will also remove this dataset from user's permissions
	err = db.SetDatasetDeleteMarker(namespace, datasetName, database.SCHEDULED, config.Server.DatasetDeleteDelayMinutes)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Load the Object Store from the request context
	currentStore, err := getStoreByName(getStores(c), namespace.ObjectStore.Name)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Rerender policies for affected users
	for _, u := range usernames {
		userPerms, err := db.GetPermissionsByUser(&namespace.ObjectStore, u, false)
		if err != nil {
			logrus.Warnf("DATASET DELETE FAILURE: Could not get permissions for user %s while rendering policy", u)
		}

		err = currentStore.SetUserPolicy(u, userPerms)
		if err != nil {
			logrus.Warnf("DATASET DELETE FAILURE: Could not rerender policy for user %s", u)
		}
	}

	c.Status(http.StatusNoContent)
}

// ListDataset lists all of the datasets the user has permissions on in the specified namespace
// @Summary List datasets
// @Schemes
// @Description Lists all datasets to which the user has access (via permissions) in the specified namespace
// @Tags Dataset
// @Accept json
// @Produce json
// @Param        namespaceName   path      string  true  "Namespace Name"
// @Success 200 {object} []database.Dataset
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/ [get]
func ListDataset(c *gin.Context) {
	_, db := getAppConfig(c)

	username := getUserInfo(c).Username
	namespaceName := c.Param("namespace")

	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, username, true)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Deduplicate datasets that a user may have access to through multiple groups
	// A namespace/dataset combo is unique
	datasets := []*database.Dataset{}
	var dedupList []string
	for _, perm := range perms {
		if perm.Dataset.Namespace.Name == namespaceName {
			skip := false
			for _, ds := range dedupList {
				dsId := fmt.Sprintf("%s/%s", perm.Dataset.Namespace.Name, perm.Dataset.Name)
				if dsId == ds {
					skip = true
					break
				}
			}
			if !skip {
				dedupList = append(dedupList,
					fmt.Sprintf("%s/%s", perm.Dataset.Namespace.Name, perm.Dataset.Name))
				datasets = append(datasets, perm.Dataset)
			}

		}
	}

	c.JSON(http.StatusOK, datasets)
}

// GetDataset gets information about a given dataset
// @Summary Get dataset
// @Schemes
// @Description Get a dataset by name. The user must at least have read access to load the dataset.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName   path      string  true  "Dataset Name"
// @Success 200 {object} database.Dataset
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName} [get]
func GetDataset(c *gin.Context) {
	_, db := getAppConfig(c)

	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	u := getUserInfo(c)
	// username := u.Username

	// Load the specified namespace
	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	dataset, err := db.GetDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// for _, perm := range dataset.Permissions {
	// 	if perm.Group.GroupName == db.GetUserDefaultGroup(username) {
	// 		c.JSON(http.StatusOK, dataset)
	// 		return
	// 	}
	// }

	// If the user is IN a group that has access to the dataset
	for _, perm := range dataset.Permissions {
		for _, group := range u.Groups {
			if perm.Group.GroupName == group {
				c.JSON(http.StatusOK, dataset)
				return
			}
		}
	}

	HandleError(c, database.ErrNotPermitted)
}

func RestoreDataset(c *gin.Context) {
	// RestoreDataset un-marks a dataset for deletion and re-renders user permissions to the underlying data
	config, db := getAppConfig(c)

	namespaceName := c.Param("namespace")
	datasetName := c.Param("name")
	userInfo := getUserInfo(c)

	// Only Admins can restore a dataset delete operation
	if isAdmin := validateAdmin(userInfo.Role); !isAdmin {
		HandleError(c, ErrUnauthorized)
		return
	}

	// Load the specified namespace
	namespace, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Get all users that have access to this dataset so we can re-render their policies
	// for the object store containing the dataset.
	usernames, err := db.GetUsersWithPermissionsToDataset(namespace, datasetName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Reset delete marker
	err = db.SetDatasetDeleteMarker(namespace, datasetName, database.NOT_SCHEDULED, config.Server.DatasetDeleteDelayMinutes)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Load the Object Store from the request context
	currentStore, err := getStoreByName(getStores(c), namespace.ObjectStore.Name)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Rerender policies for affected users
	for _, u := range usernames {
		userPerms, err := db.GetPermissionsByUser(&namespace.ObjectStore, u, false)
		if err != nil {
			logrus.Warnf("DATASET RESTORE FAILURE: Could not get permissions for user %s while rendering policy", u)
		}

		err = currentStore.SetUserPolicy(u, userPerms)
		if err != nil {
			logrus.Warnf("DATASET RESTORE FAILURE: Could not rerender policy for user %s", u)
		}
	}

	c.Status(http.StatusNoContent)
}
