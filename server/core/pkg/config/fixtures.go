package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func SetupMinioConfigTest(t *testing.T) (string, error) {
	defaultobj := &ObjectStore{
		Name:        "default",
		Description: "Default object store test",
		Type:        "minio",
		Endpoint:    "http://localhost",
	}
	defaultns := &Namespace{
		Name:        "default",
		Description: "Default namespace test",
		Bucket:      "data-test-bucket",
		ObjectStore: "default",
	}

	filename, err := writeTestConfigFile(defaultns, defaultobj)
	if err != nil {
		errors.Wrap(err, "failed to create minio test config file")
	}

	t.Cleanup(func() {
		os.Remove(filename)
	})

	return filename, nil
}

func SetupS3ConfigTest(t *testing.T) (string, error) {
	defaultobj := &ObjectStore{
		Name:        "default",
		Description: "Default object store",
		Type:        "s3",
		Endpoint:    "https://s3.amazonaws.com",
		Region:      "us-east-1",
		Profile:     "test-creds-1",
		RoleArn:     "arn:aws:iam::123456789012:role/myHossUserRole",
	}
	defaultns := &Namespace{
		Name:        "default",
		Description: "Default namespace",
		Bucket:      "my-default-bucket-1",
		ObjectStore: "default",
	}

	filename, err := writeTestConfigFile(defaultns, defaultobj)
	if err != nil {
		errors.Wrap(err, "failed to create minio test config file")
	}

	t.Cleanup(func() {
		os.Remove(filename)
	})

	return filename, nil

}

func writeTestConfigFile(ns *Namespace, objs *ObjectStore) (string, error) {

	tf, err := ioutil.TempFile("", "test_config_*.yaml")
	if err != nil {
		errors.Wrap(err, "failed to create test config file")
	}

	objectStores := []ObjectStore{}
	objectStores = append(objectStores, *objs)

	namespaces := []Namespace{}
	namespaces = append(namespaces, *ns)

	server := &Server{Dev: true}
	conf := &Configuration{ObjectStores: objectStores, Namespaces: namespaces, Server: *server}

	raw, err := yaml.Marshal(conf)
	if err != nil {
		errors.Wrap(err, "failed to marshal test config file")
	}
	tf.Write(raw)
	tf.Close()

	return tf.Name(), nil
}
