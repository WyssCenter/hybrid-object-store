package sync

import (
	"context"
	"encoding/json"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SqsApiSyncExchange struct {
	client   *sqs.Client
	queueUrl *string
}

// sendMessage sends a message to the api-sync exchange
func (ase *SqsApiSyncExchange) SendMessage(msg *ApiEventMsg) error {

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize api sync message")
	}
	msgString := string(msgBytes)
	msgGroupId := "HOSS-Service"

	// Create a unique DeduplicationId so we don't lose any messages
	dupId := uuid.New().String()
	sendMessageInput := sqs.SendMessageInput{
		MessageBody:            &msgString,
		QueueUrl:               ase.queueUrl,
		MessageGroupId:         &msgGroupId,
		MessageDeduplicationId: &dupId,
	}

	_, err = ase.client.SendMessage(
		context.TODO(),
		&sendMessageInput,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to publish api sync message")
	}

	return nil
}

func (ase *SqsApiSyncExchange) Close() {
	// not implemented for SQS
}

func LoadSqsApiSyncExchange(queueConfig *config.SQSQueueConfig) *SqsApiSyncExchange {
	var err error
	ase := &SqsApiSyncExchange{}

	// load credentials
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(queueConfig.Region),
		awsconfig.WithSharedConfigProfile(queueConfig.Profile),
	)
	if err != nil {
		logrus.Fatalf("unable to load SDK config, %v", err)
	}

	ase.client = sqs.NewFromConfig(cfg)
	urlResult, err := ase.client.GetQueueUrl(
		context.TODO(),
		&sqs.GetQueueUrlInput{
			QueueName: &queueConfig.QueueName,
		},
	)
	if err != nil {
		logrus.Fatalf("unable to get SQS queue URL, %v", err)
	}
	ase.queueUrl = urlResult.QueueUrl

	return ase
}
