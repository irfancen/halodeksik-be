package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
	FindAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error)
}

type ManufacturerRepositoryImpl struct {
	db *sql.DB
}

func NewManufacturerRepositoryImpl(db *sql.DB) *ManufacturerRepositoryImpl {
	return &ManufacturerRepositoryImpl{db: db}
}

func (repo *ManufacturerRepositoryImpl) Create(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error) {
	const create = `INSERT INTO manufacturers(name, image)
	VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create, manufacturer.Name, manufacturer.Image)
	var created entity.Manufacturer
	err := row.Scan(
		&created.Id, &created.Name, &created.Image, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)

	return &created, err
}

func (repo *ManufacturerRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error) {
	const getAllWithoutParams = `
		SELECT id, name, created_at, updated_at, deleted_at FROM manufacturers
		`

	rows, err := repo.db.QueryContext(ctx, getAllWithoutParams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Manufacturer, 0)
	for rows.Next() {
		var manufacturer entity.Manufacturer
		if err := rows.Scan(
			&manufacturer.Id, &manufacturer.Name, &manufacturer.CreatedAt, &manufacturer.UpdatedAt, &manufacturer.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &manufacturer)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
