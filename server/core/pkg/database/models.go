package database

import (
	"fmt"
	"os"
	"time"
)

const (
	PERM_READ       = "r"
	PERM_READ_WRITE = "rw"
)

const (
	OBJECT_STORE_TYPE_MINIO = "minio"
	OBJECT_STORE_TYPE_S3    = "s3"
)

const (
	DEFAULT_GROUP_SUFFIX = "hoss-default-group"
)

const (
	SYNC_TYPE_SIMPLEX = "simplex"
	SYNC_TYPE_DUPLEX  = "duplex"
)

// User object containing a user's username and the group they're a part of
type User struct {
	Id int64 `json:"-"`
	// Username is the unique username for the user
	Username string `json:"username"`

	// Memberships is a list of membership records indicating what groups the user is in
	Memberships []*Membership `json:"memberships" pg:"rel:has-many"`
}

// String prints the user record
func (u User) String() string {
	return fmt.Sprintf("User<%d %s>", u.Id, u.Username)
}

// Group object containing a list of user's who are part of the group and the dataset permissions granted to the group
type Group struct {
	Id int64 `json:"-"`
	// GroupName is the unique name of the group. Each user gets a 'default' group created
	// that they are the only member of. All permissions are then on groups.
	GroupName string `json:"group_name"`

	// Memberships is a list of membership records
	Memberships []*Membership `json:"memberships,omitempty" pg:"rel:has-many"`
	// Permissions is a list of permission records
	Permissions []*Permission `json:"permissions,omitempty" pg:"rel:has-many"`
}

// String prints the group record
func (g Group) String() string {
	return fmt.Sprintf("Group<%d %s>", g.Id, g.GroupName)
}

// Membership object connects Users to Groups they are members of
type Membership struct {
	GroupId int64 `json:"-"`
	// Group is the group this membership relationship is with
	Group *Group `json:"group,omitempty" pg:"rel:has-one"`

	UserId int64 `json:"-"`
	// User is the user who is a member of the indicated group
	User *User `json:"user,omitempty" pg:"rel:has-one"`
}

// ObjectStore object containing an object store's configuration information
// @Description Object store configuration information
type ObjectStore struct {
	Id int64 `json:"-"`
	// Name is the unique name given to identify an object store
	Name string `json:"name"`
	// Description is a short description of the object store
	Description string `json:"description"`
	// Endpoint is the object store host (e.g. https://s3.amazonaws.com, http://localhost)
	Endpoint string `json:"endpoint"`
	// ObjectStoreType is the type of object store, currently 'minio' or 's3' is supported
	ObjectStoreType string `json:"type" pg:"type:objectstoretype"`

	// ** Fields after this point are "optional" and not all stores may implement them **
	// Region is the region the object store is in if applicable
	Region string `json:"region,omitempty"`
	// Profile is the name of the profile to load from the ~/.aws/credentials file
	Profile string `json:"profile,omitempty"`
	// RoleArn is the arn of the role to use with STS if applicable
	RoleArn string `json:"role_arn,omitempty"`
	// NotificationArn is the arn of the queue to receive BucketNotification events if applicable
	NotificationArn string `json:"notification_arn,omitempty"`
}

// String prints the ObjectStore record
func (os ObjectStore) String() string {
	return fmt.Sprintf("ObjectStore<%d %s>", os.Id, os.Name)
}

// Namespace object containing a namespace's configuration information
// @Description Namespace configuration information
type Namespace struct {
	Id int64 `json:"-"`
	// Name is the unique name given to identify a namespace
	Name string `json:"name"`
	// Description is a short description of the namespace
	Description   string `json:"description"`
	ObjectStoreId int64  `json:"-"`
	// ObjectStore is the object store this namespace is backed by
	ObjectStore ObjectStore `json:"object_store" pg:"rel:has-one"`
	// BucketName is the name of the bucket in the object store where this namespace's datasets are stored
	BucketName string `json:"bucket_name"`
}

// String prints the namespace record
func (ns Namespace) String() string {
	return fmt.Sprintf("Namespace<%d %s>", ns.Id, ns.Name)
}

