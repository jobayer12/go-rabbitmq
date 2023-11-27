package main

import (
	"eventdrivenrabbit/internal"
	"log"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("jobayer", "jobayer", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}

	go func() {
		for message := range messageBus {
			log.Println("new Message", message)
			if err := message.Ack(false); err != nil {
				log.Println("Acknowledge message failed")
				continue
			}
			log.Printf("Acknowledge message %s\n", message.MessageId)

		}
	}()
	<-blocking
}
