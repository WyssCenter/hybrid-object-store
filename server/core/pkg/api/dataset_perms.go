package api

import (
	"net/http"

	"github.com/gigantum/hoss-core/pkg/sync"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UpdateUserDatasetPerms determines the user's default groupname and then updates the permissions for the default group
// @Summary Update the permissions for a user
// @Schemes
// @Description Update the permissions to a dataset for a user, by modifying the permissions on the user's default group
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName   path      string  true  "Dataset Name"
// @Param	username   path      string  true  "Username to modify"
// @Param	accessLevel   path      string  true  "Access Level ('r' or 'rw')"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName}/user/{username}/access/{accessLevel} [put]
func UpdateUserDatasetPerms(c *gin.Context) {
	_, db := getAppConfig(c)

	namespace := c.Param("namespace")
	datasetName := c.Param("name")
	username := c.Param("username")
	accessLevel := c.Param("accesslevel")

	updateDatasetPerms(c, namespace, datasetName, db.GetUserDefaultGroup(username), accessLevel)
}

// UpdateGroupDatasetPerms updates the permissions for a group on the dataset
// @Summary Update the permissions for a group
// @Schemes
// @Description Update the permissions to a dataset for a group. This will re-render all of the policies for users in the group.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName   path      string  true  "Dataset Name"
// @Param	groupName   path      string  true  "Name of the group to modify"
// @Param	accessLevel   path      string  true  "Access Level ('r' or 'rw')"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName}/group/{groupName}/access/{accessLevel} [put]
func UpdateGroupDatasetPerms(c *gin.Context) {

	namespace := c.Param("namespace")
	datasetName := c.Param("name")
	groupName := c.Param("groupname")
	accessLevel := c.Param("accesslevel")

	// You can only add the public group with read-only access
	if groupName == "public" && accessLevel == "rw" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The public group can only be added with read-only permissions."})
		return
	}

	updateDatasetPerms(c, namespace, datasetName, groupName, accessLevel)
}

// updateDatasetPerms updates permissions on a dataset
func updateDatasetPerms(c *gin.Context, namespaceName, datasetName, groupName, accessLevel string) {
	_, db := getAppConfig(c)
	userInfo := getUserInfo(c)

	// Load the specified namespace
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

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		HandleError(c, ErrUnauthorized)
		return
	}

	currentPerm, err := db.GetDatasetPermissions(namespace, datasetName, groupName)
	if err != nil {
		HandleError(c, err)
		return
	}

	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		HandleError(c, err)
		return
	}

	tx := &Transaction{}

	tx.AddFunction(func() error {
		return db.UpdateDatasetPermissions(namespace, datasetName, groupName, accessLevel)
	})

	// user re-render rollbacks need to happen after the dataset permissions have been
	// rolled back, and rollbacks occur in reverse order
	for _, membership := range group.Memberships {
		tx.AddRollback(func() error {
			perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, membership.User.Username, false)
			if err != nil {
				return err
			}
			err = currentStore.SetUserPolicy(membership.User.Username, perms)
			if err != nil {
				return err
			}
			return nil
		})
	}
	tx.AddRollback(func() error {
		if currentPerm == "" {
			return db.RemoveDatasetPermissions(namespace, datasetName, groupName)
		} else {
			return db.UpdateDatasetPermissions(namespace, datasetName, groupName, currentPerm)
		}
	})

	// for each user in the group, update their policy with new group permissions
	for _, membership := range group.Memberships {
		tx.AddFunction(func() error {
			perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, membership.User.Username, false)
			if err != nil {
				return err
			}
			err = currentStore.SetUserPolicy(membership.User.Username, perms)
			if err != nil {
				return err
			}
			return nil
		})
	}

	// If request is from a service account, don't sync since this request is the result of a sync operation
	if !userInfo.IsService {
		tx.AddFunction(func() error {
			dataset, err := db.GetDataset(namespace, datasetName)
			if err != nil {
				return err
			}

			return sync.SyncPermissionsHandler(c, currentStore, namespace, dataset, groupName, accessLevel)
		})
	} else {
		logrus.Info("Ignoring Service Account sync operation.")
	}

	err = tx.Execute()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveUserDatasetPerms determines the user's default groupname and then removes the permissions for the default group