// Dataset object containing a dataset's information and the users who have been granted permissions
type Dataset struct {
	Id          int64 `json:"-"`
	NamespaceId int64 `json:"-"`
	//Namespace is the namespace this dataset is in
	Namespace *Namespace `json:"namespace" pg:"rel:has-one"`
	// Name is the unique name inside the namespace for this dataset
	Name string `json:"name"`
	// Description is a short description about the dataset
	Description string `json:"description"`
	// Created is the UTC datetime when the dataset was created
	Created time.Time `json:"created"`
	// RootDirectory is the prefix inside the namespace's bucket where this dataset is stored
	RootDirectory string `json:"root_directory"`
	// DeleteOn is the datetime when the delete will occur
	DeleteOn time.Time `json:"delete_on"`
	// DeleteStatus indicates if this dataset is marked for delete, in process, etc. ('NOT_SCHEDULED','SCHEDULED', 'IN_PROGRESS', 'ERROR')
	DeleteStatus string `json:"delete_status"`

	OwnerId int64 `json:"-"`
	// Owner is the user who created the dataset.
	Owner *User `json:"owner" pg:"rel:has-one"`

	// SyncEnabled is a flag indicating if syncing this dataset is enabled
	SyncEnabled bool `json:"sync_enabled" pg:",use_zero"`
	// SyncType is the type of sync relationship, if SyncEnabled is true ('simplex' or 'duplex')
	SyncType string `json:"sync_type" pg:",use_zero"`
	// SyncPolicy is the sync policy document that filters while messages should be synced
	SyncPolicy string `json:"sync_policy" pg:",use_zero"`

	// Permissions is a list of permission relationships between groups and this dataset
	Permissions []*Permission `json:"permissions,omitempty" pg:"rel:has-many"`
}

// String prints the dataset record
func (ds Dataset) String() string {
	return fmt.Sprintf("Dataset<%d %d %s>", ds.Id, ds.NamespaceId, ds.Name)
}

// Permission object is the mapping between a group and a dataset, with the type of permission granted
type Permission struct {
	GroupId int64 `json:"-"`
	// Group is the Group in this permission relationship
	Group *Group `json:"group,omitempty" pg:"rel:has-one"`

	DatasetId int64 `json:"-"`
	// Datset is the dataset in this permission relationship
	Dataset *Dataset `json:"dataset,omitempty" pg:"rel:has-one"`

	// Permission is the type of relationship ('r' or 'rw')
	Permission string `json:"permission" pg:"type:permission"`
}

// String prints the permission record
func (p Permission) String() string {
	var user string
	if p.Group != nil {
		user = p.Group.String()
	} else {
		user = fmt.Sprintf("%d", p.GroupId)
	}

	var dataset string
	if p.Dataset != nil {
		dataset = p.Dataset.String()
	} else {
		dataset = fmt.Sprintf("%d", p.DatasetId)
	}

	return fmt.Sprintf("Permission<%s %s %s>", user, dataset, p.Permission)
}

// SyncConfiguration defines the source and destination for sync notification messages
type SyncConfiguration struct {
	Id int64 `json:"-"`

	SourceNamespaceId int64      `json:"-"`
	SourceNamespace   *Namespace `json:"-" pg:"rel:has-one"`

	// TargetCoreService is the url to the core service that contains the namespace to which you
	// are linking this namespace. It can be the same or different server. (e.g. https://hoss.mycompany.com/core/v1)
	TargetCoreService string `json:"target_core_service"`
	// TargetNamespace is the name of the namespace to which you are linking this namespace
	TargetNamespace string `json:"target_namespace"`
	// SyncType is the type of sync relationship to configure ('simplex' or 'duplex')
	SyncType string `json:"sync_type" pg:",use_zero"`
}

// String prints the sync configuration record
func (sc *SyncConfiguration) String() string {
	flag := "->"
	if sc.SyncType == SYNC_TYPE_DUPLEX {
		flag = "<->"
	}

	return fmt.Sprintf("<SyncConfiguration %s:%s %s %s:%s>",
		os.Getenv("EXTERNAL_HOSTNAME")+"/core/v1",
		sc.SourceNamespace.Name,
		flag,
		sc.TargetCoreService,
		sc.TargetNamespace,
	)
}

// SyncConfigurationMeta holds the LastUpdated timestamp that defines
// when the sync configuration information was last modified in the database
type SyncConfigurationMeta struct {
	Id int64 `json:"-"`

	LastUpdate time.Time `json:"last_updated"`
}
