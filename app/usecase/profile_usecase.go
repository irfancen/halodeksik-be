package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProfileUseCase interface {
	GetUserProfileByUserId(ctx context.Context, userId int64) (*entity.UserProfile, error)
	GetDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.DoctorProfile, error)
	UpdateUserProfile(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error)
	UpdateDoctorProfile(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error)
}

type ProfileUseCaseImpl struct {
	repo repository.ProfileRepository
}

func (uc *ProfileUseCaseImpl) UpdateUserProfile(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error) {
	panic("implement me")
}

func (uc *ProfileUseCaseImpl) UpdateDoctorProfile(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error) {
	//TODO implement me
	panic("implement me")
}

func NewProfileUseCaseImpl(repo repository.ProfileRepository) *ProfileUseCaseImpl {
	return &ProfileUseCaseImpl{repo: repo}
}

func (uc *ProfileUseCaseImpl) GetUserProfileByUserId(ctx context.Context, userId int64) (*entity.UserProfile, error) {
	profile, err := uc.repo.FindUserProfileByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(profile, "Id", userId)
		}
		return nil, err
	}
	return profile, nil
}

func (uc *ProfileUseCaseImpl) GetDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.DoctorProfile, error) {
	profile, err := uc.repo.FindDoctorProfileByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(profile, "Id", userId)
		}
		return nil, err
	}
	return profile, nil
}
