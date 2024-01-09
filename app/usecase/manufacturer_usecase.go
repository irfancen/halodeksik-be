package usecase

import (
	"context"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type ManufacturerUseCase interface {
	Add(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
	GetAllManufacturersWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
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
