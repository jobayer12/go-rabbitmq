package main

import (
	"context"
	"eventdrivenrabbit/internal"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Println("New message", msg)
				time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Println("Ack message failed")
					return err
				}
				log.Println("Acknowledge message", msg.MessageId)
				return nil
			})
		}
	}()

	//go func() {
	//	for message := range messageBus {
	//		log.Println("New Message", message)
	//
	//		if !message.Redelivered {
	//			message.Nack(false, true)
	//			continue
	//		}
	//
	//		if err := message.Ack(false); err != nil {
	//			log.Println("Failed to ack message")
	//			continue
	//		}
	//
	//		//if err := message.Ack(false); err != nil {
	//		//	log.Println("Acknowledge message failed")
	//		//	continue
	//		//}
	//		log.Printf("Acknowledge message %s\n", message.MessageId)
	//
	//	}
	//}()
	log.Println("Consuming, use CTRL+C to exit")
	<-blocking
}
