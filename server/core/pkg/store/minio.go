package store

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStore is an Object Store type for interfacing with Minio
type MinioStore struct {
	client *minio.Client
	config *config.Configuration
	store  *database.ObjectStore
}

const policyTemplateMinio = `
{
    "Version": "2012-10-17",
    "Statement": [
        {{STATEMENTS}}
    ]  
}
`

const readTemplateMinio = `{
    "Effect": "Allow",
    "Action": ["s3:GetObject"],
    "Resource": ["arn:aws:s3:::{{.Bucket}}/{{.DatasetName}}/*"]
},
{
    "Effect": "Allow",
    "Action": ["s3:ListBucket"],
    "Resource": "arn:aws:s3:::{{.Bucket}}",
    "Condition": {
        "StringLike": {
        "s3:prefix": "{{.DatasetName}}/*"
        }
    }
}
`

const readWriteTemplateMinio = `{
    "Effect": "Allow",
    "Action": ["s3:*"],
    "Resource": ["arn:aws:s3:::{{.Bucket}}/{{.DatasetName}}/*"]
},
{
    "Effect": "Allow",
    "Action": ["s3:ListBucket"],
    "Resource": ["arn:aws:s3:::{{.Bucket}}"],
    "Condition": {
        "StringLike": {
        "s3:prefix": ["{{.DatasetName}}/*"]
        }
    }
}
`

const denyTemplateMinio = `{
    "Effect": "Deny",
    "Action": ["s3:*"],
    "Resource": ["arn:aws:s3:::*"]
}
`

// GetType returns the type of the ObjectStore for handling interface values
func (m *MinioStore) GetType() string {
	return database.OBJECT_STORE_TYPE_MINIO
}

// GetName returns the type of the ObjectStore for handling interface values
func (m *MinioStore) GetName() string {
	return m.store.Name
}

// UserPolicyName returns the name of a canned/session policy for a user
func (m *MinioStore) UserPolicyName(username string) string {
	return username
}

// Load returns an object store based on a namespace's configuration
// Additionally any 1-time configuration, session setup, etc. should
// happen here. After this function is called, the ObjectStore should
// be ready for use.
func (m *MinioStore) Load(c *config.Configuration, o *database.ObjectStore) error {
	if c == nil {
		return errors.New("Failed to load object store. Configuration is required")
	}
	if o == nil {
		return errors.New("Failed to load object store. Object Store is required")
	}

	m.config = c
	m.store = o

	accessKey, ok := os.LookupEnv("MINIO_ROOT_USER")
	if !ok {
		return errors.New("Failed to load object store. Access Key env var is missing.")
	}

	secretKey, ok := os.LookupEnv("MINIO_ROOT_PASSWORD")
	if !ok {
		return errors.New("Failed to load object store. Secret Key env var is missing.")
	}

	endpoint, err := m.getReachableEndpoint()
	if err != nil {
		return errors.Wrap(err, "failed to get endpoint while bootstrapping minio client")
	}

	urlParsed, err := url.Parse(endpoint)
	if err != nil {
		return errors.Wrap(err, "Failed to parse object store endpoint")
	}

	useSSL := false
	if urlParsed.Scheme == "https" {
		useSSL = true
	}

	alias, err := m.getAlias()
	if err != nil {
		return errors.Wrap(err, "failed to get mc alias while bootstrapping minio client.")
	}

	isMinioClientReady := false
	for i := 0; i < 10; i++ {
		logrus.Infof("Configuring client minio for: %v, %v", alias, endpoint)

		cmd := exec.Command("mc", "alias", "set", alias, endpoint,
			accessKey, secretKey)
		err := cmd.Run()
		if err == nil {
			isMinioClientReady = true
			break
		}
		logrus.Infof("Minio is not ready, sleeping. - Failed to configure mc: %v", err)
		time.Sleep(5 * time.Second)
	}
	if !isMinioClientReady {
		logrus.Fatal("Couldn't connect to minio after 50 seconds")
	}

	// Initialize minio client object.
	minioClient, err := minio.New(urlParsed.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to load Minio client")
	}

	m.client = minioClient

	return nil
}

