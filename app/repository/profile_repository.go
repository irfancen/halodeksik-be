package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type ProfileRepository interface {
	FindUserProfileByUserId(ctx context.Context, userId int64) (*entity.UserProfile, error)
	FindDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.DoctorProfile, error)
	UpdateUserProfileByUserId(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error)
	UpdateDoctorProfileByUserId(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error)
}

type ProfileRepositoryImpl struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepositoryImpl {
	repo := ProfileRepositoryImpl{db: db}
	return &repo
}

func (repo *ProfileRepositoryImpl) FindUserProfileByUserId(ctx context.Context, userId int64) (*entity.UserProfile, error) {
	const getUserProfileByUserId = `
	SELECT user_id, name, profile_photo, date_of_birth, created_at, updated_at, deleted_at FROM user_profiles
	WHERE user_id = $1
	`

	row := repo.db.QueryRowContext(ctx, getUserProfileByUserId,
		userId,
	)

	var profile entity.UserProfile

	err := row.Scan(
		&profile.UserId,
		&profile.Name,
		&profile.ProfilePhoto,
		&profile.DateOfBirth,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (repo *ProfileRepositoryImpl) FindDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.DoctorProfile, error) {
	const getDoctorProfileByUserId = `
	SELECT user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online, created_at, updated_at, deleted_at FROM doctor_profiles
	WHERE user_id = $1
	`

	row := repo.db.QueryRowContext(ctx, getDoctorProfileByUserId,
		userId,
	)

	var profile entity.DoctorProfile

	err := row.Scan(
		&profile.UserId, &profile.Name,
		&profile.ProfilePhoto,
		&profile.StartingYear,
		&profile.DoctorCertificate,
		&profile.DoctorSpecializationId,
		&profile.ConsultationFee,
		&profile.IsOnline,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (repo *ProfileRepositoryImpl) UpdateUserProfileByUserId(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error) {
	const updateUserProfileByUserId = `
	UPDATE user_profiles
	SET name = $1, profile_photo = $2, date_of_birth = $3, updated_at = now() WHERE user_id = $4
	RETURNING user_id, name, profile_photo, date_of_birth, created_at, updated_at, deleted_at
	`

	row := repo.db.QueryRowContext(ctx, updateUserProfileByUserId,
		profile.Name, profile.ProfilePhoto, profile.DateOfBirth, profile.UserId,
	)

	var updatedProfile entity.UserProfile

	err := row.Scan(
		&updatedProfile.UserId,
		&updatedProfile.Name,
		&updatedProfile.ProfilePhoto,
		&updatedProfile.DateOfBirth,
		&updatedProfile.CreatedAt,
		&updatedProfile.UpdatedAt,
		&updatedProfile.DeletedAt,
	)

	return &updatedProfile, err

}

func (repo *ProfileRepositoryImpl) UpdateDoctorProfileByUserId(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error) {
	const updateDoctorProfileByUserId = `
	UPDATE doctor_profiles
	SET name = $1, profile_photo = $2, starting_year = $3, doctor_certificate = $4, doctor_specialization_id = $5, consultation_fee = $6 WHERE user_id = $7
	RETURNING user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online, created_at, updated_at, deleted_at
	`
	row := repo.db.QueryRowContext(ctx, updateDoctorProfileByUserId,
		profile.Name, profile.ProfilePhoto, profile.StartingYear, profile.DoctorCertificate, profile.DoctorSpecializationId,
		profile.ConsultationFee, profile.UserId,
	)

	var updatedProfile entity.DoctorProfile

	err := row.Scan(
		&updatedProfile.UserId,
		&updatedProfile.Name,
		&updatedProfile.ProfilePhoto,
		&updatedProfile.StartingYear,
		&updatedProfile.DoctorCertificate,
		&updatedProfile.DoctorSpecializationId,
		&updatedProfile.ConsultationFee,
		&updatedProfile.IsOnline,
		&updatedProfile.CreatedAt,
		&updatedProfile.UpdatedAt,
		&updatedProfile.DeletedAt,
	)

	return &updatedProfile, err

}
