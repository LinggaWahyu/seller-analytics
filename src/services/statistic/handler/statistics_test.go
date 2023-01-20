package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/usecase"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/usecase/mocks"
	"gorm.io/datatypes"
)

func Test_handler_Statistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.StatisticsUsecase
		date     string
		wantCode int
		want     GetStatisticResponse
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/statistic", nil)

				values := req.URL.Query()
				values.Add("date", "2022-01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.StatisticsUsecase {
				m := mocks.NewMockStatisticsUsecase(ctrl)
				m.EXPECT().GetStatistics(gomock.Any(), gomock.Any()).Return(&domain.Statistics{
					Date: datatypes.Date(time.Date(2022, 01, 01, 0, 0, 0, 0, time.Local)),
				}, nil)
				return m
			},
			want: GetStatisticResponse{
				Data: &domain.Statistics{},
			},
		},
		{
			name:     "invalid date format",
			wantCode: http.StatusBadRequest,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/statistic", nil)
				values := req.URL.Query()
				values.Add("date", "2022,01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.StatisticsUsecase {
				m := mocks.NewMockStatisticsUsecase(ctrl)
				return m
			},
			want: GetStatisticResponse{
				Error: "invalid date format, expect yyyy-mm-dd",
			},
		},
		{
			name:     "error",
			wantCode: http.StatusInternalServerError,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/statistic", nil)
				values := req.URL.Query()
				values.Add("date", "2022-01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.StatisticsUsecase {
				m := mocks.NewMockStatisticsUsecase(ctrl)
				m.EXPECT().GetStatistics(gomock.Any(), gomock.Any()).Return(nil, errors.New("mock error"))
				return m
			},
			want: GetStatisticResponse{
				Error: "something happened on our end, please try at a later time",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewStatisticsHandler(Params{
				StatisticsUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response GetStatisticResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
