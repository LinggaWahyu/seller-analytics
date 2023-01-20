package messagequeue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// Publisher, interface to publish to an mq topic, T is the message format to be published
type Publisher[T any] interface {
	Publish(ctx context.Context, publish PublishConfig, message T) error
}

// rabbitMQPublisher, concrete implementation of Publisher publishing to rabbitMQ exchange
type rabbitMQPublisher[T any] struct {
	exchange     string
	rabbitMQConn *amqp091.Connection
}

// NewRabbitMQPublisher, constructor returning rabbitMQPublisher as Publisher
func NewRabbitMQPublisher[T any](config PublisherConfig, conn *amqp091.Connection) Publisher[T] {
	repo := &rabbitMQPublisher[T]{
		exchange:     config.Exchange.Name,
		rabbitMQConn: conn,
	}

	if err := repo.initPublisher(config); err != nil {
		log.Fatal(err)
	}

	return repo
}

// initPublisher, performs declarations for rabbitmq AMQP model necessary for publisher
func (repo *rabbitMQPublisher[T]) initPublisher(config PublisherConfig) error {
	log.Println("Initialize rabbitMQ architecture...")

	ch, err := repo.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declaring exchange
	xchg := config.Exchange
	err = ch.ExchangeDeclare(xchg.Name, xchg.Kind, xchg.Durable, xchg.AutoDelete, xchg.Internal, xchg.NoWait, xchg.Args)
	if err != nil {
		return err
	}

	log.Println("Finish rabbitMQ publisher setup")

	return nil
}

// Publish, allows publishing to designated rabbitmq exchanges for messages of type T
func (repo *rabbitMQPublisher[T]) Publish(ctx context.Context, publish PublishConfig, message T) error {
	ch, err := repo.rabbitMQConn.Channel()
	if err != nil {
		return err
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Note we only support json messages for now
	err = ch.PublishWithContext(ctx, repo.exchange, publish.Key, publish.Mandatory, publish.Immediate, amqp091.Publishing{
		ContentType: "application/json",
		Body:        messageBytes,
	})
	if err != nil {
		return err
	}

	return nil
}
