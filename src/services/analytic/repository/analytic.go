package repository

import (
	"context"
	"time"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AnalyticRepository interface {
	GetAnalyticByDate(ctx context.Context, date time.Time) (*domain.Analytic, error)
	CreateAnalytic(ctx context.Context, analytic domain.Analytic) (*domain.Analytic, error)
	UpdateAnalytic(ctx context.Context, analytic domain.Analytic) (*domain.Analytic, error)
}

type analyticRepository struct {
	db *gorm.DB
}

func NewAnalyticRepository(db *gorm.DB) AnalyticRepository {
	return &analyticRepository{
		db: db,
	}
}

// GetAnalyticByDate, get analytic by date
func (ar *analyticRepository) GetAnalyticByDate(ctx context.Context, date time.Time) (*domain.Analytic, error) {
	result := domain.Analytic{}

	query := ar.db.WithContext(ctx)
	if err := query.Where("Date = ?", datatypes.Date(date)).First(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	result.DateString = time.Time(result.Date).Format(domain.AnalyticDateFormat)
	return &result, nil
}

// CreateAnalytic create analytic
func (ar *analyticRepository) CreateAnalytic(ctx context.Context, analytic domain.Analytic) (*domain.Analytic, error) {
	if err := ar.db.Create(&analytic).Error; err != nil {
		return nil, err
	}
	analytic.DateString = time.Time(analytic.Date).Format(domain.AnalyticDateFormat)
	return &analytic, nil
}

// UpdateAnalytic update analytic
func (ar *analyticRepository) UpdateAnalytic(ctx context.Context, analytic domain.Analytic) (*domain.Analytic, error) {
	res, err := ar.GetAnalyticByDate(ctx, time.Time(analytic.Date))
	if err != nil {
		return nil, err
	}

	if analytic.AverageOrderValue != 0 {
		res.AverageOrderValue = analytic.AverageOrderValue
	}
	if analytic.SalesConvertionRate != 0 {
		res.SalesConvertionRate = analytic.SalesConvertionRate
	}
	if analytic.CancellationOrderRate != 0 {
		res.CancellationOrderRate = analytic.CancellationOrderRate
	}
	ar.db.Save(&res)

	res.DateString = time.Time(res.Date).Format(domain.AnalyticDateFormat)
	return res, nil
}
