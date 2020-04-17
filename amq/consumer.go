package amq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"ohmytech.io/platform/config"
)

// Consume :
func Consume(queueName string, customConfig *config.AmqConfig) (<-chan amqp.Delivery, error) {
	_, ch, err := Connection(customConfig)
	if err != nil {
		log.Printf("%s: %s", "Connection to RabbitMQ", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("%s: %s", "Failed to declare a queue", err)
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("%s: %s", "Failed to register a consumer", err)
		return nil, err
	}

	fmt.Println("Consume > ", msgs)

	return msgs, nil
}
