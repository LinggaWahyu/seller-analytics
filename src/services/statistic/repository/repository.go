package repository

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(yugabyte.NewDatabase),
	fx.Provide(messagequeue.NewRabbitMQ),
	fx.Provide(messagequeue.NewRabbitMQSubscriber[domain.PayloadEventOrder]),
	fx.Provide(messagequeue.NewRabbitMQPublisher[domain.PayloadEventStatistic]),
	fx.Provide(NewStatisticsRepository),
	fx.Invoke(AutoMigrateEntities),
)

func AutoMigrateEntities(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Statistics{}); err != nil {
		return err
	}
	return nil
}
