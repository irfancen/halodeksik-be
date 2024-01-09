package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type DoctorSpecializationUseCase interface {
	GetAllSpecsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
}

type DoctorSpecializationUseCaseImpl struct {
	repo repository.DoctorSpecializationRepository
}

func NewDoctorSpecializationUseCaseImpl(repo repository.DoctorSpecializationRepository) *DoctorSpecializationUseCaseImpl {
	return &DoctorSpecializationUseCaseImpl{repo: repo}
}

func (uc *DoctorSpecializationUseCaseImpl) GetAllSpecsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	doctorSpecs, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = doctorSpecs
	paginatedItems.TotalItems = int64(len(doctorSpecs))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(doctorSpecs))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}
