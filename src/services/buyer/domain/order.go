package domain

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"gorm.io/datatypes"
)

const OrderDateFormat = "2006-01-02"

type Order struct {
	yugabyte.Model

	BuyerID      uint    `json:"buyer_id"`
	Status       string  `json:"status"`
	Amount       float64 `json:"amount"`
	InvoiceNo    string  `json:"invoice_number"`
	OrderDateStr string  `json:"order_date"`

	OrderDate    datatypes.Date `json:"-"`
	OrderDetails []OrderDetail  `json:"order_details,omitempty"`
}

type OrderDetail struct {
	yugabyte.Model
	ProductID       uint    `json:"product_id"`
	Product         Product `json:"product"`
	ProductQuantity int     `json:"product_quantity"`
	OrderID         uint    `json:"order_id"`
}

type Product struct {
	yugabyte.Model
	ProductName string  `json:"product_name"`
	Price       float32 `json:"price"`
}

type PayloadEventOrder struct {
	OrderID          int64   `json:"order_id"`
	OrderDate        string  `json:"order_date"`
	OrderStatus      int64   `json:"order_status"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalProductSold int64   `json:"total_product_sold"`
}

var UpdateOrderStatusVal = map[string]string{
	"cancelled": "cancelled",
	"completed": "completed",
}

