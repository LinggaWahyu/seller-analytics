package repository

import (
	"context"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/messagequeue"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StatisticsRepository interface {
	GetByDate(ctx context.Context, date time.Time) (*domain.Statistics, error)
	Create(ctx context.Context, stat domain.Statistics) (*domain.Statistics, error)
	Update(ctx context.Context, stat domain.Statistics) (*domain.Statistics, error)
	PublishEvent(ctx context.Context, event domain.PayloadEventStatistic) error
}

type statisticsRepository struct {
	db               *gorm.DB
	repoCoreRabbitMQ messagequeue.Publisher[domain.PayloadEventStatistic]
}

func NewStatisticsRepository(db *gorm.DB, repoCoreRabbitMQ messagequeue.Publisher[domain.PayloadEventStatistic]) StatisticsRepository {
	return &statisticsRepository{
		db:               db,
		repoCoreRabbitMQ: repoCoreRabbitMQ,
	}
}

func (sr *statisticsRepository) GetByDate(ctx context.Context, date time.Time) (*domain.Statistics, error) {
	result := domain.Statistics{}

	query := sr.db.WithContext(ctx)
	if err := query.Where("Date = ?", datatypes.Date(date)).First(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	result.DateStr = date.Format(domain.StatisticDateFormat)
	return &result, nil
}

func (sr *statisticsRepository) Create(ctx context.Context, req domain.Statistics) (*domain.Statistics, error) {
	if err := sr.db.Create(&req).Error; err != nil {
		return nil, err
	}
	req.DateStr = time.Time(req.Date).Format(domain.StatisticDateFormat)
	return &req, nil
}

func (sr *statisticsRepository) Update(ctx context.Context, req domain.Statistics) (*domain.Statistics, error) {
	res, err := sr.GetByDate(ctx, time.Time(req.Date))

	if err != nil {
		return nil, err
	}

	res.TotalRevenue = req.TotalRevenue
	res.TotalProductSold = req.TotalProductSold
	res.CompletedOrder = req.CompletedOrder
	res.CancelledOrder = req.CancelledOrder
	res.TotalOrder = req.TotalOrder
	res.DateStr = req.DateStr
	res.Date = req.Date
	sr.db.Save(&res)

	res.DateStr = time.Time(res.Date).Format(domain.StatisticDateFormat)
	return res, nil
}

func (sr *statisticsRepository) PublishEvent(ctx context.Context, event domain.PayloadEventStatistic) error {
	err := sr.repoCoreRabbitMQ.Publish(ctx, messagequeue.PublishConfig{}, event)
	if err != nil {
		return err
	}

	return nil
}