// CreateDataset creates a root folder for a dataset inside a namespace
func (m *MinioStore) CreateDataset(name string, n *database.Namespace) error {

	metadata := NewMetadataFile(name)

	//Check if dataset already exists
	_, err := m.client.StatObject(context.TODO(),
		n.BucketName,
		metadata.Key(),
		minio.StatObjectOptions{})
	if err == nil {
		return errors.New("Failed to create dataset. Dataset already exists")
	}

	if merr, ok := err.(minio.ErrorResponse); !ok {
		return errors.Wrap(err, "When creating new dataset, failed to check if dataset exists")
	} else {
		if merr.Code != "NoSuchKey" {
			return errors.Wrap(err, "When creating new dataset, failed to check if dataset exists")
		}
	}

	// Create dataset metadata file to "create" the dataset
	metaRaw, err := yaml.Marshal(metadata)
	if err != nil {
		return errors.Wrap(err, "When creating new dataset, failed to write metadata")
	}
	metaReader := bytes.NewReader(metaRaw)

	putOpts := minio.PutObjectOptions{}
	_, err = m.client.PutObject(context.TODO(),
		n.BucketName,
		metadata.Key(),
		metaReader,
		int64(len(metaRaw)),
		putOpts)
	if err != nil {
		return errors.Wrap(err, "When creating new dataset, failed to write metadata")
	}

	return nil
}

// DeleteDataset deletes a dataset directory from an object store within a
// namespace. prefix is typically set to the rootDirectory value from the
// database, which includes a trailing slash.
func (m *MinioStore) DeleteDataset(prefix string, n *database.Namespace) error {
	// TODO: Future work should shift this to workers as a dataset may have many
	// 	 	 objects and take a long time to delete.

	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for obj := range m.client.ListObjects(context.Background(),
			n.BucketName, minio.ListObjectsOptions{
				Prefix:    prefix,
				Recursive: true,
			}) {

			if obj.Err != nil {
				log.Fatalln(obj.Err)
			}
			objectsCh <- obj
		}
	}()

	removeOpts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	var errKeys []string
	for rErr := range m.client.RemoveObjects(context.Background(), n.BucketName,
		objectsCh, removeOpts) {
		errKeys = append(errKeys, rErr.ObjectName)
	}
	if len(errKeys) > 0 {
		errStr := strings.Join(errKeys[:], ",")
		return errors.New(fmt.Sprintf("Failed to delete some objects while removing the dataset `%s`: %s", prefix, errStr))
	}

	return nil
}

