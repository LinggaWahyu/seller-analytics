package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewStatisticsHandler),
	fx.Provide(ProvideGinEngine),
	fx.Invoke(SubscribeOrder),
)

func ProvideGinEngine(handler Handler) *gin.Engine {
	router := gin.Default()

	//router get statistics
	router.GET("/statistic", handler.Statistics)

	return router
}
