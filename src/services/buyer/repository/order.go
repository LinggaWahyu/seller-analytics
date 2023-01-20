package repository

import (
	"context"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetProducts(ctx context.Context) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id uint) (*domain.Product, error)
	GetOrdersByBuyerID(ctx context.Context, buyerId uint) ([]domain.Order, error)
	UpdateOrderById(ctx context.Context, order domain.Order) (*domain.Order, error)
	InsertOrder(ctx context.Context, order domain.Order) (*domain.Order, error)
	GetOngoingOrders(ctx context.Context, buyerId uint) (bool, error)
	PublishOrderEvent(ctx context.Context, event domain.PayloadEventOrder) error
	GetOrderByID(ctx context.Context, id uint) (*domain.Order, error)
}

type orderRepository struct {
	db               *gorm.DB
	repoCoreRabbitMQ messagequeue.Publisher[domain.PayloadEventOrder]
}

func NewOrderRepository(db *gorm.DB, repoCoreRabbitMQ messagequeue.Publisher[domain.PayloadEventOrder]) OrderRepository {
	return &orderRepository{
		db:               db,
		repoCoreRabbitMQ: repoCoreRabbitMQ,
	}
}

// GetProducts
func (or *orderRepository) GetProducts(ctx context.Context) ([]domain.Product, error) {
	var res []domain.Product

	if err := or.db.Find(&res).Error; err != nil && err != gorm.ErrRecordNotFound {
		return res, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return res, nil
}

// GetProductByID
func (or *orderRepository) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	var res domain.Product

	query := or.db.WithContext(ctx)
	if err := query.Where("id = ?", id).First(&res).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &res, nil
}

// GetOrders
func (or *orderRepository) GetOrdersByBuyerID(ctx context.Context, buyerId uint) ([]domain.Order, error) {
	var res []domain.Order

	query := or.db.WithContext(ctx)
	if err := query.Where("buyer_id = ?", buyerId).Preload("OrderDetails").Preload("OrderDetails.Product").Find(&res).Error; err != nil && err != gorm.ErrRecordNotFound {
		return res, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	for i := range res {
		res[i].OrderDateStr = time.Time(res[i].OrderDate).Format(domain.OrderDateFormat)
	}

	return res, nil
}

func (or *orderRepository) UpdateOrderById(ctx context.Context, order domain.Order) (*domain.Order, error) {
	if err := or.db.Model(&order).UpdateColumns(domain.Order{Status: order.Status}).Error; err != nil {
		return nil, err
	}
	order.OrderDateStr = time.Time(order.OrderDate).Format(domain.OrderDateFormat)
	return &order, nil
}

// InsertOrder
func (or *orderRepository) InsertOrder(ctx context.Context, order domain.Order) (*domain.Order, error) {
	if err := or.db.Create(&order).Error; err != nil {
		return nil, err
	}
	order.OrderDateStr = time.Time(order.OrderDate).Format(domain.OrderDateFormat)
	return &order, nil
}

// GetOngoingOrders
func (or *orderRepository) GetOngoingOrders(ctx context.Context, buyerId uint) (bool, error) {
	var res []domain.Order

	query := or.db.WithContext(ctx)
	if err := query.Where("status = ? AND buyer_id = ?", domain.OrderStatusNew, buyerId).Find(&res).Error; err != nil {
		return false, err
	}

	if len(res) > 0 {
		return true, nil
	}

	for i := range res {
		res[i].OrderDateStr = time.Time(res[i].OrderDate).Format(domain.OrderDateFormat)
	}

	return false, nil
}

// PublishOrderEvent
func (or *orderRepository) PublishOrderEvent(ctx context.Context, event domain.PayloadEventOrder) error {
	err := or.repoCoreRabbitMQ.Publish(ctx, messagequeue.PublishConfig{}, event)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByID
func (or *orderRepository) GetOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	var res domain.Order

	query := or.db.WithContext(ctx)
	if err := query.Where("id = ?", id).Preload("OrderDetails").Preload("OrderDetails.Product").First(&res).Error; err != nil {
		return nil, err
	}
	res.OrderDateStr = time.Time(res.OrderDate).Format(domain.OrderDateFormat)

	return &res, nil
}
