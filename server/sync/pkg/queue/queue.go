package queue

import (
	"github.com/pkg/errors"

	"github.com/gigantum/hoss-sync/pkg/config"
)

// Queue provides a generic interface for a message queue implementation
type Queue interface {
	// Send gets the channel used to send messages to the queue
	Send() chan<- config.Message

	// Receive gets the channel used to receive messages from the queue
	Receive() <-chan config.Message
}

// LoadNotificationQueue loads the specific queue implementation that receives notification
// messages from an object store or core service
func LoadNotificationQueue(configuration *config.Configuration, queueConfig *config.NotificationQueueConfig) (Queue, error) {
	switch queueConfig.Type {
	case "amqp":
		var queueSettings config.AMQPQueueConfig
		if err := config.UnmarshalSettings(queueConfig.Settings, &queueSettings); err != nil {
			return nil, errors.Wrap(err, "Could not load AMQP queue settings")
		}
		return AMQPNotifications(&queueSettings), nil
	case "sqs":
		var queueSettings config.SQSQueueConfig
		if err := config.UnmarshalSettings(queueConfig.Settings, &queueSettings); err != nil {
			return nil, errors.Wrap(err, "Could not load SQS queue settings")
		}
		queueSettings.Profile = configuration.SqsProfile
		return SQSNotifications(&queueSettings), nil
	default:
		return nil, errors.New("Notification queue type not supported")
	}
}
