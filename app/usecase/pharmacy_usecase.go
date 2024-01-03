package usecase

import (
	"context"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PharmacyUseCase interface {
	Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type PharmacyUseCaseImpl struct {
	repo repository.PharmacyRepository
}

func NewPharmacyUseCseImpl(repo repository.PharmacyRepository) *PharmacyUseCaseImpl {
	return &PharmacyUseCaseImpl{repo: repo}
}

func (uc *PharmacyUseCaseImpl) Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	created, err := uc.repo.Create(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *PharmacyUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacies, err := uc.repo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, totalPages, err := uc.repo.CountFindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = pharmacies
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(pharmacies))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}
