package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type DoctorSpecializationRepository interface {
	Create(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	FindById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error)
	FindAllWithoutParams(ctx context.Context) ([]*entity.DoctorSpecialization, error)
}

type DoctorSpecializationRepositoryImpl struct {
	db *sql.DB
}

func NewDoctorSpecializationRepositoryImpl(db *sql.DB) *DoctorSpecializationRepositoryImpl {
	return &DoctorSpecializationRepositoryImpl{db: db}
}

func (repo *DoctorSpecializationRepositoryImpl) Create(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	const create = `INSERT INTO doctor_specializations(name, image)
	VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create, specialization.Name, specialization.Image)
	var created entity.DoctorSpecialization
	err := row.Scan(
		&created.Id, &created.Name, &created.Image, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)

	return &created, err
}

func (repo *DoctorSpecializationRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error) {
	const findById = `SELECT id, name, image, created_at, updated_at, deleted_at FROM doctor_specializations WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, findById, id)
	var specialization entity.DoctorSpecialization
	err := row.Scan(
		&specialization.Id, &specialization.Name, &specialization.Image, &specialization.CreatedAt, &specialization.UpdatedAt, &specialization.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &specialization, err
}

func (repo *DoctorSpecializationRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.DoctorSpecialization, error) {
	const getDoctorSpecs = `
	SELECT id, name, created_at, updated_at, deleted_at FROM doctor_specializations
	ORDER BY id
	`

	rows, err := repo.db.QueryContext(ctx, getDoctorSpecs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.DoctorSpecialization, 0)
	for rows.Next() {
		var ds entity.DoctorSpecialization
		if err := rows.Scan(
			&ds.Id, &ds.Name, &ds.CreatedAt, &ds.UpdatedAt, &ds.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &ds)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
