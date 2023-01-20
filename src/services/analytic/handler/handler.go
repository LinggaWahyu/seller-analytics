package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAnalyticHandler),
	fx.Provide(ProvideGinEngine),
	fx.Invoke(SubscribeStatistic),
)

func ProvideGinEngine(handler Handler) *gin.Engine {
	router := gin.Default()

	//get analytic by date
	router.GET("/analytic", handler.GetAnalyticByDate)

	return router
}
