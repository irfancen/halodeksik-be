package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ManufacturerUseCase interface {
	GetAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error)
}

type ManufacturerUseCaseImpl struct {
	repo repository.ManufacturerRepository
}

func NewManufacturerUseCaseImpl(repo repository.ManufacturerRepository) *ManufacturerUseCaseImpl {
	return &ManufacturerUseCaseImpl{repo: repo}
}

func (uc *ManufacturerUseCaseImpl) GetAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error) {
	manufacturers, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}
	return manufacturers, nil
}
