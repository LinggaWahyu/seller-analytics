package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository/mocks"
)

func Test_IsUserAuthenticated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		buyerId uint
		wantErr bool
		want    bool
		repo    func() repository.BuyerRepository
	}{
		{
			name:    "error",
			buyerId: 1,
			wantErr: true,
			want:    false,
			repo: func() repository.BuyerRepository {
				m := mocks.NewMockBuyerRepository(ctrl)
				m.EXPECT().Get(gomock.Any(), domain.Buyer{
					Model: yugabyte.Model{
						ID: 1,
					},
				}).Return(nil, errors.New("some error"))
				return m
			},
		},
		{
			name:    "user not authenticated",
			buyerId: 1,
			want:    false,
			repo: func() repository.BuyerRepository {
				m := mocks.NewMockBuyerRepository(ctrl)
				m.EXPECT().Get(gomock.Any(), domain.Buyer{
					Model: yugabyte.Model{
						ID: 1,
					},
				}).Return(nil, nil)
				return m
			},
		},
		{
			name:    "user authenticated",
			buyerId: 1,
			want:    true,
			repo: func() repository.BuyerRepository {
				m := mocks.NewMockBuyerRepository(ctrl)
				m.EXPECT().Get(gomock.Any(), domain.Buyer{
					Model: yugabyte.Model{
						ID: 1,
					},
				}).Return(&domain.Buyer{
					Username: "some user",
				}, nil)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewBuyerUsecase(tt.repo())

			res, err := sut.IsUserAuthenticated(context.TODO(), tt.buyerId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, res, tt.want)
		})
	}
}

func Test_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		username string
		wantErr  bool
		want     *domain.Buyer
		repo     func() repository.BuyerRepository
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewBuyerUsecase(tt.repo())

			res, err := sut.Login(context.TODO(), tt.username)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, res, tt.want)
		})
	}
}
