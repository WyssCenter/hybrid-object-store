package api

import (
	"net/http"
	"strconv"

	"github.com/gigantum/hoss-core/pkg/sync"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type namespaceInput struct {
	// Name is the unique name of the namespace that is used to identify it in the API
	Name string `json:"name" binding:"required"`
	// Description is a short description of the namespace to create
	Description string `json:"description" binding:"required"`
	// ObjectStoreName is the name of the object store that will back this namespace
	ObjectStoreName string `json:"object_store_name" binding:"required"`
	// BucketName is the name of the bucket that will store this namespace's datasets
	BucketName string `json:"bucket_name" binding:"required"`
}

// CreateNamespace creates a new namespace
// @Summary Create a namespace
// @Schemes
// @Description Create a new namespace. The authorized user must have the admin role.
// @Tags Namespace
// @Accept json
// @Produce json
// @Param	namespaceInput		body	api.namespaceInput	true	"Namespace Input"
// @Success 200 {object} database.Namespace
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{name} [post]
func CreateNamespace(c *gin.Context) {
	_, db := getAppConfig(c)

	userInfo := getUserInfo(c)
	if isAdmin := validateAdmin(userInfo.Role); !isAdmin {
		HandleError(c, ErrUnauthorized)
		return
	}

	nsInput := namespaceInput{}
	err := c.Bind(&nsInput)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = db.CreateNamespace(nsInput.Name, nsInput.Description, nsInput.ObjectStoreName, nsInput.BucketName)
	if err != nil {
		HandleError(c, err)
		return
	}

	ns, err := db.GetNamespace(nsInput.Name)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Notify the sync service so service account STS credentials can be refreshed
	// to include access to datasets in this namespace immediately. For simplicity, only warn on errors because if
	// the notification doesn't work, sync service will eventually auto-refresh.
	currentStore, err := getStoreByName(getStores(c), ns.ObjectStore.Name)
	if err != nil {
		msg := "Failed to load the object store while preparing Namespace Create API Sync Notification." +
			"Sync service will not function with the new namespace until credentials timeout and are automatically refreshed."
		logrus.Warnf(msg)
	} else {
		err = sync.CreateNamespaceHandler(c, currentStore, ns)
		if err != nil {
			msg := "Failed to send Namespace Create API Sync Notification." +
				"Sync service will not function with the new namespace until credentials timeout and are automatically refreshed."
			logrus.Warnf(msg)
		}
	}

	c.JSON(http.StatusCreated, ns)
}

// GetNamespace gets a namespace by its name
// @Summary Get a namespace
// @Schemes
// @Description Get a namespace based on its name
// @Tags Namespace
// @Accept json
// @Produce json
// @Param        name   path      string  true  "Namespace Name"
// @Success 200 {object} database.Namespace
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{name} [get]
func GetNamespace(c *gin.Context) {
	_, db := getAppConfig(c)

	namespaceName := c.Param("namespace")

	ns, err := db.GetNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ns)
}

// DeleteNamespace deletes an existing namespace
// @Summary Delete a namespace
// @Schemes
// @Description Delete a namespace. Note, the namespace must be empty before deleting (i.e. no datasets remain) and syncing should be disabled.
// @Description The authorized user must have the admin role.
// @Tags Namespace
// @Accept json
// @Produce json
// @Param        name   path      string  true  "Namespace Name"
// @Success 204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Security BearerToken
// @Router /namespace/{name} [delete]
func DeleteNamespace(c *gin.Context) {
	_, db := getAppConfig(c)

	namespaceName := c.Param("namespace")

	userInfo := getUserInfo(c)
	if isAdmin := validateAdmin(userInfo.Role); !isAdmin {
		HandleError(c, ErrUnauthorized)
		return
	}

	err := db.DeleteNamespace(namespaceName)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListNamespaces lists all available namespaces
// @Summary List all available namespaces
// @Schemes
// @Description Get information about all available namespaces in this server.
// @Tags Namespace
// @Accept json
// @Produce json
// @Success 200 {object} []database.Namespace
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Security BearerToken
// @Router /namespace/ [get]
func ListNamespaces(c *gin.Context) {
	_, db := getAppConfig(c)

	limitStr, err := strconv.Atoi(c.DefaultQuery("limit", "25"))
	if err != nil {
		HandleError(c, err)
		return
	}

	offsetStr, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		HandleError(c, err)
		return
	}

	namespaces, err := db.ListNamespaces(int(limitStr), int(offsetStr))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, namespaces)
}
