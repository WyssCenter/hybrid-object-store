package queue

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gigantum/hoss-sync/pkg/config"
	"github.com/gigantum/hoss-sync/pkg/message"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// AMQPQueue defines a RabbitMQ backed notification queue
type AMQPQueue struct {
	// The settings for the queue
	queueConfig  *config.AMQPQueueConfig
	queueName    string
	exchangeName string

	// The queue references
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery

	// The channel that is used for the Queue interface
	decodedMsgs chan config.Message

	// The type of message into which data should be unmarshaled
	messageType string
}

// Send is not implemented for a notification queue
func (q *AMQPQueue) Send() chan<- config.Message {
	logrus.Fatal("AMQP queue sending not enabled")
	return nil
}

// Receive returns the channel containing decoded notification messages from the queue
func (q *AMQPQueue) Receive() <-chan config.Message {
	return q.decodedMsgs
}

func AMQPNotifications(queueConfig *config.AMQPQueueConfig) Queue {
	var err error
	q := &AMQPQueue{
		queueConfig:  queueConfig,
		exchangeName: queueConfig.ExchangeName,
		queueName:    queueConfig.QueueName,
		messageType:  queueConfig.MessageType,
		decodedMsgs:  make(chan config.Message),
	}
	for i := 0; i < 5; i++ {
		q.conn, err = amqp.Dial(os.Getenv("AMQP_URL"))
		if err == nil {
			break
		}
		log.Printf("Error dialing RabbitMQ, trying again: %s", err.Error())
		time.Sleep(5 * time.Second)
	}
	failOnError(err, "Failed to connect to RabbitMQ")

	// Create the channel
	q.channel, err = q.conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Create the queue
	_, err = q.channel.QueueDeclare(
		q.queueName, // name
		true,        // durable
		true,        // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create the exchange
	err = q.channel.ExchangeDeclare(
		q.exchangeName, // name
		"direct",       // kind
		true,           // durable
		false,          // auto delete
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare exchange")

	// Connect the queue to the exchange
	err = q.channel.QueueBind(
		q.queueName, // name
		q.queueName, // key
		"hoss",      // exchange
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to bind queue")

	// Start the consumer reading from the queue
	q.msgs, err = q.channel.Consume(
		q.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	// Goroutine to decode incoming messages into the internal format
	go func() {
		for {
			data, more := <-q.msgs
			if !more {
				logrus.Error("AMQP Notification queue broken")
				return
			}
			if q.messageType == "bucket_notification" {
				var note message.BucketNotification
				err := json.Unmarshal(data.Body, &note)
				if err != nil {
					logrus.Error("Problem decoding message: " + err.Error())
				} else {
					for _, record := range note.Records {
						record.Endpoint = queueConfig.SourceEndpoint
						q.decodedMsgs <- &record
					}
				}
			} else if q.messageType == "api_notification" {
				var msg message.ApiSyncNotification
				err := json.Unmarshal(data.Body, &msg)
				if err != nil {
					logrus.Error("Problem decoding message: " + err.Error())
				} else {
					q.decodedMsgs <- &msg
				}
			} else {
				logrus.Error("Unsupported message type set: " + q.messageType)
			}
		}
	}()

	return q
}
