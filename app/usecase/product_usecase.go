package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductUseCase interface {
	Add(ctx context.Context, product entity.Product) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, product entity.Product) (*entity.Product, error)
	Remove(ctx context.Context, id int64) error
}

type ProductUseCaseImpl struct {
	repo repository.ProductRepository
}

func NewProductUseCaseImpl(repo repository.ProductRepository) *ProductUseCaseImpl {
	return &ProductUseCaseImpl{repo: repo}
}

func (uc *ProductUseCaseImpl) Add(ctx context.Context, product entity.Product) (*entity.Product, error) {
	created, err := uc.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ProductUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	product, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(product, "Id", id)
		}
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	products, err := uc.repo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, totalPages, err := uc.repo.CountFindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = products
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(products))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}

func (uc *ProductUseCaseImpl) Edit(ctx context.Context, id int64, product entity.Product) (*entity.Product, error) {
	if _, err := uc.GetById(ctx, id); err != nil {
		return nil, err
	}
	product.Id = id
	updated, err := uc.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *ProductUseCaseImpl) Remove(ctx context.Context, id int64) error {
	if _, err := uc.GetById(ctx, id); err != nil {
		return err
	}

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
