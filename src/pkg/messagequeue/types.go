package messagequeue

import "github.com/rabbitmq/amqp091-go"

// RabbitMQConfig, config to establish a RabbitMQ connection
type RabbitMQConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

// SubscriberConfig, config to setup a subscriber instance
type SubscriberConfig struct {
	Exchange AMQPExchangeConfig
	Queue    AMQPQueueConfig
	Binding  AMQPBindConfig
}

// PublisherConfig, config to setup a publisher instance
type PublisherConfig struct {
	Exchange AMQPExchangeConfig
}

// AMQPExchangeConfig, config to declare an exchange in rabbitmq
type AMQPExchangeConfig struct {
	Name string
	Kind string

	NoWait     bool
	Durable    bool
	AutoDelete bool
	Internal   bool

	Args amqp091.Table
}

// AMQPQueueConfig, config to declare a queue in rabbitmq
type AMQPQueueConfig struct {
	Name       string
	NoWait     bool
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	Args       amqp091.Table
}

// AMQPBindConfig, config to declare a binding in rabbitmq
type AMQPBindConfig struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool

	Args amqp091.Table
}

// PublishConfig, config to determine how to publish to a topic via Publisher
type PublishConfig struct {
	Key       string
	Mandatory bool
	Immediate bool
}

// SubscribeConfig, config to determine how to subscribe to a topic via Subscriber
type SubscribeConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp091.Table
}
