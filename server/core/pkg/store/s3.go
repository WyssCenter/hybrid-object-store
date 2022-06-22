package store

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	ststypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// S3Store is an Object Store type for interfacing with AWS S3
type S3Store struct {
	awsConfig awsconfig.Config
	client    *s3.Client
	iam       *iam.Client
	sts       *sts.Client
	config    *config.Configuration
	store     *database.ObjectStore
	accountId string
}

const policyTemplateS3 = `{"Version": "2012-10-17","Statement": [{{STATEMENTS}}]}`

const readTemplateS3 = `{"Effect": "Allow","Action": "s3:GetObject","Resource": "arn:aws:s3:::{{.Bucket}}/{{.DatasetName}}/*"},{"Effect": "Allow","Action": "s3:ListBucket","Resource": "arn:aws:s3:::{{.Bucket}}","Condition": {"StringLike": {"s3:prefix": "{{.DatasetName}}/*"}}}`

// TODO: Lock down the read/write template more
const readWriteTemplateS3 = `{"Effect": "Allow","Action": "s3:*","Resource": "arn:aws:s3:::{{.Bucket}}/{{.DatasetName}}/*"},{"Effect": "Allow","Action": "s3:ListBucket","Resource": "arn:aws:s3:::{{.Bucket}}","Condition": {"StringLike": {"s3:prefix": "{{.DatasetName}}/*"}}}`

const denyTemplateS3 = `{"Effect": "Deny","Action": "s3:*","Resource": "arn:aws:s3:::*"}`

// GetType returns the type of the ObjectStore for handling interface values
func (s *S3Store) GetType() string {
	return database.OBJECT_STORE_TYPE_S3
}

// GetName returns the type of the ObjectStore for handling interface values
func (s *S3Store) GetName() string {
	return s.store.Name
}

// UserPolicyName returns the name of a canned/session policy for a user
func (s *S3Store) UserPolicyName(username string) string {
	return "hoss-user-policy-" + username
}

// Load returns an object store based on the provided namespace's configuration
func (s *S3Store) Load(c *config.Configuration, o *database.ObjectStore) error {
	if c == nil {
		return errors.New("Failed to load object store. Configuration is required")
	}
	if o == nil {
		return errors.New("Failed to load object store. Object Store record is required")
	}

	s.config = c
	s.store = o

	if s.store.Profile == "" {
		return errors.New("Failed to load object store. Profile is required")
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(s.store.Region),
		awsconfig.WithSharedConfigProfile(s.store.Profile),
	)
	if err != nil {
		logrus.Fatalf("unable to load SDK config, %v", err)
	}
	s.awsConfig = cfg

	s.client = s3.NewFromConfig(cfg)

	s.sts = sts.NewFromConfig(cfg)
	cid, err := s.sts.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		logrus.Fatalf("unable to load retrieve account ID, %v", err)
	}
	s.accountId = *cid.Account

	s.iam = iam.NewFromConfig(cfg)

	return nil
}

// CreateDataset creates a root folder for a dataset
func (s *S3Store) CreateDataset(name string, n *database.Namespace) error {

	metadata := NewMetadataFile(name)

	// Check if dataset already exists
	key := metadata.Key()
	_, err := s.client.HeadObject(context.TODO(),
		&s3.HeadObjectInput{
			Bucket: aws.String(n.BucketName),
			Key:    aws.String(key),
		})
	if err == nil {
		return errors.New("Failed to create dataset. Dataset already exists")
	}

	// TODO: DMK add additional check for a real error here vs just the key not existing?

	// Create dataset metadata file to "create" the dataset
	metaRaw, err := yaml.Marshal(metadata)
	if err != nil {
		return errors.Wrap(err, "When creating new dataset, failed to write metadata")
	}
	metaReader := bytes.NewReader(metaRaw)

	_, err = s.client.PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(n.BucketName),
			Key:    aws.String(key),
			Body:   metaReader,
		})
	if err != nil {
		return errors.Wrap(err, "When creating new dataset, failed to write metadata")
	}

	return nil
}

