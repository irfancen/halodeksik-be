package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type UserUseCase interface {
	AddAdmin(ctx context.Context, admin entity.User) (*entity.User, error)
	GetById(ctx context.Context, id int64) (*entity.User, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type UserUseCaseImpl struct {
	repo repository.UserRepository
	util util.AuthUtil
}

func NewUserUseCaseImpl(repo repository.UserRepository, util util.AuthUtil) *UserUseCaseImpl {
	return &UserUseCaseImpl{repo: repo, util: util}
}

func (uc *UserUseCaseImpl) AddAdmin(ctx context.Context, admin entity.User) (*entity.User, error) {
	if user, err := uc.repo.FindByEmail(ctx, admin.Email); err == nil {
		return nil, apperror.NewAlreadyExist(user, "Email", admin.Email)
	}

	newPassword, err := uc.util.HashAndSalt(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = newPassword
	admin.UserRoleId = 2
	admin.IsVerified = true

	created, err := uc.repo.Create(ctx, admin)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *UserUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.User, error) {
	user, err := uc.repo.FindById(ctx, id)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(user, "Id", id)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	users, err := uc.repo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, totalPages, err := uc.repo.CountFindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(users)), int64(*param.PageId), users,
	)

	return paginatedItems, nil
}
