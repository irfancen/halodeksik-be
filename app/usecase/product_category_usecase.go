package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductCategoryUseCase interface {
	GetAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error)
}

type ProductCategoryUseCaseImpl struct {
	repo repository.ProductCategoryRepository
}

func NewProductCategoryUseCaseImpl(repo repository.ProductCategoryRepository) *ProductCategoryUseCaseImpl {
	return &ProductCategoryUseCaseImpl{repo: repo}
}

func (uc *ProductCategoryUseCaseImpl) GetAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error) {
	categories, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