// DeleteDataset deletes a dataset directory
func (s *S3Store) DeleteDataset(rootDir string, n *database.Namespace) error {
	// TODO: Future work should shift this to workers as a dataset may have many
	// 	 	 objects and take a long time to delete.
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(n.BucketName),
		Prefix: aws.String(rootDir),
	}

	// Create the Paginator for the ListObjectsV2 operation.
	p := s3.NewListObjectsV2Paginator(s.client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		// Limit to 1000 results per page (this is the max)
		o.Limit = 1000
	})

	var errObjs []s3types.Error
	for p.HasMorePages() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			logrus.Fatalf("failed to get page of objects while deleting, %v", err)
		}

		// check that page has contents
		if len(page.Contents) == 0 {
			logrus.Warningf("page has no objects, skipping delete")
			continue
		}

		// convert "list objects response" to "delete objects request"
		deleteRequest := s3.DeleteObjectsInput{
			Bucket: aws.String(n.BucketName),
			Delete: &s3types.Delete{Objects: make([]s3types.ObjectIdentifier, len(page.Contents))},
		}
		for i, object := range page.Contents {
			deleteRequest.Delete.Objects[i] = s3types.ObjectIdentifier{
				Key: object.Key,
			}
		}

		output, err := s.client.DeleteObjects(context.TODO(), &deleteRequest)
		if err != nil {
			return errors.Wrap(err, "error while deleting a page of objects")
		}

		errObjs = append(errObjs, output.Errors...)
	}

	if len(errObjs) > 0 {
		var errKeys []string
		var errKeysDetail []string
		for _, o := range errObjs {
			errKeys = append(errKeys, *o.Key)
			errKeysDetail = append(errKeysDetail, fmt.Sprintf("%s (%s)", *o.Key, *o.Message))
		}
		errStr := strings.Join(errKeysDetail[:], ",")
		logrus.Errorf("Failed to delete some objects while removing the dataset at prefix `%s`: %s", rootDir, errStr)
		errStr = strings.Join(errKeys[:], ",")
		return errors.New(fmt.Sprintf("Failed to delete some objects while removing the dataset at prefix `%s`: %s", rootDir, errStr))
	}

	return nil
}

// SetUserPolicy re-renders a user's policy and applies it to the store
func (s *S3Store) SetUserPolicy(username string, permissions []*database.Permission) error {
	// Render policy
	policy, err := RenderTemplate(policyTemplateS3, readTemplateS3, readWriteTemplateS3,
		denyTemplateS3, permissions)
	if err != nil {
		return errors.Wrap(err, "Failed to update policy")
	}

	policyVersions, err := s.iam.ListPolicyVersions(
		context.TODO(),
		&iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(s.getPolicyArn(username)),
		},
	)

	if err != nil {
		if !strings.Contains(err.Error(), "api error NoSuchEntity") {
			// Some other error occured
			return errors.Wrap(err, fmt.Sprintf("Failed to list user policy verions while updating: %v", err))
		}

		// No Policy exists, so create the initial version (v1)
		description := fmt.Sprintf("HOSS policy for %s", username)
		tags := []iamtypes.Tag{
			{Key: aws.String("HOSS-Resource"), Value: aws.String("Policy")},
			{Key: aws.String("Name"), Value: aws.String(s.UserPolicyName(username))},
		}

		_, err = s.iam.CreatePolicy(context.TODO(), &iam.CreatePolicyInput{
			PolicyName:     aws.String(s.UserPolicyName(username)),
			PolicyDocument: aws.String(policy),
			Description:    &description,
			Tags:           tags,
		})
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create user policy: %v", err))
		}
	} else {
		// Get the current policy to check if it has changed
		currentPolicyVersion, err := s.iam.GetPolicyVersion(
			context.TODO(),
			&iam.GetPolicyVersionInput{
				PolicyArn: aws.String(s.getPolicyArn(username)),
				VersionId: policyVersions.Versions[0].VersionId,
			},
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to lookup existing user policy version %s for creating a new policy version: %v", *policyVersions.Versions[0].VersionId, err))
		}

		currentPolicy, err := url.QueryUnescape(*currentPolicyVersion.PolicyVersion.Document)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to decode existing user policy version %s for creating a new policy version: %v", *policyVersions.Versions[0].VersionId, err))
		}

		if trimWhitespace(policy) == trimWhitespace(currentPolicy) {
			logrus.Infof("Policy for %s is up-to-date. Skipping policy update.", username)
			return nil
		}

		// If there are already 5 versions, delete the oldest one
		// NOTE: Event though AWS says there can only be up to 5 policy versions, during testing
		//       we saw one instance of 6 policy versions existing. This code will handle this and
		//       reduce the total numer of policy verions to 4.
		if len(policyVersions.Versions) >= 5 {
			for i := 4; i < len(policyVersions.Versions); i++ {
				_, err = s.iam.DeletePolicyVersion(context.TODO(), &iam.DeletePolicyVersionInput{
					PolicyArn: aws.String(s.getPolicyArn(username)),
					VersionId: policyVersions.Versions[i].VersionId,
				})
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("Failed to delete user policy version %s for creating a new policy version: %v", *policyVersions.Versions[i].VersionId, err))
				}
			}
		}

		// Create a new version and set it as the default version
		_, err = s.iam.CreatePolicyVersion(context.TODO(), &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(s.getPolicyArn(username)),
			PolicyDocument: aws.String(policy),
			SetAsDefault:   true,
		})
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create new user policy version: %v", err))
		}
	}

	return nil
}

