package sync

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type AmqpApiSyncExchange struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// sendMessage sends a message to the api-sync exchange
func (ase *AmqpApiSyncExchange) SendMessage(msg *ApiEventMsg) error {

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize api sync message")
	}

	err = ase.channel.Publish(
		"hoss",              // exchange
		"api_notifications", // routing key
		true,                // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBytes,
		})

	if err != nil {
		return errors.Wrap(err, "Failed to publish api sync message")
	}

	return nil
}

func (ase *AmqpApiSyncExchange) Close() {
	ase.channel.Close()
	ase.connection.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func LoadAmqpApiSyncExchange(queueConfig *config.AMQPQueueConfig) *AmqpApiSyncExchange {
	var err error
	ase := &AmqpApiSyncExchange{}

	for i := 0; i < 5; i++ {
		ase.connection, err = amqp.Dial(queueConfig.Url)
		if err == nil {
			break
		}
		log.Printf("Error dialing RabbitMQ, trying again: %s", err.Error())
		time.Sleep(5 * time.Second)
	}
	failOnError(err, "Failed to connect to RabbitMQ when creating API Sync Exchange")

	ase.channel, err = ase.connection.Channel()
	failOnError(err, "Failed to open a channel when creating API Sync Exchange")

	err = ase.channel.ExchangeDeclare(
		"hoss",   // name
		"direct", // kind
		true,     // durable
		false,    // auto delete
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare exchangewhen creating API Sync Exchange")

	return ase
}
