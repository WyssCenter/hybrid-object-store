package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-sync/pkg/config"

	service "github.com/gigantum/hoss-service"
)

func main() {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil && debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	}

	configuration := config.Load("")
	CheckForServices(configuration)

	ctx := context.TODO() // TODO create signal handler

	// Get the service JWT and start the refresh routine
	tokens := service.GetRenewingServiceJWT(configuration.AuthEndpoint, configuration.RefreshIntervals.AuthToken)
	go tokens.RefreshRoutine()

	// Start the UpdateMuxer for monitoring SyncConfiguration changes
	populatedConfigs := &PopulatedCoreServiceConfigurations{}
	go populatedConfigs.UpdateMuxer(ctx, configuration, tokens)

	// Start monitoring for bucket events
	Demuxer(ctx, configuration, tokens, populatedConfigs)
}

// CheckForServices verifies that the dependent services have started and are accepting connections
func CheckForServices(configuration *config.Configuration) {
	var err error
	// Auth Service check
	for i := 0; i < 12; i++ {
		response, err := http.Get(configuration.AuthEndpoint + "/ping")
		if err == nil {
			if response.StatusCode == 200 {
				break
			}
		}

		logrus.Info("Auth service is not ready, sleeping")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		logrus.WithField("error", err.Error()).Fatal("Couldn't connect to Auth service after 60 seconds")
	}

	// Core Service check - this will also ensure that opensearch is ready, as the core service waits for search
	for _, coreService := range configuration.CoreServices {
		// Hack to support running on localhost
		coreService = strings.Replace(coreService, "localhost/core", "core:8080", 1)

		for i := 0; i < 12; i++ {
			response, err := http.Get(coreService + "/discover")
			if err == nil {
				if response.StatusCode == 200 {
					break
				}
			}

			logrus.Info("Core service is not ready, sleeping")
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			logrus.WithField("error", err.Error()).Fatal("Couldn't connect to Core service after 60 seconds")
		}
	}
}
