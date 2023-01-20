package repository

import (
	"context"

	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"gorm.io/gorm"
)

// BuyerRepository, interface for buyer repository
type BuyerRepository interface {
	Get(ctx context.Context, buyer domain.Buyer) (*domain.Buyer, error)
	Create(ctx context.Context, buyer domain.Buyer) (*domain.Buyer, error)
	GetByUsername(ctx context.Context, username string) (*domain.Buyer, error)
}

// buyerRepository, concrete implementation of buyer repository
type buyerRepository struct {
	db *gorm.DB
}

// NewBuyerRepository, constructor function for buyer repository
func NewBuyerRepository(db *gorm.DB) BuyerRepository {
	return &buyerRepository{
		db: db,
	}
}

// Get, gets buyer by primary key, note that only buyer.ID is used for the query
func (br *buyerRepository) Get(ctx context.Context, buyer domain.Buyer) (*domain.Buyer, error) {
	query := br.db.WithContext(ctx)
	if err := query.First(&buyer, buyer.ID).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &buyer, nil
}

// GetByUsername, gets buyer by username
func (br *buyerRepository) GetByUsername(ctx context.Context, username string) (*domain.Buyer, error) {
	result := domain.Buyer{}

	query := br.db.WithContext(ctx)
	if err := query.Where("username = ?", username).First(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, nil
}

// Create, insert buyer entity into database
func (br *buyerRepository) Create(ctx context.Context, buyer domain.Buyer) (*domain.Buyer, error) {
	if err := br.db.Create(&buyer).Error; err != nil {
		return nil, err
	}
	return &buyer, nil
}