//DeleteUserPolicy completely removes a canned policy for a user from the system
func (s *S3Store) DeleteUserPolicy(username string) error {
	policyVersions, err := s.iam.ListPolicyVersions(
		context.TODO(),
		&iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(s.getPolicyArn(username)),
		},
	)

	if err != nil {
		if !strings.Contains(err.Error(), "api error NoSuchEntity") {
			// Some other error occured
			return errors.Wrap(err, fmt.Sprintf("Failed to list user policy versions while deleting: %v", err))
		} else {
			return nil // policy doesn't exist
		}
	}

	// Delete all non-default policy versions
	for i := range policyVersions.Versions {
		if !policyVersions.Versions[i].IsDefaultVersion {
			_, err = s.iam.DeletePolicyVersion(context.TODO(), &iam.DeletePolicyVersionInput{
				PolicyArn: aws.String(s.getPolicyArn(username)),
				VersionId: policyVersions.Versions[i].VersionId,
			})
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to delete user policy version %s: %v", *policyVersions.Versions[i].VersionId, err))
			}
		}
	}

	// Delete the default policy
	_, err = s.iam.DeletePolicy(context.TODO(), &iam.DeletePolicyInput{
		PolicyArn: aws.String(s.getPolicyArn(username)),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to delete user policy: %v", err))
	}

	return nil
}

// GetSTSCredentials gets temporary STS credentials for the object store for the current user
func (s *S3Store) GetSTSCredentials(jwt string, claims map[string]interface{}, username string) (*Credentials, error) {
	exp := int32(claims["exp"].(float64))
	durationSeconds := exp - int32(time.Now().UTC().Unix())

	// Duration minimum in STS is 15 minutes
	if durationSeconds < 15*60 {
		durationSeconds = 900
	}
	// Duration maximum is set in the role. For now we hardcode to 12 hours, assuming
	// the role has been set up per the instructions
	if durationSeconds > 12*60*60 {
		durationSeconds = 12 * 60 * 60
	}

	sessionSuffix := time.Now().Format(time.RFC3339)
	sessionSuffix = strings.ReplaceAll(sessionSuffix, ":", "")
	sessionSuffix = strings.ReplaceAll(sessionSuffix, "-", "")
	arInput := sts.AssumeRoleInput{
		RoleArn:         aws.String(s.store.RoleArn),
		RoleSessionName: aws.String(s.UserPolicyName(username) + "_" + sessionSuffix),
		DurationSeconds: aws.Int32(durationSeconds),
		PolicyArns: []ststypes.PolicyDescriptorType{
			{Arn: aws.String(s.getPolicyArn(username))},
		},
	}

	result, err := s.sts.AssumeRole(context.TODO(), &arInput)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to get temporary object store credentials: %v", err))
	}

	expiration := result.Credentials.Expiration.Format(time.RFC3339)
	creds := Credentials{AccessKeyId: *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
		Expiration:      expiration,
		Region:          s.store.Region,
		Endpoint:        s.store.Endpoint,
	}

	return &creds, nil
}

// EventsEnables checks to see if Bucket Event Notifications have been enabled for the given dataset
func (s *S3Store) EventsEnabled(namespace *database.Namespace, dataset *database.Dataset) (bool, error) {
	bucket := namespace.BucketName

	input := s3.GetBucketNotificationConfigurationInput{Bucket: &bucket}
	bucketNotification, err := s.client.GetBucketNotificationConfiguration(context.Background(), &input)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get bucket notification configuration")
	}

	for _, queueConfig := range bucketNotification.QueueConfigurations {
		for _, rule := range queueConfig.Filter.Key.FilterRules {
			if strings.ToLower(string(rule.Name)) == "prefix" && *rule.Value == dataset.RootDirectory {
				return true, nil
			}
		}
	}

	return false, nil
}

