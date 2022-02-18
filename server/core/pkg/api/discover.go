package api

import (
	"os"
	"strings"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gin-gonic/gin"
)

// Discover provides information about the current server
// @Summary Fetch server information for client configuration
// @Schemes
// @Tags Discover
// @Description This open endpoint returns information about the server's configuration and provides a health check
// @Accept json
// @Produce json
// @Success 200 {object} object{alive=bool,dev=bool,version=string,build=string,services=[]string,auth_service=string}
// @Router /discover [get]
func Discover(c *gin.Context) {
	val, ok := c.Get("config")
	alive := "false"
	dev := false
	auth_service := ""
	delete_delay_minutes := 0
	if ok {
		alive = "true"
		coreConfig := val.(*config.Configuration)
		dev = coreConfig.Server.Dev
		auth_service = coreConfig.Server.AuthService
		delete_delay_minutes = coreConfig.Server.DatasetDeleteDelayMinutes

		// When the auth service is running locally, this will be an inteternal route 'auth:8080',
		// which we don't want to send to clients. Instead set to EXTERNAL_HOSTNAME.
		if strings.Contains(auth_service, "auth:8080") {
			auth_service = os.Getenv("EXTERNAL_HOSTNAME") + "/auth/v1"
		}
	}

	build_hash, err := getBuildHash()
	if err != nil {
		HandleError(c, err)
		return
	}

	version, err := getVersion()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"alive":                alive,
		"dev":                  dev,
		"version":              version,
		"build":                build_hash,
		"services":             getAvailableServices(),
		"auth_service":         auth_service,
		"delete_delay_minutes": delete_delay_minutes,
	})
}
