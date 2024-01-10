package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type DoctorSpecializationUseCase interface {
	Add(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	GetById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error)
	GetAllSpecsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
}

type DoctorSpecializationUseCaseImpl struct {
	repo     repository.DoctorSpecializationRepository
	uploader appcloud.FileUploader
}

func NewDoctorSpecializationUseCaseImpl(repo repository.DoctorSpecializationRepository, uploader appcloud.FileUploader) *DoctorSpecializationUseCaseImpl {
	return &DoctorSpecializationUseCaseImpl{repo: repo, uploader: uploader}
}

func (uc *DoctorSpecializationUseCaseImpl) Add(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	var (
		err      error
		fileName string
	)
	fileHeader := ctx.Value(appconstant.FormImage)

	if fileHeader != nil {
		fileName, err = uc.uploader.Upload(ctx, fileHeader, specialization.GetEntityName())
		if err != nil {
			return nil, err
		}
	}

	if !util.IsEmptyString(fileName) {
		specialization.Image = fileName
	}

	created, err := uc.repo.Create(ctx, specialization)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *DoctorSpecializationUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error) {
	specialization, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(specialization, "Id", id)
		}
		return nil, err
	}
	return specialization, nil
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
