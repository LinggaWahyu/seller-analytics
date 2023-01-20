package domain

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"gorm.io/datatypes"
)

const AnalyticDateFormat = "2006-01-02"

type Analytic struct {
	yugabyte.Model
	AverageOrderValue     float64 `json:"average_order_value"`
	SalesConvertionRate   float32 `json:"sales_conversion_rate"`
	CancellationOrderRate float32 `json:"cancellation_order_rate"`

	DateString string         `json:"date" gorm:"-"`
	Date       datatypes.Date `json:"-"`
}

type StatisticEvent struct {
	TotalRevenue   float64 `json:"total_revenue"`
	CompletedOrder int64   `json:"completed_order"`
	CanceledOrder  int64   `json:"canceled_order"`
	TotalOrder     int64   `json:"total_order"`
	Date           string  `json:"date"`
}
