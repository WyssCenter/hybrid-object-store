package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration contains configuration info for the core service
type Configuration struct {
	ObjectStores []ObjectStore `yaml:"object_stores"`
	Namespaces   []Namespace   `yaml:"namespaces"`
	Queues       []Queue       `yaml:"queues"`
	Server       Server        `yaml:"server"`
}

// Namespace contains configuration info for a namespace that should be
// bootstrapped at start
type Namespace struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Bucket      string `yaml:"bucket"`
	ObjectStore string `yaml:"object_store"`
}

// ObjectStore contains configuration info for an object store that should be
// bootstrapped at start
type ObjectStore struct {
	Name            string `yaml:"name"`
	Description     string `yaml:"description"`
	Type            string `yaml:"type"`
	Endpoint        string `yaml:"endpoint"`
	Region          string `yaml:"region"`
	Profile         string `yaml:"profile"`
	RoleArn         string `yaml:"role_arn"`
	NotificationArn string `yaml:"notification_arn"`
}

type Queue struct {
	Type        string                 `json:"type"`
	Settings    map[string]interface{} `json:"settings"`
	ObjectStore string                 `yaml:"object_store"`
}

type AMQPQueueConfig struct {
	Url string `json:"url"`
}

type SQSQueueConfig struct {
	QueueName string `yaml:"queue_name"`
	Region    string `json:"region"`
	Profile   string `json:"profile"`
}

// Server contains configuration info for the core service
type Server struct {
	Dev                        bool   `yaml:"dev"`
	AuthService                string `yaml:"auth_service"`
	ElasticsearchEndpoint      string `yaml:"elasticsearch_endpoint"`
	SyncFrequencyMinutes       int    `yaml:"sync_frequency_minutes"`
	DatasetDeleteDelayMinutes  int    `yaml:"dataset_delete_delay_minutes"`
	DatasetDeletePeriodSeconds int    `yaml:"dataset_delete_period_seconds"`
}

// Load creates a default config and then initializes it with values from
// the default config file location.
func Load(path string) *Configuration {
	if path == "" {
		path = findConfig()
	}
	config := &Configuration{}
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		if err = yaml.Unmarshal(bytes, &config); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return config
}

// findConfig locates the config file (depending on testing/prod configs)
func findConfig() string {
	if _, err := os.Stat(filepath.Join("/opt", "config.yaml")); os.IsNotExist(err) {
		// `/opt/config.yaml` does not exist. Try to load relative
		path, err := filepath.Abs("../../config.yaml")
		if err != nil {
			log.Fatal(err)
		}
		return path
	}

	return filepath.Join("/opt", "config.yaml")
}

// UnmarshalSettings converts the generic map to the specific interface given
func UnmarshalSettings(settings map[string]interface{}, target interface{}) error {
	b, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, target)
}
