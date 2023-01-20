package main

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/config"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/handler"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(fx.Options(
		config.Module,
		repository.Module,
		usecase.Module,
		handler.Module,
		fx.Invoke(gin.ServeHTTP),
	)).Run()
}
