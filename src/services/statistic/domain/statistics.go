package domain

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"gorm.io/datatypes"
)

const StatisticDateFormat = "2006-01-02"

type PayloadEventOrder struct {
	OrderID          int64   `json:"order_id"`
	OrderDate        string  `json:"order_date"`
	OrderStatus      int64   `json:"order_status"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalProductSold int64   `json:"total_product_sold"`
}

type PayloadEventStatistic struct {
	TotalRevenue   float64 `json:"total_revenue"`
	CompletedOrder int64   `json:"completed_order"`
	CanceledOrder  int64   `json:"canceled_order"`
	TotalOrder     int64   `json:"total_order"`
	Date           string  `json:"date"`
}

type Statistics struct {
	yugabyte.Model
	TotalRevenue     int64 `json:"total_revenue"`
	TotalProductSold int64 `json:"total_product_sold"`
	CompletedOrder   int64 `json:"completed_order"`
	CancelledOrder   int64 `json:"cancelled_order"`
	TotalOrder       int64 `json:"total_order"`

	DateStr string         `json:"date" gorm:"-"`
	Date    datatypes.Date `json:"-"`
}
