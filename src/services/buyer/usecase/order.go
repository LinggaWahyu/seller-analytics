package usecase

import (
	"context"
	"errors"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository"
)

type OrderUsecase interface {
	Products(ctx context.Context) ([]domain.Product, error)
	ProductByID(ctx context.Context, id uint) (*domain.Product, error)
	UpdateOrderStatus(ctx context.Context, orderId uint, status string) (*domain.Order, error)
	CreateOrder(ctx context.Context, req domain.Order) (*domain.Order, error)
	OrderByID(ctx context.Context, id uint) (*domain.Order, error)
	OrdersByBuyer(ctx context.Context, buyerId uint) ([]domain.Order, error)
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

// Products are method to return list of products
func (ou *orderUsecase) Products(ctx context.Context) ([]domain.Product, error) {
	return nil, errors.New("unimplemented")
}

// ProductByID are method to return a product
func (ou *orderUsecase) ProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	return nil, errors.New("unimplemented")
}

// UpdateOrder is an update method for order
func (ou *orderUsecase) UpdateOrderStatus(ctx context.Context, orderId uint, status string) (*domain.Order, error) {
	return nil, errors.New("unimplemented")
}

// CreateOrder is an update method for order
func (ou *orderUsecase) CreateOrder(ctx context.Context, req domain.Order) (*domain.Order, error) {
	return nil, errors.New("unimplemented")
}

// OrderByID are method to return list of orders
func (ou *orderUsecase) OrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	return nil, errors.New("unimplemented")
}

// Orders are method to return list of orders
func (ou *orderUsecase) OrdersByBuyer(ctx context.Context, buyerId uint) ([]domain.Order, error) {
	return nil, errors.New("unimplemented")
}
