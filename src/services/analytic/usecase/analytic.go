package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/repository"
)

type AnalyticUsecase interface {
	GetAnalyticByDate(ctx context.Context, date time.Time) (*domain.Analytic, error)
	HandleStatisticEvent(statisticEvent domain.StatisticEvent)
}

type analyticUsecase struct {
	analyticRepo repository.AnalyticRepository
}

func NewAnalyticsUsecase(analyticRepo repository.AnalyticRepository) AnalyticUsecase {
	return &analyticUsecase{
		analyticRepo: analyticRepo,
	}
}

func (au *analyticUsecase) GetAnalyticByDate(ctx context.Context, date time.Time) (*domain.Analytic, error) {
	return nil, errors.New("unimplemented")
}

func (au *analyticUsecase) HandleStatisticEvent(statisticEvent domain.StatisticEvent) {
	// TODO write code here
}
