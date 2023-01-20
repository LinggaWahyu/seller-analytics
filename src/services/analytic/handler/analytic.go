package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/usecase"
	statdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"go.uber.org/fx"
)

type Handler interface {
	GetAnalyticByDate(ctx *gin.Context)
}

type handler struct {
	AnalyticUsecase usecase.AnalyticUsecase
}

type Params struct {
	fx.In
	AnalyticUsecase usecase.AnalyticUsecase
}

func NewAnalyticHandler(param Params) Handler {
	return &handler{
		AnalyticUsecase: param.AnalyticUsecase,
	}
}

func (h *handler) GetAnalyticByDate(ctx *gin.Context) {

	strDate := ctx.Query("date")
	date := time.Now()
	var err error
	if strDate != "" {
		date, err = time.Parse(domain.AnalyticDateFormat, strDate)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, GetAnalyticByDateResponse{
				Error: "invalid date format, expect yyyy-mm-dd",
			})
			return
		}
	}

	res, err := h.AnalyticUsecase.GetAnalyticByDate(ctx, date)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetAnalyticByDateResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, GetAnalyticByDateResponse{
		Data: res,
	})
}

func SubscribeStatistic(repoCoreRabbitMQ messagequeue.Subscriber[statdomain.PayloadEventStatistic], usecase usecase.AnalyticUsecase) {
	go func() {
		// TODO write code here
	}()
}
