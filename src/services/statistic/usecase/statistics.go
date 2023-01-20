package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/repository"
)

type StatisticsUsecase interface {
	GetStatistics(ctx context.Context, date time.Time) (*domain.Statistics, error)
	HandleOrderEvent(domain.PayloadEventOrder)
}

type statisticsUsecase struct {
	statisticsRepo repository.StatisticsRepository
}

func NewStatisticsUsecase(statisticsRepo repository.StatisticsRepository) StatisticsUsecase {
	return &statisticsUsecase{
		statisticsRepo: statisticsRepo,
	}
}

func (su *statisticsUsecase) GetStatistics(ctx context.Context, date time.Time) (*domain.Statistics, error) {
	return nil, errors.New("unimplemented")
}

func (su *statisticsUsecase) HandleOrderEvent(msg domain.PayloadEventOrder) {
	// TODO write code here
}
