package api

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	JWT       string
	Claims    map[string]interface{}
	Username  string
	Role      string
	Groups    []string
	IsService bool
}

// getAppConfig gets the Configuration and Database references for API handler methods to use
func getAppConfig(c *gin.Context) (*config.Configuration, *database.Database) {
	val, ok := c.Get("config")
	if !ok {
		log.Fatal("Failed to load config in context")
	}

	coreConfig := val.(*config.Configuration)

	val, ok = c.Get("db")
	if !ok {
		log.Fatal("Failed to load database in context")
	}

	db := val.(*database.Database)

	return coreConfig, db

}

// getUserInfo gets the UserInfo of the currently authenticated user
func getUserInfo(c *gin.Context) UserInfo {
	val, ok := c.Get("user")
	if ok {
		userInfo := val.(UserInfo)
		return userInfo
	} else {
		return UserInfo{}
	}
}

const (
	ROLE_ADMIN      = "admin"
	ROLE_PRIVILEGED = "privileged"
	ROLE_USER       = "user"
	ROLE_SERVICE    = "service"
)

// validatePrivileged checks if the user has a privileged or admin role
func validatePrivileged(role string) bool {
	return role == ROLE_PRIVILEGED || role == ROLE_ADMIN || role == ROLE_SERVICE
}

// validateAdmin checks if the user has an admin role
func validateAdmin(role string) bool {
	return role == ROLE_ADMIN || role == ROLE_SERVICE
}

// getStores gets all of the ObjectStore references for API handler methods to use
func getStores(c *gin.Context) map[string]store.ObjectStore {
	val, ok := c.Get("stores")
	if !ok {
		return nil
	}

	stores := val.(map[string]store.ObjectStore)

	return stores
}

// getStoreByName returns an already loaded object store by its name from the context
func getStoreByName(stores map[string]store.ObjectStore, name string) (store.ObjectStore, error) {

	obj, ok := stores[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("The store '%s' was not found", name))
	}

	return obj, nil
}

// Get the core service endpoint for this server
func getCoreServiceEndpoint() string {
	return os.Getenv("EXTERNAL_HOSTNAME") + "/core/v1"
}
