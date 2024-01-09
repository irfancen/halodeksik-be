package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type DoctorSpecializationRepository interface {
	FindAllWithoutParams(ctx context.Context) ([]*entity.DoctorSpecialization, error)
}

type DoctorSpecializationRepositoryImpl struct {
	db *sql.DB
}

func NewDoctorSpecializationRepositoryImpl(db *sql.DB) *DoctorSpecializationRepositoryImpl {
	return &DoctorSpecializationRepositoryImpl{db: db}
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
