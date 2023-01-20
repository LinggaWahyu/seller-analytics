package repository

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
	statdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(yugabyte.NewDatabase),
	fx.Provide(messagequeue.NewRabbitMQ),
	fx.Provide(messagequeue.NewRabbitMQSubscriber[statdomain.PayloadEventStatistic]),
	fx.Provide(NewAnalyticRepository),
	fx.Invoke(AutoMigrateEntities),
)

func AutoMigrateEntities(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Analytic{}); err != nil {
		return err
	}
	return nil
}
