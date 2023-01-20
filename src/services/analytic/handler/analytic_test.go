package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/usecase"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/usecase/mocks"
	"gorm.io/datatypes"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// mock gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	ctx.AddParam("date", "20220101")
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHandler_GetAnalyticByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.AnalyticUsecase
		date     string
		wantCode int
		want     GetAnalyticByDateResponse
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/analytic", nil)

				values := req.URL.Query()
				values.Add("date", "2022-01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.AnalyticUsecase {
				m := mocks.NewMockAnalyticUsecase(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), gomock.Any()).Return(&domain.Analytic{
					Date: datatypes.Date(time.Date(2022, 01, 01, 0, 0, 0, 0, time.Local)),
				}, nil)
				return m
			},
			want: GetAnalyticByDateResponse{
				Data: &domain.Analytic{},
			},
		},
		{
			name:     "invalid date format",
			wantCode: http.StatusBadRequest,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/analytic", nil)
				values := req.URL.Query()
				values.Add("date", "2022,01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.AnalyticUsecase {
				m := mocks.NewMockAnalyticUsecase(ctrl)
				return m
			},
			want: GetAnalyticByDateResponse{
				Error: "invalid date format, expect yyyy-mm-dd",
			},
		},
		{
			name:     "error",
			wantCode: http.StatusInternalServerError,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/analytic", nil)
				values := req.URL.Query()
				values.Add("date", "2022-01-01")
				req.URL.RawQuery = values.Encode()
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.AnalyticUsecase {
				m := mocks.NewMockAnalyticUsecase(ctrl)
				m.EXPECT().GetAnalyticByDate(gomock.Any(), gomock.Any()).Return(nil, errors.New("mock error"))
				return m
			},
			want: GetAnalyticByDateResponse{
				Error: "something happened on our end, please try at a later time",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewAnalyticHandler(Params{
				AnalyticUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response GetAnalyticByDateResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)

		})
	}

}
