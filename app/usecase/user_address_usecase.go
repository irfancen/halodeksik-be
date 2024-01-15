package usecase

import (
	"context"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type AddressUseCase interface {
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, address entity.Address) (*entity.Address, error)
	Add(ctx context.Context, address entity.Address) (*entity.Address, error)
	GetById(ctx context.Context, id int64) (*entity.Address, error)
	Remove(ctx context.Context, id int64) error
}

type AddressUseCaseImpl struct {
	userAddressRepo repository.UserAddressRepository
	areaRepo        repository.AddressAreaRepository
	locUtil         util.LocationUtil
}

func NewAddressUseCaseImpl(addressRepo repository.UserAddressRepository, areaRepository repository.AddressAreaRepository, locationUtil util.LocationUtil) *AddressUseCaseImpl {
	return &AddressUseCaseImpl{
		userAddressRepo: addressRepo,
		areaRepo:        areaRepository,
		locUtil:         locationUtil,
	}
}

func (uc *AddressUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}

	addresses, err := uc.userAddressRepo.FindAllByUserId(ctx, userId.(int64), param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.userAddressRepo.CountFindAllUserId(ctx, userId.(int64), param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = addresses
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(addresses))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}

func (uc *AddressUseCaseImpl) Add(ctx context.Context, address entity.Address) (*entity.Address, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}

	err := uc.validateCityAndProvince(ctx, address)
	if err != nil {
		return nil, err
	}

	address.ProfileId = userId.(int64)
	created, err := uc.userAddressRepo.Create(ctx, address)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *AddressUseCaseImpl) Edit(ctx context.Context, id int64, address entity.Address) (*entity.Address, error) {
	addressDb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = uc.validateCityAndProvince(ctx, address)
	if err != nil {
		return nil, err
	}

	address.Id = addressDb.Id

	updated, err := uc.userAddressRepo.Update(ctx, address)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *AddressUseCaseImpl) validateCityAndProvince(ctx context.Context, address entity.Address) error {
	city, err := uc.areaRepo.FindCityById(ctx, address.CityId)

	if err != nil {
		return apperror.NewNotFound(city, "id", address.CityId)
	}

	if city.ProvinceId != address.ProvinceId {
		return apperror.ErrInvalidCityProvinceCombi
	}

	province, err := uc.areaRepo.FindProvinceById(ctx, address.ProvinceId)
	if err != nil {
		return apperror.NewNotFound(province, "id", address.ProvinceId)
	}

	err = uc.locUtil.ValidateLatLong(city.Name, province.Name, address.Latitude, address.Longitude)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AddressUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Address, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}

	addressDb, err := uc.userAddressRepo.FindById(ctx, id)
	if err != nil {
		return nil, apperror.NewNotFound(addressDb, "id", id)
	}

	if addressDb.ProfileId != userId {
		return nil, apperror.ErrUnauthorized
	}

	return addressDb, nil

}

func (uc *AddressUseCaseImpl) Remove(ctx context.Context, id int64) error {
	if _, err := uc.GetById(ctx, id); err != nil {
		return err
	}

	err := uc.userAddressRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
