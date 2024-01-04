package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PharmacyProductUseCase interface {
	Add(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
}

type PharmacyProductUseCaseImpl struct {
	pharmacyProductRepo repository.PharmacyProductRepository
	pharmacyRepo        repository.PharmacyRepository
	productRepo         repository.ProductRepository
}

func NewPharmacyProductUseCaseImpl(
	pharmacyProductRepo repository.PharmacyProductRepository,
	pharmacyRepo repository.PharmacyRepository,
	productRepo repository.ProductRepository,
) *PharmacyProductUseCaseImpl {
	return &PharmacyProductUseCaseImpl{
		pharmacyProductRepo: pharmacyProductRepo,
		pharmacyRepo:        pharmacyRepo,
		productRepo:         productRepo,
	}
}

func (uc *PharmacyProductUseCaseImpl) Add(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	if pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyProduct.PharmacyId); err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyProduct.PharmacyId)
		}
		return nil, err
	}

	if product, err := uc.productRepo.FindById(ctx, pharmacyProduct.ProductId); err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(product, "Id", pharmacyProduct.ProductId)
		}
		return nil, err
	}

	created, err := uc.pharmacyProductRepo.Create(ctx, pharmacyProduct)
	if err != nil {
		return nil, err
	}
	return created, nil
}
