package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PharmacyUseCase interface {
	Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
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
