package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository/mocks"
)

func Test_orderUsecase_Products(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Product
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			fields: fields{
				orderRepo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []domain.Product{
				{
					Model: yugabyte.Model{
						ID: 1,
					},
					ProductName: "Product 1",
					Price:       100,
				},
			},
			wantErr: false,
			mock: func() {
				mockRepo.EXPECT().GetProducts(gomock.Any()).Return([]domain.Product{
					{
						Model: yugabyte.Model{
							ID: 1,
						},
						ProductName: "Product 1",
						Price:       100,
					},
				}, nil).Times(1)
			},
		},
		{
			name: "error",
			fields: fields{
				orderRepo: mockRepo,
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				mockRepo.EXPECT().GetProducts(gomock.Any()).Return(nil, errors.New("expected error")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			got, err := ou.Products(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderUsecase.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderUsecase_ProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockRepo := mocks.NewMockOrderRepository(ctrl)

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Product
		wantErr bool
		mock    func()
	}{}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			got, err := ou.ProductByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.ProductByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderUsecase.ProductByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderUsecase_OrdersByBuyer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockRepo := mocks.NewMockOrderRepository(ctrl)

	// orderDate := datatypes.Date(time.Now())

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Order
		wantErr bool
		mock    func()
	}{}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			got, err := ou.OrdersByBuyer(tt.args.ctx, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.Orders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderUsecase.Orders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderUsecase_UpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockRepo := mocks.NewMockOrderRepository(ctrl)

	// orderDate := datatypes.Date(time.Now())
	// td := context.TODO()
	// ctx := context.WithValue(td, domain.BuyerKey, uint(1))

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx     context.Context
		orderId uint
		status  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Order
		wantErr bool
		mock    func()
	}{}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			res, err := ou.UpdateOrderStatus(tt.args.ctx, tt.args.orderId, tt.args.status)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want, res)

		})
	}
}

func Test_orderUsecase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockRepo := mocks.NewMockOrderRepository(ctrl)

	// orderDate := datatypes.Date(time.Now())

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx context.Context
		req domain.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Order
		wantErr bool
		mock    func()
	}{}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			if _, err := ou.CreateOrder(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_orderUsecase_OrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockRepo := mocks.NewMockOrderRepository(ctrl)

	// orderDate := datatypes.Date(time.Now())

	type fields struct {
		orderRepo repository.OrderRepository
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Order
		wantErr bool
		mock    func()
	}{}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			ou := &orderUsecase{
				orderRepo: tt.fields.orderRepo,
			}
			got, err := ou.OrderByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderUsecase.OrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderUsecase.OrderByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrderUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)

	type args struct {
		orderRepo repository.OrderRepository
	}
	tests := []struct {
		name string
		args args
		want OrderUsecase
	}{
		{
			name: "success",
			args: args{
				orderRepo: mockRepo,
			},
			want: &orderUsecase{
				orderRepo: mockRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderUsecase(tt.args.orderRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}
