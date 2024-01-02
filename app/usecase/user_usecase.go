package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type UserUseCase interface {
	AddAdmin(ctx context.Context, admin entity.User) (*entity.User, error)
	GetById(ctx context.Context, id int64) (*entity.User, error)
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
