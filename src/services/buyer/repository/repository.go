package repository

import (
	"strconv"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(yugabyte.NewDatabase),
	fx.Provide(messagequeue.NewRabbitMQ),
	fx.Provide(messagequeue.NewRabbitMQPublisher[domain.PayloadEventOrder]),
	fx.Provide(NewBuyerRepository),
	fx.Provide(NewOrderRepository),
	fx.Invoke(AutoMigrateEntities),
	fx.Invoke(PrepareProductData),
)

// AutoMigrateEntities, auto migrate database schema from domain models to database
func AutoMigrateEntities(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Buyer{}, &domain.Order{}, &domain.OrderDetail{}, &domain.Product{}); err != nil {
		return err
	}
	return nil
}

func PrepareProductData(db *gorm.DB) error {
	//var products []domain.Product

	for i := 0; i < 2; i++ {
		product := domain.Product{
			Price:       100,
			ProductName: "Product " + strconv.Itoa(i+1),
		}

		//products = append(products, product)
		if err := db.Where(product).FirstOrCreate(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
