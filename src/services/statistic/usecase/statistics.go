package usecase

import (
	"context"
	"log"
	"time"

	buyerdomain "github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/repository"
	"gorm.io/datatypes"
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
	res, err := su.statisticsRepo.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (su *statisticsUsecase) HandleOrderEvent(msg domain.PayloadEventOrder) {
	ctx := context.Background()

	orderDate, err := time.Parse(domain.StatisticDateFormat, msg.OrderDate)
	if err != nil {
		log.Println("[HandleOrderEvent] error", err)
		return
	}

	res, err := su.statisticsRepo.GetByDate(ctx, orderDate)
	if err != nil {
		log.Println("[HandleOrderEvent] error", err)
		return
	}

	var resFinal *domain.Statistics
	if res == nil {
		// statistic object not found, create new one
		resFinal, err = su.statisticsRepo.Create(ctx, updateStatisticsData(domain.Statistics{}, msg))
		if err != nil {
			log.Println("[HandleOrderEvent] error", err)
			return
		}
	} else {
		// statistic object found, update
		resFinal, err = su.statisticsRepo.Update(ctx, updateStatisticsData(*res, msg))
		if err != nil {
			log.Println("[HandleOrderEvent] error", err)
			return
		}
	}

	evt := domain.PayloadEventStatistic{
		TotalRevenue:   float64(resFinal.TotalRevenue),
		CompletedOrder: resFinal.CompletedOrder,
		CanceledOrder:  resFinal.CancelledOrder,
		TotalOrder:     resFinal.TotalOrder,
		Date:           resFinal.DateStr,
	}

	err = su.statisticsRepo.PublishEvent(ctx, evt)
	if err != nil {
		log.Println("[HandleOrderEvent] error", err)
	}
}

func updateStatisticsData(statistics domain.Statistics, msg domain.PayloadEventOrder) (result domain.Statistics) {
	date, err := time.Parse(domain.StatisticDateFormat, msg.OrderDate)
	if err != nil {
		log.Println("[updateStatisticsData] error", err)
		return result
	}

	result = domain.Statistics{
		TotalRevenue:     statistics.TotalRevenue,
		TotalProductSold: statistics.TotalProductSold,
		CompletedOrder:   statistics.CompletedOrder,
		CancelledOrder:   statistics.CancelledOrder,
		TotalOrder:       statistics.TotalOrder,
		DateStr:          msg.OrderDate,
		Date:             datatypes.Date(date),
	}

	switch msg.OrderStatus {

	case buyerdomain.OrderStatusCompletedInt:
		result.TotalRevenue += int64(msg.TotalRevenue)
		result.TotalProductSold += msg.TotalProductSold
		result.CompletedOrder += 1
		return result

	case buyerdomain.OrderStatusCancelledInt:
		result.CancelledOrder += 1
		return result

	default:
		result.TotalOrder += 1
		return result
	}
}
