package handler

import (
	httpdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
)

type LoginRequest struct {
	Username string
}
type LoginResponse = httpdomain.ResponseModel[domain.Buyer]

type GetProductsRequest struct {
	Id int64
}

type ProductsResponse = httpdomain.ResponseModel[[]domain.Product]
type ProductByIDResponse = httpdomain.ResponseModel[domain.Product]

type UpdateOrderRequest struct {
	ID     uint   `json:"order_id"`
	Status string `json:"status"`
}

type CreateOrderRequest struct {
	BuyerID  uint                            `json:"buyer_id"`
	Products []CreateOrderRequestProductData `json:"products"`
}

type CreateOrderRequestProductData struct {
	ProductID  uint  `json:"product_id"`
	ProductQty int32 `json:"product_qty"`
}

type OrderResponse = httpdomain.ResponseModel[domain.Order]
type OrdersResponse = httpdomain.ResponseModel[[]domain.Order]
