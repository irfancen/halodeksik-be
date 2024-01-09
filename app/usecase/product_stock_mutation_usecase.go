package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductStockMutationUseCase interface {
	Add(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type ProductStockMutationUseCaseImpl struct {
	productStockMutationRepo repository.ProductStockMutationRepository
	pharmacyProductRepo      repository.PharmacyProductRepository
}

func NewProductStockMutationUseCaseImpl(productStockMutationRepo repository.ProductStockMutationRepository, pharmacyProductRepo repository.PharmacyProductRepository) *ProductStockMutationUseCaseImpl {
	return &ProductStockMutationUseCaseImpl{productStockMutationRepo: productStockMutationRepo, pharmacyProductRepo: pharmacyProductRepo}
}

func (uc *ProductStockMutationUseCaseImpl) findPharmacyProductById(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
		}
		return nil, err
	}
	return pharmacyProduct, nil
}

func (uc *ProductStockMutationUseCaseImpl) Add(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error) {
	pharmacyProduct, err := uc.findPharmacyProductById(ctx, stockMutation.PharmacyProductId)
	if err != nil {
		return nil, err
	}

	if stockMutation.ProductStockMutationTypeId == appconstant.StockMutationTypeReduction &&
		pharmacyProduct.Stock-stockMutation.Stock < 0 {
		return nil, apperror.ErrInsufficientProductStock
	}

	created, err := uc.productStockMutationRepo.Create(ctx, stockMutation)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ProductStockMutationUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	productStockMutation, err := uc.productStockMutationRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRepo.CountFindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = productStockMutation
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(productStockMutation))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}
