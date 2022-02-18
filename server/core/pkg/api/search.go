package api

import (
	"net/http"

	"github.com/gigantum/hoss-core/pkg/opensearch"
	"github.com/gin-gonic/gin"
)

// CreateOrUpdateMetadataDocument creates or updates a document in the metadata search index
// @Summary Create or update a document in the metadata search index
// @Schemes
// @Description Creates or updates a document in the metadata search index
// @Description **NOTE: This endpoint is only available to the service account**
// @Tags Service Account
// @Accept json
// @Produce json
// @Param	documentInput		body	opensearch.MetadataIndexPayload	true	"Document Input"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search/document/metadata [put]
func CreateOrUpdateMetadataDocument(c *gin.Context) {
	config, _ := getAppConfig(c)
	user := getUserInfo(c)

	if !user.IsService {
		HandleError(c, ErrUnauthorized)
	}

	payload := opensearch.MetadataIndexPayload{}
	err := c.Bind(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = opensearch.CreateOrUpdateDocument(config.Server.ElasticsearchEndpoint, &payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteMetadataDocument deletes a document in the metadata search index
// @Summary Delete a document in the metadata search index
// @Schemes
// @Description Deletes a document in the metadata search index
// @Description **NOTE: This endpoint is only available to the service account**
// @Tags Service Account
// @Accept json
// @Produce json
// @Param	documentInput		body	opensearch.MetadataIndexPayload	true	"Document Input"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /search/document/metadata [delete]
func DeleteMetadataDocument(c *gin.Context) {
	config, _ := getAppConfig(c)
	user := getUserInfo(c)

	if !user.IsService {
		HandleError(c, ErrUnauthorized)
	}

	payload := opensearch.MetadataIndexPayload{}
	err := c.Bind(&payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = opensearch.DeleteDocument(config.Server.ElasticsearchEndpoint, &payload)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
