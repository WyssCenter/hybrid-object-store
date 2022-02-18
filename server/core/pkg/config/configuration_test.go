package config

import (
	"testing"

	"github.com/gigantum/hoss-core/pkg/test"
)

func TestMinioDefaultConfig(t *testing.T) {
	testfile, err := SetupMinioConfigTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	c := Load(testfile)

	test.AssertEqual(t, c.Server.Dev, true)

	test.AssertEqual(t, c.Namespaces[0].Name, "default")
	test.AssertEqual(t, c.Namespaces[0].Description, "Default namespace test")
	test.AssertEqual(t, c.Namespaces[0].Bucket, "data-test-bucket")
	test.AssertEqual(t, c.Namespaces[0].ObjectStore, "default")

	test.AssertEqual(t, c.ObjectStores[0].Name, "default")
	test.AssertEqual(t, c.ObjectStores[0].Description, "Default object store test")
	test.AssertEqual(t, c.ObjectStores[0].Type, "minio")
	test.AssertEqual(t, c.ObjectStores[0].Endpoint, "http://localhost")
	test.AssertEqual(t, c.ObjectStores[0].Region, "")
	test.AssertEqual(t, c.ObjectStores[0].Profile, "")
}

func TestS3DefaultConfig(t *testing.T) {
	testfile, err := SetupS3ConfigTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	c := Load(testfile)

	test.AssertEqual(t, c.Server.Dev, true)

	test.AssertEqual(t, c.Namespaces[0].Name, "default")
	test.AssertEqual(t, c.Namespaces[0].Description, "Default namespace")
	test.AssertEqual(t, c.Namespaces[0].Bucket, "my-default-bucket-1")
	test.AssertEqual(t, c.Namespaces[0].ObjectStore, "default")

	test.AssertEqual(t, c.ObjectStores[0].Name, "default")
	test.AssertEqual(t, c.ObjectStores[0].Description, "Default object store")
	test.AssertEqual(t, c.ObjectStores[0].Type, "s3")
	test.AssertEqual(t, c.ObjectStores[0].Endpoint, "https://s3.amazonaws.com")
	test.AssertEqual(t, c.ObjectStores[0].Region, "us-east-1")
	test.AssertEqual(t, c.ObjectStores[0].Profile, "test-creds-1")
	test.AssertEqual(t, c.ObjectStores[0].RoleArn, "arn:aws:iam::123456789012:role/myHossUserRole")

}

func TestDefaultConfig(t *testing.T) {
	c := Load("")

	test.AssertEqual(t, c.Server.Dev, true)

	test.AssertEqual(t, c.Namespaces[0].Name, "default")
	test.AssertEqual(t, c.Namespaces[0].Description, "Default namespace")
	test.AssertEqual(t, c.Namespaces[0].Bucket, "data")
	test.AssertEqual(t, c.Namespaces[0].ObjectStore, "default")

	test.AssertEqual(t, c.ObjectStores[0].Name, "default")
	test.AssertEqual(t, c.ObjectStores[0].Description, "Default object store")
	test.AssertEqual(t, c.ObjectStores[0].Type, "minio")
	test.AssertEqual(t, c.ObjectStores[0].Endpoint, "http://localhost")
	test.AssertEqual(t, c.ObjectStores[0].Region, "")
	test.AssertEqual(t, c.ObjectStores[0].Profile, "")
}
