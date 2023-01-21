package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/usecase"
	"go.uber.org/fx"
)

type Handler interface {
	Statistics(*gin.Context)
}

type handler struct {
	StatisticsUsecase usecase.StatisticsUsecase
}

type Params struct {
	fx.In
	StatisticsUsecase usecase.StatisticsUsecase
}

func NewStatisticsHandler(param Params) Handler {
	return &handler{
		StatisticsUsecase: param.StatisticsUsecase,
	}
}

func (h *handler) Statistics(ctx *gin.Context) {
	strDate := ctx.Query("date")

	date := time.Now()
	var err error
	if strDate != "" {
		date, err = time.Parse(domain.StatisticDateFormat, strDate)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, GetStatisticResponse{
				Error: "invalid date format, expect yyyy-mm-dd",
			})
			return
		}
	}

	res, err := h.StatisticsUsecase.GetStatistics(ctx, date)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, GetStatisticResponse{
			Error: "something happened on our end, please try at a later time",
		})
		return
	}

	ctx.JSON(http.StatusOK, GetStatisticResponse{
		Data: res,
	})
}

func SubscribeOrder(
	repoCoreRabbitMQ messagequeue.Subscriber[domain.PayloadEventOrder],
	usecase usecase.StatisticsUsecase) {
	go func() {
		err := repoCoreRabbitMQ.Subscribe(messagequeue.SubscribeConfig{
			AutoAck: true,
		}, func(msg domain.PayloadEventOrder) {
			if msg.OrderDate == "" {
				log.Println("invalid message: date can't be empty")
			}
			usecase.HandleOrderEvent(msg)
		})

		if err != nil {
			log.Println(err)
		}
	}()
}
