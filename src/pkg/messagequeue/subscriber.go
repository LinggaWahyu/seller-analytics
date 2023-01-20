package messagequeue

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

// Subscriber, interface to subscribe to an mq topic, T is the message format to be received
type Subscriber[T any] interface {
	Subscribe(subscribe SubscribeConfig, handlerFunc func(msg T)) error
}

// rabbitMQSubscriber, concrete implementation of Subscriber subscribing to rabbitMQ queue
type rabbitMQSubscriber[T any] struct {
	queue        string
	rabbitMQConn *amqp091.Connection
}

// NewRabbitMQSubscriber, constructor returning rabbitMQSubscriber as Subscriber
func NewRabbitMQSubscriber[T any](config SubscriberConfig, conn *amqp091.Connection) Subscriber[T] {
	repo := &rabbitMQSubscriber[T]{
		queue:        config.Queue.Name,
		rabbitMQConn: conn,
	}

	if err := repo.initSubscriber(config); err != nil {
		log.Fatal(err)
	}

	return repo
}

// initSubscriber, performs declarations for rabbitmq AMQP model necessary for subscriber
func (repo *rabbitMQSubscriber[T]) initSubscriber(config SubscriberConfig) error {
	log.Println("Initialize rabbitMQ architecture...")

	ch, err := repo.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declaring exchange
	err = ch.ExchangeDeclare(config.Exchange.Name, config.Exchange.Kind, config.Exchange.Durable, config.Exchange.AutoDelete, config.Exchange.Internal, config.Exchange.NoWait, config.Exchange.Args)
	if err != nil {
		return err
	}

	// declaring queue
	_, err = ch.QueueDeclare(config.Queue.Name, config.Queue.Durable, config.Queue.AutoDelete, config.Queue.Exclusive, config.Queue.NoWait, config.Queue.Args)
	if err != nil {
		return err
	}

	// bind exchanges and config.Queues
	err = ch.QueueBind(config.Binding.Name, config.Binding.Key, config.Binding.Exchange, config.Binding.NoWait, config.Binding.Args)
	if err != nil {
		return err
	}

	log.Println("Finish rabbitMQ architecture")

	return nil
}

// Subscribe, allows subscribing to designated rabbitmq queues for messages of type T
func (repo *rabbitMQSubscriber[T]) Subscribe(subscribe SubscribeConfig, handlerFunc func(msg T)) error {
	ch, err := repo.rabbitMQConn.Channel()
	if err != nil {
		return err
	}

	msgChan, err := ch.Consume(repo.queue, subscribe.Consumer, subscribe.AutoAck, subscribe.Exclusive, subscribe.NoLocal, subscribe.NoWait, subscribe.Args)
	if err != nil {
		return err
	}

	for msg := range msgChan {
		var event T

		if err = json.Unmarshal(msg.Body, &event); err != nil {
			log.Println(errors.Wrapf(err, "error unmarshall body"))
			continue
		}

		handlerFunc(event)
	}

	return nil
}
