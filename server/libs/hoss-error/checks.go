package errors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(cause error, message string) error {
	return New(message + ": " + cause.Error())
}

func IsNotExists(err error) bool {
	return err == ErrNotFound
}

func IsAlreadyExists(err error) bool {
	return err == ErrExists
}

func IsInvalidInput(err error) bool {
	return err == ErrInvalidInput
}

func IsUnauthorized(err error) bool {
	return err == ErrUnauthorized
}

// SetResponse sets the error message for the request
func SetResponse(c *gin.Context, err error) {
	switch err {
	case ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
	case ErrExists:
		c.JSON(http.StatusForbidden, gin.H{"error": "Resource already exists"})
	case ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	default:
		logrus.Infof("Unhandled error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unhandled error"})
	}
}

func ConvertError(err error) error {
	// pg errors
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key value") {
		return ErrExists
	} else if strings.Contains(errMsg, "no rows in result set") {
		return ErrNotFound
	} else if strings.Contains(errMsg, "invalid input value") {
		return ErrInvalidInput
	}

	return err
}
