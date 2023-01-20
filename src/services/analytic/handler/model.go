package handler

import (
	httpdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
)

type GetAnalyticRequest struct {
	Date string `json:"date"`
}
type GetAnalyticByDateResponse = httpdomain.ResponseModel[domain.Analytic]