// SetUserPolicy re-renders a user's policy and applies it to the store
func (m *MinioStore) SetUserPolicy(username string, permissions []*database.Permission) error {
	// Render policy
	policy, err := RenderTemplate(policyTemplateMinio, readTemplateMinio,
		readWriteTemplateMinio, denyTemplateMinio, permissions)
	if err != nil {
		return errors.Wrap(err, "Failed to update policy")
	}

	alias, err := m.getAlias()
	if err != nil {
		return errors.Wrap(err, "failed to remove user policy")
	}

	// Check if the policy has actually changed
	var stdOutBufChk bytes.Buffer
	var stdErrBufChk bytes.Buffer
	cmdChk := exec.Command("mc", "admin", "policy", "info", alias, m.UserPolicyName(username))
	cmdChk.Stdout = &stdOutBufChk
	cmdChk.Stderr = &stdErrBufChk
	cmdErr := cmdChk.Run()
	if cmdErr != nil {
		// A not found error always occurs on first policy gen. Ignore this error.
		if strings.Contains(stdErrBufChk.String(), "The canned policy does not exist") {
			// The policy has never been created so just move on and create it.
		} else if strings.Contains(stdErrBufChk.String(), "not found") {
			// The policy was not found so just move on and create it.
		} else {
			// Something else happened
			return errors.Wrap(cmdErr, fmt.Sprintf("Failed to verify existing policy before update:%s\n\n%s", cmdChk.Stdout, cmdChk.Stderr))
		}
	}

	// Extract existing policy and strip all spaces
	existingPolicy := trimWhitespace(stdOutBufChk.String())

	if trimWhitespace(policy) == existingPolicy {
		logrus.Infof("Policy for %s is up-to-date. Skipping policy update.", username)
		return nil
	}

	// Write policy to temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp-policy-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
		return errors.Wrap(err, "Failed to update policy")
	}
	if _, err = tmpFile.Write([]byte(policy)); err != nil {
		return errors.Wrap(err, "Failed to write policy file")
	}
	if err := tmpFile.Close(); err != nil {
		return errors.Wrap(err, "Failed to write policy file")
	}

	// Add the canned policy to minio using mc
	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	cmd := exec.Command("mc", "admin", "policy", "add", alias,
		m.UserPolicyName(username), tmpFile.Name())
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf
	cmdErr = cmd.Run()

	// Remove temp file after you run the command, regardless of error state
	os.Remove(tmpFile.Name())

	if cmdErr != nil {
		return errors.Wrap(cmdErr, fmt.Sprintf("Failed to apply policy:%s\n\n%s", cmd.Stdout, cmd.Stderr))
	}

	return nil
}

//DeleteUserPolicy completely removes a canned policy for a user from the system
func (m *MinioStore) DeleteUserPolicy(username string) error {
	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer

	alias, err := m.getAlias()
	if err != nil {
		return errors.Wrap(err, "failed to remove user policy")
	}

	cmd := exec.Command("mc", "admin", "policy", "remove", alias, m.UserPolicyName(username))
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to remove user policy for %s:%s\n\n%s", username,
			cmd.Stdout, cmd.Stderr))
	}

	return nil
}

// GetSTSCredentials gets temporary STS credentials for the object store for the current user
func (m *MinioStore) GetSTSCredentials(jwt string, claims map[string]interface{}, username string) (*Credentials, error) {
	exp := int(claims["exp"].(float64))
	expireInSeconds := exp - int(time.Now().UTC().Unix())
	expireAt := time.Now().UTC().Add(time.Second * time.Duration(expireInSeconds))

	// minio sts callback function
	getWebTokenExpiryFunc := func() (*credentials.WebIdentityToken, error) {
		return &credentials.WebIdentityToken{
			Token:  jwt,
			Expiry: expireInSeconds, // amount of time left before token expiration in seconds
		}, nil
	}

	sts, err := credentials.NewSTSWebIdentity(m.client.EndpointURL().String(), getWebTokenExpiryFunc)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create new STS Web Identity")
	}

	creds, err := sts.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get temp STS credentials")
	}

	expiration := expireAt.Format(time.RFC3339)
	response := Credentials{
		AccessKeyId:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		Expiration:      expiration,
		Region:          m.store.Region,
		Endpoint:        m.store.Endpoint,
	}

	return &response, nil
}

// EventsEnabled checks to see if Bucket Events have been enabled for the Dataset
func (m *MinioStore) EventsEnabled(namespace *database.Namespace, dataset *database.Dataset) (bool, error) {
	bucket := namespace.BucketName

	bucketNotification, err := m.client.GetBucketNotification(context.Background(), bucket)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get bucket notification configuration")
	}

	for _, queueConfig := range bucketNotification.QueueConfigs {
		for _, rule := range queueConfig.Filter.S3Key.FilterRules {
			if rule.Name == "prefix" && rule.Value == dataset.RootDirectory {
				return true, nil
			}
		}
	}

	return false, nil
}