// @Summary Remove access to a dataset for a user
// @Schemes
// @Description Remove the permissions to a dataset for a user, revoking access.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName   path      string  true  "Dataset Name"
// @Param	username   path      string  true  "Name of the user to modify"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName}/user/{username} [delete]
func RemoveUserDatasetPerms(c *gin.Context) {
	_, db := getAppConfig(c)

	namespace := c.Param("namespace")
	datasetName := c.Param("name")
	username := c.Param("username")

	removeDatasetPerms(c, namespace, datasetName, db.GetUserDefaultGroup(username))
}

// RemoveGroupDatasetPerms removes the permissions for a group on the dataset
// @Summary Remove access to a dataset for a group
// @Schemes
// @Description Remove the permissions to a dataset for a group. This will re-render all of the policies for users in the group.
// @Tags Dataset
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Param	datasetName   path      string  true  "Dataset Name"
// @Param	groupName   path      string  true  "Name of the group to modify"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/dataset/{datasetName}/group/{groupName} [delete]
func RemoveGroupDatasetPerms(c *gin.Context) {

	namespace := c.Param("namespace")
	datasetName := c.Param("name")
	groupName := c.Param("groupname")

	// You cannot remove the admin group
	if groupName == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The admin group can not be removed."})
		return
	}

	removeDatasetPerms(c, namespace, datasetName, groupName)
}

// removeDatasetPerms removes permissions on a dataset
func removeDatasetPerms(c *gin.Context, namespaceName, datasetName, groupName string) {
	_, db := getAppConfig(c)
	userInfo := getUserInfo(c)

	// Load the specified namespace
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

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		HandleError(c, ErrUnauthorized)
		return
	}

	currentPerm, err := db.GetDatasetPermissions(namespace, datasetName, groupName)
	if err != nil {
		HandleError(c, err)
		return
	}

	group, err := db.GetOrCreateGroup(groupName)
	if err != nil {
		HandleError(c, err)
		return
	}

	tx := &Transaction{}

	tx.AddFunction(func() error {
		return db.RemoveDatasetPermissions(namespace, datasetName, groupName)
	})

	// user re-render rollbacks need to happen after the dataset permissions have been
	// rolled back, and rollbacks occur in reverse order
	for _, membership := range group.Memberships {
		tx.AddFunction(func() error {
			perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, membership.User.Username, false)
			if err != nil {
				return err
			}
			err = currentStore.SetUserPolicy(membership.User.Username, perms)
			if err != nil {
				return err
			}
			return nil
		})
	}
	tx.AddRollback(func() error {
		if currentPerm != "" {
			// handle if they already didn't have permissions
			return db.UpdateDatasetPermissions(namespace, datasetName, groupName, currentPerm)
		}

		return nil
	})

	// for each user in the group, update their policy with new group permissions
	for _, membership := range group.Memberships {
		tx.AddFunction(func() error {
			perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, membership.User.Username, false)
			if err != nil {
				return err
			}
			err = currentStore.SetUserPolicy(membership.User.Username, perms)
			if err != nil {
				return err
			}
			return nil
		})
	}

	// If request is from a service account, don't sync since this request is the result of a sync operation
	if !userInfo.IsService {
		tx.AddFunction(func() error {
			dataset, err := db.GetDataset(namespace, datasetName)
			if err != nil {
				return err
			}

			return sync.SyncPermissionsHandler(c, currentStore, namespace, dataset, groupName, "")
		})
	} else {
		logrus.Info("Ignoring Service Account sync operation.")
	}

	err = tx.Execute()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
