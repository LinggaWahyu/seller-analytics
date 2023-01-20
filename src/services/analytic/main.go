package main

import (
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/http/gin"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/config"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/handler"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/usecase"
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
