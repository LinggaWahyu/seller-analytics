package usecase

import (
	"context"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
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
	res, err := bu.buyerRepo.Get(ctx, domain.Buyer{
		Model: yugabyte.Model{
			ID: buyerId,
		},
	})
	if err != nil {
		return false, err
	}

	if res == nil {
		return false, nil
	}
	return true, nil
}

func (bu *buyerUsecase) Login(ctx context.Context, username string) (*domain.Buyer, error) {
	res, err := bu.buyerRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if res == nil {
		// user does not exist yet, create
		res, err = bu.buyerRepo.Create(ctx, domain.Buyer{
			Username: username,
		})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
