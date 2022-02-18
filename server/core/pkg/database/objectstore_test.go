package database

import (
	"testing"

	"github.com/gigantum/hoss-core/pkg/test"
)

func TestCreateObjectStoreMinio(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("my_object_store",
		"My description",
		"minio",
		"https://mydomain.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	objLoaded, err := db.GetObjectStore("my_object_store")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}

	test.AssertEqual(t, objLoaded.Name, "my_object_store")
	test.AssertEqual(t, objLoaded.Description, "My description")
	test.AssertEqual(t, objLoaded.Endpoint, "https://mydomain.com")
	test.AssertEqual(t, objLoaded.ObjectStoreType, "minio")
	test.AssertEqual(t, objLoaded.Profile, "")
	test.AssertEqual(t, objLoaded.Region, "")
}

func TestCreateObjectStoreInvalid(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("my_object_store",
		"My description",
		"not-a-supported-type",
		"https://mydomain.com",
		"",
		"",
		"",
		"")
	if err == nil {
		t.Fatal("Expected to fail to create object store")
	}

}

func TestCreateObjectStoreExists(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("my_object_store",
		"My description",
		"minio",
		"https://mydomain.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	err = db.CreateObjectStore("my_object_store",
		"My description 2",
		"minio",
		"https://mydomain2.com",
		"",
		"",
		"",
		"")
	if err == nil {
		t.Fatal("Created namespace again when an exception should have been raised")
	}

	test.AssertEqual(t, err.Error(), "record already exists")

}

func TestGetObjectStoreDoesNotExist(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err = db.GetObjectStore("my_object_store_does_not_exist")
	if err == nil {
		t.Fatal("Expected error because object store does not exist")
	}

	test.AssertEqual(t, err.Error(), "record not found")

}

func TestDeleteObjectStore(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("my_object_store",
		"My description",
		"minio",
		"https://mydomain.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	_, err = db.GetObjectStore("my_object_store")
	if err != nil {
		t.Fatal("Expected object Store to exist")
	}

	err = db.DeleteObjectStore("my_object_store")
	if err != nil {
		t.Fatal("Expected no error but delete object store failed")
	}

	_, err = db.GetObjectStore("my_object_store")
	if err == nil {
		t.Fatal("Expected error because object store should not exist")
	}
	test.AssertEqual(t, err.Error(), "record not found")
}

func TestListObjectStores(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("obj1",
		"My description 1",
		"minio",
		"https://mydomain1.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	err = db.CreateObjectStore("obj2",
		"My description 2",
		"minio",
		"https://mydomain2.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	err = db.CreateObjectStore("obj3",
		"My description 3",
		"minio",
		"https://mydomain3.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	err = db.CreateObjectStore("obj4",
		"My description 4",
		"minio",
		"https://mydomain4.com",
		"",
		"",
		"",
		"")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	objs, err := db.ListObjectStores(25, 0)
	if err != nil {
		t.Fatal("Failed to list object stores")
	}
	if len(objs) != 5 {
		t.Fatalf("Length of objs is wrong, %v != 5", len(objs))
	}
	test.AssertEqual(t, objs[1].Name, "obj1")
	test.AssertEqual(t, objs[4].Name, "obj4")

	objs, err = db.ListObjectStores(25, 1)
	if err != nil {
		t.Fatal("Failed to list object stores")
	}
	if len(objs) != 4 {
		t.Fatalf("Length of objs is wrong, %v != 4", len(objs))
	}
	test.AssertEqual(t, objs[0].Name, "obj1")
	test.AssertEqual(t, objs[3].Name, "obj4")

	objs, err = db.ListObjectStores(2, 2)
	if err != nil {
		t.Fatal("Failed to list objs")
	}
	if len(objs) != 2 {
		t.Fatalf("Length of objs is wrong, %v != 2", len(objs))
	}
	test.AssertEqual(t, objs[0].Name, "obj2")
	test.AssertEqual(t, objs[1].Name, "obj3")

}

func TestCreateObjectStoreS3(t *testing.T) {
	db, err := SetupDatabaseTest(t)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	err = db.CreateObjectStore("my_object_store",
		"My description",
		"s3",
		"https://s3.awsamazon.com",
		"us-east-1",
		"my-hoss-user",
		"arn:aws:iam::123456789012:role/myHossUserRole",
		"arn:aws:sqs::123456789012:hoss-notifications")
	if err != nil {
		t.Fatal("Failed to create object store")
	}

	objLoaded, err := db.GetObjectStore("my_object_store")
	if err != nil {
		t.Fatal("Expected no error but get namespace failed")
	}

	test.AssertEqual(t, objLoaded.Name, "my_object_store")
	test.AssertEqual(t, objLoaded.Description, "My description")
	test.AssertEqual(t, objLoaded.Endpoint, "https://s3.awsamazon.com")
	test.AssertEqual(t, objLoaded.ObjectStoreType, "s3")
	test.AssertEqual(t, objLoaded.Region, "us-east-1")
	test.AssertEqual(t, objLoaded.Profile, "my-hoss-user")
	test.AssertEqual(t, objLoaded.RoleArn, "arn:aws:iam::123456789012:role/myHossUserRole")
	test.AssertEqual(t, objLoaded.NotificationArn, "arn:aws:sqs::123456789012:hoss-notifications")
}