// EnableEvents starts Minio sending Bucket Events for the Dataset files
func (m *MinioStore) EnableEvents(namespace *database.Namespace, dataset *database.Dataset) error {
	bucket := namespace.BucketName // DP ???: use m.namespace instead of passing it as an argument?
	arn := "arn:minio:sqs::_:amqp"
	alias, err := m.getAlias()
	if err != nil {
		return err
	}

	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	cmd := exec.Command("mc", "event", "add", alias+"/"+bucket, arn, "--prefix", dataset.RootDirectory)
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to enable bucket events for %s:%s\n\n%s", arn,
			cmd.Stdout, cmd.Stderr))
	}

	return nil
}

// DisableEvents stops Minio sending Bucket Events for the Dataset files
func (m *MinioStore) DisableEvents(namespace *database.Namespace, dataset *database.Dataset) error {
	bucket := namespace.BucketName
	arn := "arn:minio:sqs::_:amqp"
	alias, err := m.getAlias()
	if err != nil {
		return err
	}

	var stdOutBuf bytes.Buffer
	var stdErrBuf bytes.Buffer
	// `--event` flag explicitly included because mc is not properly setting default values for the
	// `remove` command like it does for the `add` command. Can remove when fixed: https://github.com/minio/mc/pull/3752
	cmd := exec.Command("mc", "event", "remove", alias+"/"+bucket, arn, "--prefix", dataset.RootDirectory, "--event", "put,delete,get")
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to disable bucket events for %s:%s\n\n%s", arn,
			cmd.Stdout, cmd.Stderr))
	}

	return nil
}

// getAlias returns the mc alias for this minio struct
func (m *MinioStore) getAlias() (string, error) {
	if m.store == nil {
		return "", errors.New("Cannot get mc alias without first loading the store.")
	}

	return "hoss-" + m.store.Name, nil
}

// getReachableEndpoint is a helper function that converts localhost to the internal DNS route to minio.
//   This is used to handle the case where you are running on localhost and can't get to minio via the proxy
//   that is running on localhost.
func (m *MinioStore) getReachableEndpoint() (string, error) {
	u, err := url.Parse(m.store.Endpoint)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse endpoint while generating internal minio endpoint")
	}

	// If the "TESTING" env var is set, keep using localhost so debugging/tests running on host work.
	_, isTesting := os.LookupEnv("TESTING")

	if !isTesting && u.Host == "localhost" {
		// This namespace is using minio running on localhost, set to internal endpoint
		return "http://minio:9000", nil
	}

	return m.store.Endpoint, nil
}

// PatchMinioEvents lists all datasets in a namespace that is backed by a minio
// object store. Then it will disable/enable bucket events on the datasets.
// For some unknown reason this is required when running minio in gateway
// mode on a container restart. This issue is tracked here:
// https://github.com/minio/minio/issues/13816
// If this is fixed upstream, we can remove this functionality.
func PatchMinioEvents(db *database.Database, stores map[string]ObjectStore) error {
	nsList, err := db.ListNamespaces(1000, 0)
	if err != nil {
		return errors.Wrap(err, "failed to list namespaces while trying to patch minio events")
	}

	for _, ns := range nsList {
		if ns.ObjectStore.ObjectStoreType == "minio" {
			logrus.Infof("Patching minio events for namespace: %s", ns.Name)

			dsList, err := db.ListDatasetsInNamespace(ns)
			if err != nil {
				return errors.Wrap(err, "failed to list datasets in namespace while trying to patch minio events")
			}
			for _, ds := range dsList {
				currentStore, ok := stores[ds.Namespace.ObjectStore.Name]
				if !ok {
					return errors.Wrap(err, "failed to load object store")
				}

				// If events are enabled, disable them
				isEnabled, err := currentStore.EventsEnabled(ns, ds)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("failed to check if dataset %s sync is enabled: %s", ds.Name, err.Error()))
				}
				if isEnabled {
					err = currentStore.DisableEvents(ns, ds)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("Failed to disable dataset %s sync: %s", ds.Name, err.Error()))
					}
				}

				// Enable events
				err = currentStore.EnableEvents(ns, ds)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("Failed to enalbe dataset %s sync: %s", ds.Name, err.Error()))
				}

			}
		}
	}
	return nil
}
