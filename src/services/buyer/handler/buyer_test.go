package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/usecase"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/usecase/mocks"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// mock gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
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

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.BuyerUsecase
		username string
		wantCode int
		want     LoginResponse
	}{
		{
			name:     "invalid binding",
			wantCode: http.StatusBadRequest,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/buyer/login", bytes.NewReader(nil))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.BuyerUsecase {
				return mocks.NewMockBuyerUsecase(ctrl)
			},
			want: LoginResponse{
				Error: "invalid body type",
			},
		},
		{
			name:     "login error",
			wantCode: http.StatusInternalServerError,
			request: func() *http.Request {
				loginReq := LoginRequest{
					Username: "testuser",
				}
				breq, _ := json.Marshal(loginReq)
				req, _ := http.NewRequest(http.MethodPost, "/buyer/login", bytes.NewReader(breq))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.BuyerUsecase {
				m := mocks.NewMockBuyerUsecase(ctrl)
				m.EXPECT().Login(gomock.Any(), "testuser").Return(nil, errors.New("mock error"))
				return m
			},
			want: LoginResponse{
				Error: "something happened on our end, please try at a later time",
			},
		},
		{
			name:     "login success",
			wantCode: http.StatusOK,
			request: func() *http.Request {
				loginReq := LoginRequest{
					Username: "testuser",
				}
				breq, _ := json.Marshal(loginReq)
				req, _ := http.NewRequest(http.MethodPost, "/buyer/login", bytes.NewReader(breq))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			usecase: func() usecase.BuyerUsecase {
				m := mocks.NewMockBuyerUsecase(ctrl)
				m.EXPECT().Login(gomock.Any(), "testuser").Return(&domain.Buyer{
					Username: "testuser",
				}, nil)
				return m
			},
			want: LoginResponse{
				Data: &domain.Buyer{
					Username: "testuser",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				BuyerUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response LoginResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)

		})
	}

}

func Test_handler_Products(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		wantCode int
		want     ProductsResponse
	}{
		{
			name:     "error",
			wantCode: http.StatusInternalServerError,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/products/", nil)
				return req
			},
			usecase: func() usecase.OrderUsecase {
				m := mocks.NewMockOrderUsecase(ctrl)
				m.EXPECT().Products(gomock.Any()).Return([]domain.Product{}, errors.New("expected error")).Times(1)
				return m
			},
			want: ProductsResponse{
				Error: "something happened on our end, please try at a later time",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response ProductsResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)

		})
	}
}

func Test_handler_ProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		username string
		wantCode int
		want     ProductByIDResponse
	}{
		{
			name: "success",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
				return req
			},
			usecase: func() usecase.OrderUsecase {
				m := mocks.NewMockOrderUsecase(ctrl)
				resp := &domain.Product{
					Model: yugabyte.Model{
						ID: 1,
					},
					ProductName: "Product 1",
					Price:       100,
				}
				m.EXPECT().ProductByID(gomock.Any(), gomock.Any()).Return(resp, nil).Times(1)
				return m
			},
			wantCode: http.StatusOK,
			want: ProductByIDResponse{
				Data: &domain.Product{
					Model:       yugabyte.Model{ID: 1},
					ProductName: "Product 1",
					Price:       100,
				},
			},
		},
		{
			name: "error",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
				return req
			},
			usecase: func() usecase.OrderUsecase {
				m := mocks.NewMockOrderUsecase(ctrl)
				m.EXPECT().ProductByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("expected error")).Times(1)
				return m
			},
			wantCode: http.StatusInternalServerError,
			want: ProductByIDResponse{
				Error: "something happened on our end, please try at a later time",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response ProductByIDResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}

func Test_handler_OrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		username string
		wantCode int
		want     OrderResponse
	}{
		{
			name: "error no auth",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/orders/1", nil)
				return req
			},
			usecase: func() usecase.OrderUsecase {
				m := mocks.NewMockOrderUsecase(ctrl)

				return m
			},
			wantCode: http.StatusBadRequest,
			want: OrderResponse{
				Error: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response OrderResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}

func Test_handler_CreateOrder(t *testing.T) {

	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		username string
		wantCode int
		want     OrderResponse
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response OrderResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}

func Test_handler_UpdateOrderStatus(t *testing.T) {
	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		username string
		wantCode int
		want     OrderResponse
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response OrderResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}

func Test_handler_Orders(t *testing.T) {
	tests := []struct {
		name     string
		request  func() *http.Request
		usecase  func() usecase.OrderUsecase
		username string
		wantCode int
		want     OrderResponse
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			sut := NewBuyerHandler(Params{
				OrderUsecase: tt.usecase(),
			})

			router := ProvideGinEngine(sut)
			router.ServeHTTP(recorder, tt.request())

			var response OrderResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tt.want, response)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
