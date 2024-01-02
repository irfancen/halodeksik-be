package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductCategoryUseCase interface {
	GetAllProductCategoriesWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
}

type ProductCategoryUseCaseImpl struct {
	repo repository.ProductCategoryRepository
}

func NewProductCategoryUseCaseImpl(repo repository.ProductCategoryRepository) *ProductCategoryUseCaseImpl {
	return &ProductCategoryUseCaseImpl{repo: repo}
}

func (uc *ProductCategoryUseCaseImpl) GetAllProductCategoriesWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	categories, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = categories
	paginatedItems.TotalItems = int64(len(categories))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(categories))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}
