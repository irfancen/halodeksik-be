package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
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
	EditAdmin(ctx context.Context, id int64, user entity.User) (*entity.User, error)
	RemoveAdmin(ctx context.Context, id int64) error
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

	currentUserId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	currentUserRoleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)
	if currentUserRoleId != appconstant.UserRoleIdAdmin && currentUserId != user.Id {
		return nil, apperror.ErrForbiddenViewEntity
	}

	return user, nil
}

func (uc *UserUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	users, err := uc.repo.FindAll(ctx, param)
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

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(users)), int64(*param.PageId), users,
	)

	return paginatedItems, nil
}

func (uc *UserUseCaseImpl) EditAdmin(ctx context.Context, id int64, user entity.User) (*entity.User, error) {
	userdb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if userdb.UserRoleId != appconstant.UserRoleIdPharmacyAdmin {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if user.Email == "" && user.Password == "" {
		return userdb, nil
	}

	if user.Email != "" {
		userdb.Email = user.Email
	}

	if user.Password != "" {
		newPassword, err := uc.util.HashAndSalt(user.Password)
		if err != nil {
			return nil, err
		}
		userdb.Password = newPassword
	}

	updated, err := uc.repo.Update(ctx, *userdb)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *UserUseCaseImpl) RemoveAdmin(ctx context.Context, id int64) error {
	userdb, err := uc.GetById(ctx, id)
	if err != nil {
		return err
	}

	if userdb.UserRoleId != appconstant.UserRoleIdPharmacyAdmin {
		return apperror.ErrForbiddenModifyEntity
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
