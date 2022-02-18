package sync

import (
	"os"
	"strings"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ApiSyncExchange provides a generic interface for a sync exchange implementation
type ApiSyncExchange interface {
	// SendMessage sends an API notification message to the queue
	SendMessage(*ApiEventMsg) error

	// Used for testing purposes
	Close()
}

// LoadApiSyncExchange loads the API sync exchange implementation that sends messages
// to a specific type of queue
func LoadApiSyncExchange(cfg *config.Configuration) (map[string]ApiSyncExchange, error) {
	ase := make(map[string]ApiSyncExchange)
	for _, queueConfig := range cfg.Queues {
		switch queueConfig.Type {
		case "amqp":
			var queueSettings config.AMQPQueueConfig
			if err := config.UnmarshalSettings(queueConfig.Settings, &queueSettings); err != nil {
				return nil, errors.Wrap(err, "Could not load AMQP queue settings")
			}
			queueSettings.Url = strings.Replace(queueSettings.Url, "${RABBITMQ_USER}", os.Getenv("RABBITMQ_USER"), -1)
			queueSettings.Url = strings.Replace(queueSettings.Url, "${RABBITMQ_PASS}", os.Getenv("RABBITMQ_PASS"), -1)

			ase[queueConfig.ObjectStore] = LoadAmqpApiSyncExchange(&queueSettings)
		case "sqs":
			var queueSettings config.SQSQueueConfig
			if err := config.UnmarshalSettings(queueConfig.Settings, &queueSettings); err != nil {
				return nil, errors.Wrap(err, "Could not load SQS queue settings")
			}
			ase[queueConfig.ObjectStore] = LoadSqsApiSyncExchange(&queueSettings)
		default:
			logrus.Error("Notification queue type not supported")
		}
	}
	return ase, nil
}
