package amq

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
	"ohmytech.io/platform/config"
)

// Connection :
func Connection(customConfig *config.AmqConfig) (conn *amqp.Connection, chanl *amqp.Channel, err error) {
	var conf *config.AmqConfig

	if nil != customConfig {
		conf = customConfig
	} else {
		conf = &config.GetConfig().Amq
	}

	conn, err = amqp.Dial(conf.Connector)
	if err != nil {
		log.Printf("%s: %s", "Failed to connect to RabbitMQ", err)
		return nil, nil, errors.New("Failed to connect to RabbitMQ")
	}

	chanl, err = conn.Channel()
	if err != nil {
		log.Printf("%s: %s", "Failed to open a channel", err)
		return nil, nil, errors.New("Failed to open a channel")
	}

	return conn, chanl, nil
}
