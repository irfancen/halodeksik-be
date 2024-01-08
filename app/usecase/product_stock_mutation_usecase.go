package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductStockMutationUseCase interface {
	Add(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error)
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
