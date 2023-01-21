package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/repository/mocks"
	"gorm.io/datatypes"
)

func Test_analyticUsecase_GetAnalyticByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name    string
		date    time.Time
		want    *domain.Analytic
		wantErr bool
		repo    func() repository.AnalyticRepository
	}{
		{
			name: "sukses",
			date: date,
			want: &domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   80,
				CancellationOrderRate: 20,
				Date:                  datatypes.Date(date),
			},
			repo: func() repository.AnalyticRepository {
				m := mocks.NewMockAnalyticRepository(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), date).Return(&domain.Analytic{
					AverageOrderValue:     100,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  datatypes.Date(date),
				}, nil)
				return m
			},
		},
		{
			name:    "error",
			date:    date,
			want:    nil,
			wantErr: true,
			repo: func() repository.AnalyticRepository {
				m := mocks.NewMockAnalyticRepository(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), date).Return(nil, errors.New("mock error"))
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewAnalyticsUsecase(tt.repo())
			got, err := au.GetAnalyticByDate(context.TODO(), tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("analyticUsecase.GetAnalyticByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("analyticUsecase.GetAnalyticByDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_analyticUsecase_HandleStatisticEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dateTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	dateString := dateTime.Format(domain.AnalyticDateFormat)
	date := datatypes.Date(dateTime)

	tests := []struct {
		name     string
		analytic domain.StatisticEvent
		repo     func() repository.AnalyticRepository
	}{
		{
			name: "sukses",
			analytic: domain.StatisticEvent{
				TotalRevenue:   100,
				CompletedOrder: 4,
				CanceledOrder:  1,
				TotalOrder:     5,
				Date:           dateString,
			},
			repo: func() repository.AnalyticRepository {
				m := mocks.NewMockAnalyticRepository(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), dateTime).Return(&domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   100,
					CancellationOrderRate: 0,
					Date:                  date,
				}, nil)
				m.EXPECT().UpdateAnalytic(gomock.Any(), domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  date,
				}).Return(&domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  date,
				}, nil)
				return m
			},
		},
		{
			name: "error update",
			analytic: domain.StatisticEvent{
				TotalRevenue:   100,
				CompletedOrder: 4,
				CanceledOrder:  1,
				TotalOrder:     5,
				Date:           dateString,
			},
			repo: func() repository.AnalyticRepository {
				m := mocks.NewMockAnalyticRepository(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), dateTime).Return(&domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   100,
					CancellationOrderRate: 0,
					Date:                  date,
				}, nil)
				m.EXPECT().UpdateAnalytic(gomock.Any(), domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  date,
				}).Return(nil, errors.New("mock error"))
				return m
			},
		},
		{
			name: "no record, create new one",
			analytic: domain.StatisticEvent{
				TotalRevenue:   100,
				CompletedOrder: 4,
				CanceledOrder:  1,
				TotalOrder:     5,
				Date:           dateString,
			},
			repo: func() repository.AnalyticRepository {
				m := mocks.NewMockAnalyticRepository(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), dateTime).Return(nil, nil)
				m.EXPECT().CreateAnalytic(gomock.Any(), domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  date,
				}).Return(&domain.Analytic{
					AverageOrderValue:     25,
					SalesConvertionRate:   80,
					CancellationOrderRate: 20,
					Date:                  date,
				}, nil)
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewAnalyticsUsecase(tt.repo())
			au.HandleStatisticEvent(tt.analytic)
		})
	}
}