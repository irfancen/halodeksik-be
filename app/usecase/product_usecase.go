package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/env"
	"halodeksik-be/app/repository"
	"mime/multipart"
	"path/filepath"
)

type ProductUseCase interface {
	Add(ctx context.Context, product entity.Product) (*entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, product entity.Product) (*entity.Product, error)
	Remove(ctx context.Context, id int64) error
}

type ProductUseCaseImpl struct {
	repo        repository.ProductRepository
	uploader    appcloud.FileUploader
	cloudUrl    string
	cloudFolder string
}

func NewProductUseCaseImpl(repo repository.ProductRepository, uploader appcloud.FileUploader) *ProductUseCaseImpl {
	cloudUrl := env.Get("GCLOUD_STORAGE_CDN")
	cloudFolder := env.Get("GCLOUD_STORAGE_FOLDER_PRODUCTS")
	return &ProductUseCaseImpl{repo: repo, uploader: uploader, cloudUrl: cloudUrl, cloudFolder: cloudFolder}
}

func (uc *ProductUseCaseImpl) Add(ctx context.Context, product entity.Product) (*entity.Product, error) {
	fileHeader := ctx.Value(appconstant.FormImage).(*multipart.FileHeader)
	if fileHeader == nil {
		return nil, apperror.ErrProductImageDoesNotExistInContext
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	extension := filepath.Ext(fileHeader.Filename)
	createUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	fileName := fmt.Sprintf("%s%s", createUUID.String(), extension)

	err = uc.uploader.SendToBucket(ctx, file, fmt.Sprintf("%s/", uc.cloudFolder), fileName)
	if err != nil {
		return nil, err
	}
	product.Image = fmt.Sprintf("%s/%s/%s", uc.cloudUrl, uc.cloudFolder, fileName)

	created, err := uc.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *ProductUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	product, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(product, "Id", id)
		}
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	products, err := uc.repo.FindAll(ctx, param)
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
	paginatedItems.Items = products
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(products))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}

func (uc *ProductUseCaseImpl) Edit(ctx context.Context, id int64, product entity.Product) (*entity.Product, error) {
	productDb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Id = id
	product.Image = productDb.Image

	fileHeaderAny := ctx.Value(appconstant.FormImage)
	if fileHeaderAny != nil {
		fileHeader := fileHeaderAny.(*multipart.FileHeader)
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		extension := filepath.Ext(fileHeader.Filename)
		updateUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		fileName := fmt.Sprintf("%s%s", updateUUID.String(), extension)

		err = uc.uploader.SendToBucket(ctx, file, fmt.Sprintf("%s/", uc.cloudFolder), fileName)
		if err != nil {
			return nil, err
		}
		product.Image = fmt.Sprintf("%s/%s/%s", uc.cloudUrl, uc.cloudFolder, fileName)
	}

	updated, err := uc.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *ProductUseCaseImpl) Remove(ctx context.Context, id int64) error {
	if _, err := uc.GetById(ctx, id); err != nil {
		return err
	}

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
