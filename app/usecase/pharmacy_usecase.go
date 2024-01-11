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

type PharmacyUseCase interface {
	Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	GetById(ctx context.Context, id int64) (*entity.Pharmacy, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	Remove(ctx context.Context, id int64) error
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

	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	return pharmacy, nil
}

func (uc *PharmacyUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacies, err := uc.repo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.repo.CountFindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
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
	pharmacydb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if pharmacydb.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	pharmacy.Id = id
	updated, err := uc.repo.Update(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *PharmacyUseCaseImpl) Remove(ctx context.Context, id int64) error {
	pharmacy, err := uc.GetById(ctx, id)
	if err != nil {
		return err
	}

	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return apperror.ErrForbiddenModifyEntity
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
