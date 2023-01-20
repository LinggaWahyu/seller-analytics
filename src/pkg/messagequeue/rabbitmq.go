package messagequeue

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// NewRabbitMQ, establishes a rabbitMQ connection and returns it
func NewRabbitMQ(cfg RabbitMQConfig) (*amqp091.Connection, error) {
	log.Println("Initialize rabbitMQ connection...")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	log.Println("Finish rabbitMQ connection")

	return conn, nil
}
