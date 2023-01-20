package config

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/cfg/viper"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	mhttp "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	viper.InitDefaultConfig[Config],
	NewHTTPServerCfg,
	NewDatabaseCfg,
	NewRabbitMQCfg,
	NewPublisherCfg,
)

type Config struct {
	HTTP           mhttp.HTTPServerConfig
	Database       yugabyte.YugabyteDBConfig
	RabbitMQ       messagequeue.RabbitMQConfig
	OrderPublisher messagequeue.PublisherConfig
}

// NewHTTPServerCfg, provides http config to dependency injection
func NewHTTPServerCfg(cfg *Config) mhttp.HTTPServerConfig {
	return cfg.HTTP
}

// NewDatabaseCfg, provides database config to dependency injection
func NewDatabaseCfg(cfg *Config) yugabyte.YugabyteDBConfig {
	return cfg.Database
}

// NewRabbitMQCfg, provides rabbitmq config to dependency injection
func NewRabbitMQCfg(cfg *Config) messagequeue.RabbitMQConfig {
	return cfg.RabbitMQ
}

// NewPublisherCfg, provides mq publisher config to dependency injection
func NewPublisherCfg(cfg *Config) messagequeue.PublisherConfig {
	return cfg.OrderPublisher
}
