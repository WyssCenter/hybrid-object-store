package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetObjectStore gets a object store by its name
// @Summary Get an object store based on its name
// @Schemes
// @Description Get information about an object store that is available for namespaces to use in this server.
// @Tags Object Store
// @Accept json
// @Produce json
// @Param        name   path      string  true  "Object Store Name"
// @Success 200 {object} database.ObjectStore
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Security BearerToken
// @Router /object_store/{name} [get]
func GetObjectStore(c *gin.Context) {
	_, db := getAppConfig(c)

	objectStoreName := c.Param("object_store")

	objStore, err := db.GetObjectStore(objectStoreName)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, objStore)
}

// ListObjectStores lists all available object stores
// @Summary List all available object stores
// @Schemes
// @Description Get information about all available object store in this server.
// @Tags Object Store
// @Accept json
// @Produce json
// @Success 200 {object} []database.ObjectStore
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Security BearerToken
// @Router /object_store/ [get]
func ListObjectStores(c *gin.Context) {
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

	objectStores, err := db.ListObjectStores(int(limitStr), int(offsetStr))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, objectStores)
}
