package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	errors "github.com/gigantum/hoss-error"
	service "github.com/gigantum/hoss-service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// STSCredentials is the response from the Core Service that contains the information about the generated STS credentials
type STSCredentials struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	Expiration      string `json:"expiration"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
}

// getCredentials requests STS credentials for the given namespace from the given core service
func getCredentials(idToken string, core_service string, objstore string) (*STSCredentials, error) {
	// Hack to support running on localhost
	core_service = strings.Replace(core_service, "localhost/core", "core:8080", 1)

	// Get the target namespace information
	req, err := http.NewRequest("GET", core_service+"/object_store/"+objstore+"/sts", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, errors.New("problem with STS response: " + err.Error())
		}

		logrus.Debug(string(d))
		return nil, errors.New("problem with STS response: StatusCode != 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var creds STSCredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// getClient initializes an S3 client based on the given STS credentials
func getClient(creds *STSCredentials) (*s3.Client, error) {
	var ops []func(*config.LoadOptions) error
	ops = append(ops,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				creds.AccessKeyId,
				creds.SecretAccessKey,
				creds.SessionToken,
			),
		),
	)

	if creds.Region != "" {
		ops = append(ops, config.WithRegion(creds.Region))
	}

	// NOTE: let the sdk determine the correct endpoint when communicating with AWS
	if creds.Endpoint != "" && !strings.Contains(creds.Endpoint, "amazonaws.com") {
		// Hack to support running on localhost
		endpoint := creds.Endpoint
		if strings.HasSuffix(endpoint, "://localhost") {
			endpoint = "http://minio:9000"
		}

		ops = append(ops, config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						PartitionID:       "aws",
						URL:               endpoint,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			}),
		))
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		ops...,
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

// RenewingClient is an interface for an object that can load a S3 Client and keep it updated when credentials are about to expire
type RenewingClient interface {
	// Get the S3 Client and any error that occurred when the last updated happened
	GetClient() (*s3.Client, error)

	// Manually request that the client be refreshed, so any policy changes are used
	ForceRefresh()

	// Background routine that will update the client when credentials are about to expire
	RefreshRoutine(ctx context.Context)
}

// RenewingClientImpl implements the RenewingClient interface
type RenewingClientImpl struct {
	mu sync.RWMutex

	sts       *STSCredentials
	client    *s3.Client
	lastError error

	tokens      service.RenewingTokens
	coreService string
	objectStore string
	interval    time.Duration
}

// GetClient gets the S3 client and any error that occurred during the last update
func (i *RenewingClientImpl) GetClient() (*s3.Client, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.lastError != nil {
		return nil, i.lastError
	}

	if i.client == nil {
		return nil, errors.New("No S3 client is available")
	}

	return i.client, nil
}

// ForceRefresh fetches new credentials and creates a new S3 client
func (i *RenewingClientImpl) ForceRefresh() {
	i.updateCredentials() // sets the internal fields
}

// updateCredentials is a private method that is used to set new credentials
// New STS credentials are queried and a new S3 client is created
// Finally the write lock is acquired and the struct fields are updated
// Any errors are set in the RenewingClientImpl.lastError field
func (i *RenewingClientImpl) updateCredentials() {
	// common method to acquire the write lock and set the fields, also print a log message
	setFields := func(sts *STSCredentials, client *s3.Client, lastError error) {
		i.mu.Lock()
		i.sts, i.client, i.lastError = sts, client, lastError
		if i.lastError != nil {
			logrus.Infof("Failed to fefresh %s object store client: %s", i.objectStore, i.lastError.Error())
		} else {
			logrus.Infof("Succesfully refreshed %s object store client", i.objectStore)
		}
		i.mu.Unlock()
	}

	idToken, err := i.tokens.GetIDToken()
	if err != nil {
		err = fmt.Errorf("could not get service ID token for authentication: %w", err)
		setFields(nil, nil, err)
		return
	}

	sts, err := getCredentials(idToken, i.coreService, i.objectStore)
	if err != nil {
		err = fmt.Errorf("could not renew STS credentials: %w", err)
		setFields(nil, nil, err)
		return
	}

	client, err := getClient(sts)
	if err != nil {
		err = fmt.Errorf("could not get S3 client: %w", err)
		setFields(sts, nil, err)
		return
	}

	setFields(sts, client, nil)
}

// RefreshRoutine is a never ending method that will periodically refresh the STS credentials and S3 client
func (i *RenewingClientImpl) RefreshRoutine(ctx context.Context) {
	ticker := time.NewTicker(i.interval)
	logrus.Infof("Starting %s object store client refresh routine. Interval: %s", i.objectStore, i.interval)

	for {
		select {
		case <-ticker.C:
			logrus.Infof("Renewing STS Credentials for %s", i.objectStore)
			i.updateCredentials() // acquired write lock
		case <-ctx.Done():
			logrus.Infof("Stopping %s object store client refresh routine", i.objectStore)
			ticker.Stop()
			return
		}
	}
}

// GetRenewingClient returns a RewnewingClientImpl that will perodically refresh
// the stored STS credentials and S3 client.
// Note: it is up to the caller to start the refresh routine
func GetRenewingClient(tokens service.RenewingTokens, coreService string, objectStore string, interval time.Duration) RenewingClient {
	impl := &RenewingClientImpl{
		tokens:      tokens,
		coreService: coreService,
		objectStore: objectStore,
		interval:    interval,
	}

	impl.updateCredentials() // sets the internal fields

	return impl
}
