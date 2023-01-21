package usecase

import (
	"context"
	"log"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/repository"
	"gorm.io/datatypes"
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
	res, err := au.analyticRepo.GetAnalyticByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (au *analyticUsecase) HandleStatisticEvent(statisticEvent domain.StatisticEvent) {
	ctx := context.Background()

	date, err := time.Parse(domain.AnalyticDateFormat, statisticEvent.Date)
	if err != nil {
		log.Println("[HandleOrderEvent] error parsing date", err)
		return
	}

	res, err := au.analyticRepo.GetAnalyticByDate(ctx, date)
	if err != nil {
		log.Println("[HandleOrderEvent] error GetAnalyticByDate", err)
		return
	}

	analytic, err := calculateAnalytic(statisticEvent)
	if err != nil {
		log.Println(err)
		return
	}

	if res == nil {
		// analytic object not found, create new one
		_, err = au.analyticRepo.CreateAnalytic(ctx, analytic)
		if err != nil {
			log.Println("[HandleOrderEvent] error CreateAnalytic", err)
			return
		}
	} else {
		// analytic object found, update
		_, err = au.analyticRepo.UpdateAnalytic(ctx, analytic)
		if err != nil {
			log.Println("[HandleOrderEvent] error UpdateAnalytic", err)
			return
		}
	}

}

func calculateAnalytic(statisticEvent domain.StatisticEvent) (domain.Analytic, error) {
	var res domain.Analytic

	if statisticEvent.TotalRevenue > 0 && statisticEvent.CompletedOrder > 0 {
		res.AverageOrderValue = statisticEvent.TotalRevenue / float64(statisticEvent.CompletedOrder)
	}
	if statisticEvent.CompletedOrder > 0 && statisticEvent.TotalOrder > 0 {
		res.SalesConvertionRate = float32(statisticEvent.CompletedOrder) / float32(statisticEvent.TotalOrder) * 100
	}
	if statisticEvent.CanceledOrder > 0 && statisticEvent.TotalOrder > 0 {
		res.CancellationOrderRate = float32(statisticEvent.CanceledOrder) / float32(statisticEvent.TotalOrder) * 100
	}

	date, err := time.Parse(domain.AnalyticDateFormat, statisticEvent.Date)
	if err != nil {
		return domain.Analytic{}, err
	}
	res.Date = datatypes.Date(date)
	return res, nil
}
