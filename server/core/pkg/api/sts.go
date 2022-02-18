package api

import (
	"net/http"

	"github.com/gigantum/hoss-core/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetServiceSTSCredentials renders the service's policy and get STS credentials for all data in an object store
// @Summary Get Service Account STS credentials
// @Schemes
// @Description Re-renders the service account's policy and gets STS credentials for the specified object store
// @Description **NOTE: This endpoint is only available to the service account**
// @Tags Service Account
// @Accept json
// @Produce json
// @Param	objectStoreName   path      string  true  "Object Store Name"
// @Success 200 {object} store.Credentials
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /object_store/{objectStoreName}/sts [get]
func GetServiceSTSCredentials(c *gin.Context) {
	_, db := getAppConfig(c)
	user := getUserInfo(c)

	if !user.IsService {
		HandleError(c, ErrUnauthorized)
	}

	objectStoreName := c.Param("object_store")

	// Load the Object Store from the request context
	currentStore, err := getStoreByName(getStores(c), objectStoreName)
	if err != nil {
		HandleError(c, err)
		return
	}

	// compile permissions for all datasets in each namespace in this object store
	var perms []*database.Permission
	limit := 20
	offset := 0
	for {
		namespaces, err := db.ListNamespaces(limit, offset)
		if err != nil {
			HandleError(c, err)
			return
		}
		if len(namespaces) == 0 {
			break
		}
		for _, ns := range namespaces {
			// only add permissions for namespaces in current object store
			if ns.ObjectStore.Name == currentStore.GetName() {
				perms = append(perms, &database.Permission{
					Dataset: &database.Dataset{
						Namespace: ns,
						Name:      "*", // Grant access to all datasets within the namespace
						// Which gives proactive permissions in case a new dataset
						// is created after STS credentials are generated
					},
					Permission: database.PERM_READ_WRITE,
				})
			}
		}
		offset += limit
	}

	err = currentStore.SetUserPolicy(user.Username, perms)
	if err != nil {
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
		return
	}

	// Get Credentials from the object store
	creds, err := currentStore.GetSTSCredentials(user.JWT, user.Claims, user.Username)
	if err != nil {
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
		return
	}

	c.JSON(http.StatusOK, creds)
}

// GetUserSTSCredentials renders a user's policy and gets STS credentials for the object store
// @Summary Get STS credentials
// @Schemes
// @Description Re-renders a user's policy and gets STS credentials for the object store related to the specified namespace
// @Tags Credentials
// @Accept json
// @Produce json
// @Param	namespaceName   path      string  true  "Namespace Name"
// @Success 200 {object} store.Credentials
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /namespace/{namespaceName}/sts [get]
func GetUserSTSCredentials(c *gin.Context) {
	_, db := getAppConfig(c)
	user := getUserInfo(c)

	if user.IsService {
		HandleError(c, ErrUnauthorized)
	}

	namespaceName := c.Param("namespace")

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

	// Render the user's policy in case there have been changes
	perms, err := db.GetPermissionsByUser(&namespace.ObjectStore, user.Username, false)
	if err != nil {
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
		return
	}

	err = currentStore.SetUserPolicy(user.Username, perms)
	if err != nil {
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
		return
	}

	// Get Credentials from the object store
	creds, err := currentStore.GetSTSCredentials(user.JWT, user.Claims, user.Username)
	if err != nil {
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
		return
	}

	c.JSON(http.StatusOK, creds)
}
