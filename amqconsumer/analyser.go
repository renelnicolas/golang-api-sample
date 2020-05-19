package main

import (
	"log"

	"ohmytech.io/platform/amq"
	"ohmytech.io/platform/config"
)

// go build -o bin/parsing parsing.go && bin/parsing
const queueName = "analyser"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conf := config.New("../.env")

	msgs, err := amq.Consume(queueName, &conf.Amq)
	failOnError(err, "Failed to open a channel")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// TODO
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
