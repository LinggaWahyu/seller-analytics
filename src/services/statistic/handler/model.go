package handler

import (
	httpdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
)

type GetStatisticRequest struct {
	Date string `json:"date"`
}
type GetStatisticResponse = httpdomain.ResponseModel[domain.Statistics]
