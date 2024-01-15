package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductStockMutationRequestUseCase interface {
	Add(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
}

type ProductStockMutationRequestUseCaseImpl struct {
	productStockMutationRequestRepo repository.ProductStockMutationRequestRepository
	pharmacyProductRepo             repository.PharmacyProductRepository
	pharmacyRepo                    repository.PharmacyRepository
}

func NewProductStockMutationRequestUseCaseImpl(productStockMutationRequestRepo repository.ProductStockMutationRequestRepository, pharmacyProductRepo repository.PharmacyProductRepository, pharmacyRepo repository.PharmacyRepository) *ProductStockMutationRequestUseCaseImpl {
	return &ProductStockMutationRequestUseCaseImpl{productStockMutationRequestRepo: productStockMutationRequestRepo, pharmacyProductRepo: pharmacyProductRepo, pharmacyRepo: pharmacyRepo}
}

func (uc *ProductStockMutationRequestUseCaseImpl) findPharmacyProductById(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
		}
		return nil, err
	}
	return pharmacyProduct, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) Add(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	originPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductOriginId)
	if err != nil {
		return nil, err
	}

	destPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductDestId)
	if err != nil {
		return nil, err
	}

	originPharmacy, err := uc.pharmacyRepo.FindById(ctx, originPharmacyProduct.PharmacyId)
	if err != nil {
		return nil, err
	}

	if originPharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if originPharmacyProduct.ProductId != destPharmacyProduct.ProductId {
		return nil, apperror.ErrRequestStockMutationDifferentProduct
	}

	if originPharmacyProduct.Id == destPharmacyProduct.Id {
		return nil, apperror.ErrRequestStockMutationFromOwnPharmacy
	}

	if originPharmacyProduct.Stock < mutationRequest.Stock {
		return nil, apperror.ErrInsufficientProductStock
	}

	mutationRequest.ProductStockMutationRequestStatusId = appconstant.StockMutationRequestStatusPending
	created, err := uc.productStockMutationRequestRepo.Create(ctx, mutationRequest)
	if err != nil {
		return nil, err
	}
	return created, nil
}
