package usecase

import (
	"context"
	"errors"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/repository"
)

type BuyerUsecase interface {
	IsUserAuthenticated(ctx context.Context, buyerId uint) (bool, error)
	Login(ctx context.Context, username string) (*domain.Buyer, error)
}

type buyerUsecase struct {
	buyerRepo repository.BuyerRepository
}

func NewBuyerUsecase(buyerRepo repository.BuyerRepository) BuyerUsecase {
	return &buyerUsecase{
		buyerRepo: buyerRepo,
	}
}

func (bu *buyerUsecase) IsUserAuthenticated(ctx context.Context, buyerId uint) (bool, error) {
	return false, errors.New("unimplemented")
}

func (bu *buyerUsecase) Login(ctx context.Context, username string) (*domain.Buyer, error) {
	return nil, errors.New("unimplemented")
}
