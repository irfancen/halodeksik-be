package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type AddressAreaUseCase interface {
	GetAllProvinces(ctx context.Context) ([]*entity.Province, error)
	GetAllCities(ctx context.Context) ([]*entity.City, error)
}

type AddressAreaUseCaseImpl struct {
	repo repository.AddressAreaRepository
}

func NewAddressAreaUseCaseImpl(repo repository.AddressAreaRepository) *AddressAreaUseCaseImpl {
	return &AddressAreaUseCaseImpl{repo: repo}
}

func (uc *AddressAreaUseCaseImpl) GetAllProvinces(ctx context.Context) ([]*entity.Province, error) {
	provinces, err := uc.repo.FindAllProvince(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (uc *AddressAreaUseCaseImpl) GetAllCities(ctx context.Context) ([]*entity.City, error) {
	cities, err := uc.repo.FindAllCities(ctx)
	if err != nil {
		return nil, err
	}
	return cities, nil
}
