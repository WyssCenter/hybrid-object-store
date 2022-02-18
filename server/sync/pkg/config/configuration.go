package config

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/ghodss/yaml"
)

// Load the given configuration file, if the file is "" then load from the default location
func Load(path string) *Configuration {
	if path == "" {
		path = "/opt/config.yaml"
	}

	config := &Configuration{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Could not read config file: " + err.Error())
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatal("Could not load config file: " + err.Error())
	}

	config.RefreshIntervals.CoreService, err = time.ParseDuration(config.RefreshIntervals.CoreServiceString)
	if err != nil {
		log.Fatalf("could not parse core_service refresh interval: %s", err.Error())
	}

	config.RefreshIntervals.AuthToken, err = time.ParseDuration(config.RefreshIntervals.AuthTokenString)
	if err != nil {
		log.Fatalf("could not parse auth_token refresh interval: %s", err.Error())
	}

	config.RefreshIntervals.StsCredentials, err = time.ParseDuration(config.RefreshIntervals.StsCredsString)
	if err != nil {
		log.Fatalf("could not parse sts_creds refresh interval: %s", err.Error())
	}

	if len(config.CoreServices) == 0 {
		log.Fatal("core_services: At least one Core Service must be defined in the config file")
	}

	if config.WorkerInstanceCount == 0 {
		log.Fatal("worker_instance_count: At least one worker per monitored Core Service must be defined")
	}

	return config
}

// UnmarshalSettings converts the generic map to the specific interface given
func UnmarshalSettings(settings map[string]interface{}, target interface{}) error {
	b, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, target)
}

// Configuration is the top level configuration for the Sync Service
type Configuration struct {
	CoreServices []string `json:"core_services"`

	RefreshIntervals RefreshIntervals `json:"refresh_intervals"`

	AuthEndpoint string `json:"auth_endpoint"`

	SqsProfile string `json:"sqs_profile"`

	WorkerBufferSize    int `json:"worker_buffer_size"`
	WorkerInstanceCount int `json:"worker_instance_count"` // per core service
}

// RefreshIntervals defines the refresh intervals for various credentials needed by the sync service
type RefreshIntervals struct {
	CoreServiceString string        `json:"core_service"`
	CoreService       time.Duration `json:"-"`
	AuthTokenString   string        `json:"auth_token"`
	AuthToken         time.Duration `json:"-"`
	StsCredsString    string        `json:"sts_creds"`
	StsCredentials    time.Duration `json:"-"`
}

// NotificationQueueConfig defines a queue to monitor for notifications
type NotificationQueueConfig struct {
	Type     string                 `json:"type"`
	Settings map[string]interface{} `json:"settings"`
}

// =============================================================================
// Implementation specific configs
// =============================================================================

type AMQPQueueConfig struct {
	// The type of message to unmarshal into correct type
	MessageType string `json:"message_type"`
	//
	SourceEndpoint string `json:"source_endpoint"`
	QueueName      string `json:"queue_name"`
	ExchangeName   string `json:"exchange_name"`
}

type SQSQueueConfig struct {
	// The type of message to unmarshal into correct type
	MessageType string `json:"message_type"`
	//
	SourceEndpoint string `json:"source_endpoint"`
	QueueName      string `json:"queue_name"`

	Region  string `json:"region"`
	Profile string `json:"profile"`
}
