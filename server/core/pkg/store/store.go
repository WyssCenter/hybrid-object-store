package store

import (
	"log"
	"strings"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
)

type Credentials struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	Expiration      string `json:"expiration"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
}

// ObjectStore interface defines functions required for an object store implementation
type ObjectStore interface {
	// Load returns an object store based on a namespace's configuration
	// Additionally any 1-time configuration, session setup, etc. should
	// happen here. After this function is called, the ObjectStore should
	// be ready for use.
	Load(c *config.Configuration, o *database.ObjectStore) error

	// CreateDataset creates a root folder for a dataset inside a namespace
	CreateDataset(name string, n *database.Namespace) error

	// DeleteDataset deletes a dataset folder
	DeleteDataset(name string, n *database.Namespace) error

	// SetUserPolicy renders and sets the policy in the object store for a user based on permissions
	SetUserPolicy(username string, permissions []*database.Permission) error

	// DeleteUserPolicy nulls the policy but keeps it present.
	DeleteUserPolicy(username string) error

	// GetSTSCredentials exchanges a JWT for object store creds
	GetSTSCredentials(jwt string, claims map[string]interface{}, username string) (*Credentials, error)

	// GetType returns the type of the ObjectStore for handling interface values
	GetType() string

	// GetName returns the name of the ObjectStore
	GetName() string

	// UserPolicyName returns the name of a canned policy for a user
	UserPolicyName(username string) string

	// EventsEnables checks to see if Bucket Event Notifications have been enabled for the given dataset
	EventsEnabled(namespace *database.Namespace, dataset *database.Dataset) (bool, error)

	// EnableEvents turns on Bucket Event Notifications for the given dataset
	EnableEvents(namespace *database.Namespace, dataset *database.Dataset) error

	// DisableEvents turns off Bucket Notifications for the given dataset
	DisableEvents(namespace *database.Namespace, dataset *database.Dataset) error
}

// LoadObjectStores is a helper method to load all object stores. Since things are pretty broken if object stores fail
//                  to load, we just log.Fatal and exit.
func LoadObjectStores(c *config.Configuration, objectStores []*database.ObjectStore) map[string]ObjectStore {
	objMap := map[string]ObjectStore{}

	for _, o := range objectStores {
		switch t := o.ObjectStoreType; t {
		case database.OBJECT_STORE_TYPE_MINIO:
			s := &MinioStore{}
			err := s.Load(c, o)
			if err != nil {
				log.Fatalf("Failed to load minio object store: %v", err)
			}
			var interfaceType ObjectStore = s
			objMap[o.Name] = interfaceType

		case database.OBJECT_STORE_TYPE_S3:
			s := &S3Store{}
			err := s.Load(c, o)
			if err != nil {
				log.Fatalf("Failed to load S3 object store: %v", err)
			}
			var interfaceType ObjectStore = s
			objMap[o.Name] = interfaceType

		default:
			log.Fatal("Unsupported store type encountered while loading object stores. Must be `minio` or `s3`")
		}

	}

	return objMap
}

// trimWhitespace is a helper function that removes all whitespace from a string
func trimWhitespace(s string) string {
	s1 := strings.TrimSpace(s)
	s2 := strings.ReplaceAll(s1, "\n", "")
	s3 := strings.ReplaceAll(s2, "\t", "")
	s4 := strings.ReplaceAll(s3, " ", "")
	return s4
}
