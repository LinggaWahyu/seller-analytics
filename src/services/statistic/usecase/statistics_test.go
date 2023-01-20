package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/repository/mocks"
)

func Test_statisticsUsecase_GetStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name    string
		date    time.Time
		want    *domain.Statistics
		wantErr bool
		repo    func() repository.StatisticsRepository
	}{
		{
			name:    "error",
			date:    date,
			want:    nil,
			wantErr: true,
			repo: func() repository.StatisticsRepository {
				m := mocks.NewMockStatisticsRepository(ctrl)
				m.EXPECT().GetByDate(gomock.Any(), date).Return(nil, errors.New("mock error"))
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewStatisticsUsecase(tt.repo())
			got, err := au.GetStatistics(context.TODO(), tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("statisticsUsecase.GetStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statisticsUsecase.GetStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}
