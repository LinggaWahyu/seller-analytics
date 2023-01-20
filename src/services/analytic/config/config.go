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
	NewSubscriberCfg,
)

type Config struct {
	HTTP                mhttp.HTTPServerConfig
	Database            yugabyte.YugabyteDBConfig
	RabbitMQ            messagequeue.RabbitMQConfig
	StatisticSubscriber messagequeue.SubscriberConfig
}

func NewHTTPServerCfg(cfg *Config) mhttp.HTTPServerConfig {
	return cfg.HTTP
}

func NewDatabaseCfg(cfg *Config) yugabyte.YugabyteDBConfig {
	return cfg.Database
}

func NewRabbitMQCfg(cfg *Config) messagequeue.RabbitMQConfig {
	return cfg.RabbitMQ
}

func NewSubscriberCfg(cfg *Config) messagequeue.SubscriberConfig {
	return cfg.StatisticSubscriber
}
