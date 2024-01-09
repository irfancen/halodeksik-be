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

type ManufacturerUseCase interface {
	Add(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
	GetById(ctx context.Context, id int64) (*entity.Manufacturer, error)
	GetAllManufacturersWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
}

type ManufacturerUseCaseImpl struct {
	repo     repository.ManufacturerRepository
	uploader appcloud.FileUploader
}

func NewManufacturerUseCaseImpl(repo repository.ManufacturerRepository, uploader appcloud.FileUploader) *ManufacturerUseCaseImpl {
	return &ManufacturerUseCaseImpl{repo: repo, uploader: uploader}
}

func (uc *ManufacturerUseCaseImpl) Add(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error) {
	var (
		err      error
		fileName string
	)
	fileHeader := ctx.Value(appconstant.FormImage)

	if fileHeader != nil {
		fileName, err = uc.uploader.Upload(ctx, fileHeader, manufacturer.GetEntityName())
		if err != nil {
			return nil, err
		}
	}

	if !util.IsEmptyString(fileName) {
		manufacturer.Image = fileName
	}

	created, err := uc.repo.Create(ctx, manufacturer)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ManufacturerUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Manufacturer, error) {
	manufacturer, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(manufacturer, "Id", id)
		}
		return nil, err
	}
	return manufacturer, nil
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

func (uc *ManufacturerUseCaseImpl) Edit(ctx context.Context, id int64, manufacturer entity.Manufacturer) (*entity.Manufacturer, error) {
	var (
		err      error
		fileName string
	)

	manufacturerDb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	manufacturer.Id = id
	manufacturer.Image = manufacturerDb.Image

	fileHeader := ctx.Value(appconstant.FormImage)

	if fileHeader != nil {
		fileName, err = uc.uploader.Upload(ctx, fileHeader, manufacturer.GetEntityName())
		if err != nil {
			return nil, err
		}
	}

	if !util.IsEmptyString(fileName) {
		manufacturer.Image = fileName
	}

	updated, err := uc.repo.Update(ctx, manufacturer)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
