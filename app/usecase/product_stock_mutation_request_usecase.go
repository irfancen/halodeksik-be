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

type ProductStockMutationRequestUseCase interface {
	Add(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
	GetAllIncoming(ctx context.Context, pharmacyOriginId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllOutgoing(ctx context.Context, pharmacyDestId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
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
	destPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductDestId)
	if err != nil {
		return nil, err
	}

	originPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductOriginId)
	if err != nil {
		return nil, err
	}

	destPharmacy, err := uc.pharmacyRepo.FindById(ctx, destPharmacyProduct.PharmacyId)
	if err != nil {
		return nil, err
	}

	if destPharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if destPharmacyProduct.ProductId != originPharmacyProduct.ProductId {
		return nil, apperror.ErrRequestStockMutationDifferentProduct
	}

	if destPharmacyProduct.Id == originPharmacyProduct.Id {
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

func (uc *ProductStockMutationRequestUseCaseImpl) GetAllIncoming(ctx context.Context, pharmacyOriginId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	if pharmacyOriginId != 0 {
		pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyOriginId)
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyOriginId)
		}
		if err != nil {
			return nil, err
		}
		if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	mutationRequest, err := uc.productStockMutationRequestRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRequestRepo.CountFindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = mutationRequest
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(mutationRequest))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) GetAllOutgoing(ctx context.Context, pharmacyDestId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	if pharmacyDestId != 0 {
		pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyDestId)
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyDestId)
		}
		if err != nil {
			return nil, err
		}
		if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	mutationRequest, err := uc.productStockMutationRequestRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRequestRepo.CountFindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = mutationRequest
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(mutationRequest))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}
