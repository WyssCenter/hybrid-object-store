package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-sync/pkg/config"
	"github.com/gigantum/hoss-sync/pkg/queue"

	errors "github.com/gigantum/hoss-error"
	service "github.com/gigantum/hoss-service"
)

// Demuxer loads the notification queues defined in the configuration file and routes the
// messages from them to the appropriate worker queues for execution. The Demuxer is responsible
// for creating / deleting the workers using the WorkerManager interface
func Demuxer(ctx context.Context, configuration *config.Configuration, tokens service.RenewingTokens, populatedConfigs *PopulatedCoreServiceConfigurations) {
	// Load the different notification queues
	notifications := make(chan config.Message)

	for _, coreService := range configuration.CoreServices {
		queues, err := QueryQueueConfigurations(tokens, coreService)
		if err != nil {
			logrus.Fatal("Could not query notification queues: " + err.Error())
		}

		for _, notificationQueueSettings := range queues {
			logrus.Infof("Starting to monitor notification queue: %+v", notificationQueueSettings)

			notificationQueue, err := queue.LoadNotificationQueue(configuration, &notificationQueueSettings)
			if err != nil {
				logrus.Fatal("Could not get notification queue: " + err.Error())
			}

			// funnel messages from each notification queue into the common channel
			go func(q queue.Queue) {
				for {
					msg, more := <-q.Receive()
					if !more {
						return
					}
					notifications <- msg
				}
			}(notificationQueue)
		}
	}

	// Start goroutines to funnel messages from each populated core service config's SyncObjectQueue channel into common channel
	for _, popConfig := range populatedConfigs.populatedConfigs {
		go func(c chan config.Message) {
			for {
				msg := <-c
				notifications <- msg
			}
		}(popConfig.SyncObjectQueue)
	}

	// Main loop
	for {
		select {
		case msg := <-notifications:
			if msg.RequireReload() {
				// This message contains information that requires the latest Sync Configuration information
				// to be correctly processed. Calling populatedConfigs.ForceReload() will signal the UpdateMuxer()
				// to signal the Update Monitors to requery their Core Services for the latest Sync
				// Configuration information.
				logrus.Infof("Reloading sync configuration because of message %s", msg.String())

				populatedConfigs.ForceReload()
			}

			// Route the incoming notification message to the appropriate worker(s)
			dispatched := false
			for _, populatedConfig := range populatedConfigs.GetConfigs() {
				is_match, should_ignore := msg.Match(populatedConfig)
				if is_match {
					if !should_ignore {
						populatedConfig.WorkerQueue <- msg
					}
					dispatched = true
					break
				}
			}

			if !dispatched {
				logrus.Error("Could not find core service configuration for message: " + msg.String())
			}
		case <-ctx.Done():
			logrus.Infof("Demuxer stopping...")
			return
		}
	}
}

// QueryQueueConfigurations queries the given core service for information on the Bucket and API notification queues that need to be monitored
func QueryQueueConfigurations(tokens service.RenewingTokens, coreService string) ([]config.NotificationQueueConfig, error) {
	// Hack to support running on localhost
	coreService = strings.Replace(coreService, "localhost/core", "core:8080", 1)

	req, err := http.NewRequest("GET", coreService+"/configuration/queue", nil)
	if err != nil {
		return nil, errors.New("could not create queue configuration request: " + err.Error())
	}

	idToken, err := tokens.GetIDToken()
	if err != nil {
		return nil, errors.New("could not get service ID Token for authentication: " + err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not make queue configuration request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 304 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with queue configuration response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with queue configuration response: StatusCode != 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("problem with reading queue configuration response: " + err.Error())
	}

	queueConfigurations := []config.NotificationQueueConfig{}
	if err := json.Unmarshal(body, &queueConfigurations); err != nil {
		return nil, errors.New("problem unmarshaling the queue configurations response: " + err.Error())
	}

	return queueConfigurations, nil
}
