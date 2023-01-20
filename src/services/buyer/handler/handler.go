package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/gin/middleware"
	"go.uber.org/fx"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var Module = fx.Provide(
	NewBuyerHandler,
	ProvideGinEngine,
)

func ProvideGinEngine(handler Handler) *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sha_session", store))

	buyer := router.Group("/buyer")

	buyer.POST("/login", handler.Login)
	buyer.GET("/orders", handler.Auth(), handler.Orders)

	orders := router.Group("/orders", handler.Auth())
	orders.GET("/:id", handler.OrderByID)
	orders.POST("/", handler.CreateOrder)
	orders.PUT("/status", handler.UpdateOrderStatus)

	products := router.Group("/products")
	products.GET("/", handler.Products)
	products.GET("/:id", handler.ProductByID)

	router.Use(middleware.LogErrors())

	return router
}
