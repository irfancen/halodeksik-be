package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ManufacturerUseCase interface {
	GetAllManufacturersWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
}

type ManufacturerUseCaseImpl struct {
	repo repository.ManufacturerRepository
}

func NewManufacturerUseCaseImpl(repo repository.ManufacturerRepository) *ManufacturerUseCaseImpl {
	return &ManufacturerUseCaseImpl{repo: repo}
}

func (uc *ManufacturerUseCaseImpl) GetAllManufacturersWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	manufacturers, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = manufacturers
	paginatedItems.TotalItems = int64(len(manufacturers))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(manufacturers))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}
