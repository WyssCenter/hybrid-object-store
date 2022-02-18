package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gigantum/hoss-sync/pkg/config"
	"github.com/gigantum/hoss-sync/pkg/message"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
)

// SQSQueue defines an AWS notification queue
type SQSQueue struct {
	// The settings for the queue
	queueConfig *config.SQSQueueConfig
	queueName   string

	// The channel that is used for the Queue interface
	decodedMsgs chan config.Message

	// The type of message into which data should be unmarshaled
	messageType string
}

// Send is not implemented for a notification queue
func (q *SQSQueue) Send() chan<- config.Message {
	logrus.Fatal("SQS queue sending not enabled")
	return nil
}

// Receive returns the channel containing decoded notification messages from the queue
func (q *SQSQueue) Receive() <-chan config.Message {
	return q.decodedMsgs
}

func SQSNotifications(queueConfig *config.SQSQueueConfig) Queue {
	var err error
	q := &SQSQueue{
		queueConfig: queueConfig,
		queueName:   queueConfig.QueueName,
		messageType: queueConfig.MessageType,
		decodedMsgs: make(chan config.Message),
	}

	// load credentials
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(queueConfig.Region),
		awsconfig.WithSharedConfigProfile(queueConfig.Profile),
	)
	if err != nil {
		logrus.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	urlResult, err := client.GetQueueUrl(
		context.TODO(),
		&sqs.GetQueueUrlInput{
			QueueName: &q.queueName,
		},
	)
	if err != nil {
		logrus.Fatalf("unable to get SQS queue URL, %v", err)
	}
	queueURL := urlResult.QueueUrl

	// Goroutine to decode incoming messages into the internal format
	go func() {
		for {
			receiveMessageInput := sqs.ReceiveMessageInput{
				AttributeNames: []sqstypes.QueueAttributeName{
					sqstypes.QueueAttributeNameAll,
				},
				MessageAttributeNames: []string{
					"All",
				},
				QueueUrl:            queueURL,
				MaxNumberOfMessages: int32(1),
				VisibilityTimeout:   int32(30),
				WaitTimeSeconds:     int32(5),
			}

			msgResult, err := client.ReceiveMessage(
				context.TODO(),
				&receiveMessageInput,
			)
			if err != nil {
				logrus.Warningf("unable to get SQS message, %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if len(msgResult.Messages) > 0 {
				if q.messageType == "bucket_notification" {
					var note message.BucketNotification
					err := json.Unmarshal([]byte(*msgResult.Messages[0].Body), &note)
					if err != nil {
						logrus.Error("Problem decoding message: " + err.Error())
					} else {
						for _, record := range note.Records {
							record.Endpoint = queueConfig.SourceEndpoint
							q.decodedMsgs <- &record
						}
					}
				} else if q.messageType == "api_notification" {
					var notification message.ApiSyncNotification
					err := json.Unmarshal([]byte(*msgResult.Messages[0].Body), &notification)
					if err != nil {
						logrus.Error("Problem decoding message: " + err.Error())
					} else {
						q.decodedMsgs <- &notification
					}
				} else {
					logrus.Error("Unsupported message type set: " + q.messageType)
					continue // skip the message delete
				}

				_, err = client.DeleteMessage(
					context.TODO(),
					&sqs.DeleteMessageInput{
						QueueUrl:      queueURL,
						ReceiptHandle: msgResult.Messages[0].ReceiptHandle,
					},
				)
				if err != nil {
					logrus.Warning("Could not delete processed message: " + err.Error())
				}
			}
		}
	}()

	return q
}
