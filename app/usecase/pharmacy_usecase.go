package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PharmacyUseCase interface {
	Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	GetById(ctx context.Context, id int64) (*entity.Pharmacy, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
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

func (uc *PharmacyUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Pharmacy, error) {
	pharmacy, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", id)
		}
		return nil, err
	}
	return pharmacy, nil
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

func (uc *PharmacyUseCaseImpl) Edit(ctx context.Context, id int64, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	if _, err := uc.GetById(ctx, id); err != nil {
		return nil, err
	}
	pharmacy.Id = id
	updated, err := uc.repo.Update(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
