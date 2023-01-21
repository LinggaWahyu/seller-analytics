package usecase

import (
	"context"
	"errors"
	"log"
	"time"

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
	var res []domain.Product

	res, err := ou.orderRepo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ProductByID are method to return a product
func (ou *orderUsecase) ProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	res, err := ou.orderRepo.GetProductByID(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateOrder is an update method for order
func (ou *orderUsecase) UpdateOrderStatus(ctx context.Context, orderId uint, status string) (*domain.Order, error) {
	// first get order by id
	order, err := ou.orderRepo.GetOrderByID(ctx, orderId)
	if err != nil {
		return nil, err
	}

	// we don't want to update if it does not exist
	if order == nil {
		return nil, errors.New("not found")
	}

	if order.BuyerID != ctx.Value(domain.BuyerKey).(uint) {
		return nil, errors.New("unauthorized")
	}

	// only new status can be update
	if order.Status != domain.OrderStatusNew {
		return nil, errors.New("only new order status eligibile to be update")
	}

	order.Status = status

	// update the order
	res, err := ou.orderRepo.UpdateOrderById(ctx, *order)
	if err != nil {
		return nil, err
	}

	var totalProductSold int64

	for _, v := range order.OrderDetails {
		totalProductSold += int64(v.ProductQuantity)
	}

	orderStatus := domain.OrderStatusCompletedInt
	if status == domain.OrderStatusCancelled {
		orderStatus = domain.OrderStatusCancelledInt
	}

	evt := domain.PayloadEventOrder{
		OrderID:          int64(orderId),
		OrderDate:        time.Time(order.OrderDate).Format("2006-01-02"),
		OrderStatus:      int64(orderStatus),
		TotalRevenue:     float64(order.Amount),
		TotalProductSold: totalProductSold,
	}

	err = ou.orderRepo.PublishOrderEvent(ctx, evt)
	if err != nil {
		log.Println("error publishing order event")
	}

	return res, nil
}

// CreateOrder is an update method for order
func (ou *orderUsecase) CreateOrder(ctx context.Context, req domain.Order) (*domain.Order, error) {
	var res *domain.Order
	//Get ongoing order, if there are ongoing then err
	hasOngoingOrders, err := ou.orderRepo.GetOngoingOrders(ctx, req.BuyerID)
	if err != nil {
		return nil, err
	}

	if hasOngoingOrders {
		return res, errors.New("has ongoing orders")
	}

	// check each product
	for k, v := range req.OrderDetails {

		// get product by id
		resProduct, err := ou.orderRepo.GetProductByID(ctx, v.ProductID)
		if err != nil {
			return nil, err
		}

		//update order details
		req.OrderDetails[k].Product.ProductName = resProduct.ProductName
		req.OrderDetails[k].Product.Price = resProduct.Price
		//sum the amount by
		req.Amount += float64(resProduct.Price) * float64(v.ProductQuantity)
	}

	// insert to table order
	res, err = ou.orderRepo.InsertOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	evt := domain.PayloadEventOrder{
		OrderID:      int64(req.ID),
		OrderDate:    res.CreatedAt.Format("2006-01-02"),
		OrderStatus:  domain.OrderStatusNewInt,
		TotalRevenue: float64(req.Amount),
	}
	err = ou.orderRepo.PublishOrderEvent(ctx, evt)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// OrderByID are method to return list of orders
func (ou *orderUsecase) OrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	res, err := ou.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Orders are method to return list of orders
func (ou *orderUsecase) OrdersByBuyer(ctx context.Context, buyerId uint) ([]domain.Order, error) {
	var res []domain.Order

	res, err := ou.orderRepo.GetOrdersByBuyerID(ctx, buyerId)
	if err != nil {
		return res, err
	}

	return res, nil
}
