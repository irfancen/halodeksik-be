package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type DrugClassificationUseCase interface {
	GetAllWithoutParams(ctx context.Context) ([]*entity.DrugClassification, error)
}

type DrugClassificationUseCaseImpl struct {
	repo repository.DrugClassificationRepository
}

func NewDrugClassificationUseCaseImpl(repo repository.DrugClassificationRepository) *DrugClassificationUseCaseImpl {
	return &DrugClassificationUseCaseImpl{repo: repo}
}

func (uc *DrugClassificationUseCaseImpl) GetAllWithoutParams(ctx context.Context) ([]*entity.DrugClassification, error) {
	drugClassifications, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}
	return drugClassifications, nil
}