// EnableEvents turns on Bucket Event Notifications for the given dataset
func (s *S3Store) EnableEvents(namespace *database.Namespace, dataset *database.Dataset) error {
	bucket := namespace.BucketName
	queueArn := namespace.ObjectStore.NotificationArn
	if queueArn == "" {
		return fmt.Errorf("NotificationArn not defined for the '%s' object store", s.store.Name)
	}

	// get existing queue configs
	getQueueConfigInput := s3.GetBucketNotificationConfigurationInput{Bucket: &bucket}
	bucketNotification, err := s.client.GetBucketNotificationConfiguration(context.Background(), &getQueueConfigInput)
	if err != nil {
		return errors.Wrap(err, "Failed to get bucket notification configuration")
	}
	for _, queueConfig := range bucketNotification.QueueConfigurations {
		for _, rule := range queueConfig.Filter.Key.FilterRules {
			if strings.ToLower(string(rule.Name)) == "prefix" && *rule.Value == dataset.RootDirectory {
				// events already enabled for this dataset
				return nil
			}
		}
	}
	queueConfigs := append(bucketNotification.QueueConfigurations, s3types.QueueConfiguration{
		Events: []s3types.Event{
			"s3:ObjectCreated:*",
			"s3:ObjectRemoved:*",
		},
		QueueArn: &queueArn,
		Filter: &s3types.NotificationConfigurationFilter{
			Key: &s3types.S3KeyFilter{
				FilterRules: []s3types.FilterRule{
					{
						Name:  "prefix",
						Value: &dataset.RootDirectory,
					},
				},
			},
		},
	})

	// update queue configuration
	notificationConfiguration := s3types.NotificationConfiguration{
		QueueConfigurations: queueConfigs,
	}

	input := s3.PutBucketNotificationConfigurationInput{
		Bucket:                    &bucket,
		NotificationConfiguration: &notificationConfiguration,
	}
	_, err = s.client.PutBucketNotificationConfiguration(context.Background(), &input)
	if err != nil {
		return errors.Wrap(err, "Failed to get bucket notification configuration")
	}

	return nil
}

// DisableEvents turns off Bucket Notifications for the given dataset
func (s *S3Store) DisableEvents(namespace *database.Namespace, dataset *database.Dataset) error {
	bucket := namespace.BucketName

	// get existing queue configs
	getQueueConfigInput := s3.GetBucketNotificationConfigurationInput{Bucket: &bucket}
	bucketNotification, err := s.client.GetBucketNotificationConfiguration(context.Background(), &getQueueConfigInput)
	if err != nil {
		return errors.Wrap(err, "Failed to get bucket notification configuration")
	}

	indexDeleteQueueConfig := -1
	for i, queueConfig := range bucketNotification.QueueConfigurations {
		for _, rule := range queueConfig.Filter.Key.FilterRules {
			if strings.ToLower(string(rule.Name)) == "prefix" && *rule.Value == dataset.RootDirectory {
				indexDeleteQueueConfig = i
				break
			}
		}
		if indexDeleteQueueConfig != -1 {
			break
		}
	}
	if indexDeleteQueueConfig == -1 {
		logrus.Infof("Events already disabled for %s/%s", namespace.Name, dataset.Name)
		return nil
	}
	queueConfigs := append(bucketNotification.QueueConfigurations[:indexDeleteQueueConfig],
		bucketNotification.QueueConfigurations[indexDeleteQueueConfig+1:]...)

	notificationConfiguration := s3types.NotificationConfiguration{
		QueueConfigurations: queueConfigs,
	}

	input := s3.PutBucketNotificationConfigurationInput{
		Bucket:                    &bucket,
		NotificationConfiguration: &notificationConfiguration,
	}
	_, err = s.client.PutBucketNotificationConfiguration(context.Background(), &input)
	if err != nil {
		return errors.Wrap(err, "Failed to get bucket notification configuration")
	}

	return nil
}

// getPolicyArn gets the deterministic ARN for the user's policy
func (s *S3Store) getPolicyArn(username string) string {
	return fmt.Sprintf("arn:aws:iam::%s:policy/%s", s.accountId, s.UserPolicyName(username))
}
