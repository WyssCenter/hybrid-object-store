package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-core/pkg/database"
)

var (
	ErrUnauthorized = errors.New("user is not authorized")
)

func HandleError(c *gin.Context, err error) {
	switch err {
	case database.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
	case database.ErrExists:
		c.JSON(http.StatusForbidden, gin.H{"error": "Resource already exists"})
	case database.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	case ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is unauthorized"})
	default:
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
	}
}
