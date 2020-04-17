package amq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"ohmytech.io/platform/config"
)

// Publish :
func Publish(queueName string, obj interface{}, customConfig *config.AmqConfig) error {
	_, ch, err := Connection(customConfig)
	if err != nil {
		log.Printf("%s: %s", "Connection to RabbitMQ", err)
		return err
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
		return err
	}

	body, err := json.Marshal(obj)
	if err != nil {
		log.Printf("%s: %s", "Failed to json.Marshal", err)
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("%s: %s", "Failed to publish a message", err)
		return err
	}

	return nil
}
