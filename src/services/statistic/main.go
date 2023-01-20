package main

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/config"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/handler"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(fx.Options(
		config.Module,
		repository.Module,
		usecase.Module,
		handler.Module,
		fx.Invoke(gin.ServeHTTP),
	))
}
