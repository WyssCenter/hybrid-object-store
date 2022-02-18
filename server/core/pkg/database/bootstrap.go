package database

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-core/pkg/config"
)

// BootstrapDefaults creates the default object store and namespace if they do not exist
func BootstrapDefaults(c *config.Configuration, db *Database) error {

	for _, objStore := range c.ObjectStores {
		// Check if object store exists
		_, err := db.GetObjectStore(objStore.Name)
		if err != nil {
			// Store doesn't exist. Try to create it
			logrus.Infof("Bootstrapping object store `%s`.", objStore.Name)

			err = db.CreateObjectStore(objStore.Name,
				objStore.Description,
				objStore.Type,
				objStore.Endpoint,
				objStore.Region,
				objStore.Profile,
				objStore.RoleArn,
				objStore.NotificationArn)
			if err != nil {
				return errors.Wrap(err, "Failed to bootstrap default object store.")
			}
		}
	}

	for _, ns := range c.Namespaces {
		// Check if namespace exists
		_, err := db.GetNamespace(ns.Name)
		if err != nil {
			// Namespace doesn't exist. Try to create it
			logrus.Infof("Bootstrapping namespace `%s`.", ns.Name)

			err = db.CreateNamespace(ns.Name,
				ns.Description,
				ns.ObjectStore,
				ns.Bucket)
			if err != nil {
				return errors.Wrap(err, "Failed to bootstrap default namespace.")
			}
		}
	}

	return nil
}
