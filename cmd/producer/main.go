package main

import (
	"context"
	"eventdrivenrabbit/internal"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
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

	// all consuming will be done this connection
	consumeConn, err := internal.ConnectRabbitMQ("jobayer", "jobayer", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer consumeConn.Close()

	consumeClient, err := internal.NewRabbitMQClient(consumeConn)

	if err != nil {
		panic(err)
	}
	defer consumeClient.Close()

	//if err := client.CreateQueue("customers_created", true, false); err != nil {
	//	panic(err)
	//}
	//
	//if err := client.CreateQueue("customers_test", false, true); err != nil {
	//	panic(err)
	//}
	//
	//if err := client.CreateBinding("customers_created", "customers.created.*", "customer_events"); err != nil {
	//	panic(err)
	//}
	//
	//if err := client.CreateBinding("customers_test", "customers.*", "customer_events"); err != nil {
	//	panic(err)
	//}

	queue, err := consumeClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err := consumeClient.CreateBinding(queue.Name, queue.Name, "customer_callbacks"); err != nil {
		panic(err)
	}

	messageBus, err := consumeClient.Consume(queue.Name, "customer-api", false)
	if err != nil {
		panic(err)
	}
	go func() {
		for message := range messageBus {
			log.Printf("Message callback %s\n", message.CorrelationId)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	for i := 0; i < 100; i++ {
		if err := client.Send(ctx, "customer_events", "customers.created.us", amqp.Publishing{
			ContentType:   "plain/text",
			DeliveryMode:  amqp.Persistent,
			ReplyTo:       queue.Name,
			CorrelationId: fmt.Sprintf("customer_created_%d", i),
			Body:          []byte(`A cool message`),
		}); err != nil {
			panic(err)
		}
	}

	time.Sleep(10 * time.Second)
	log.Println(client)
	var blocking chan struct{}
	<-blocking
}
