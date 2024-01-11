package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type AddressAreaRepository interface {
	FindAllProvince(ctx context.Context) ([]*entity.Province, error)
	FindAllCities(ctx context.Context) ([]*entity.City, error)
}

type AddressAreaRepositoryImpl struct {
	db *sql.DB
}

func NewAddressAreaRepositoryImpl(db *sql.DB) *AddressAreaRepositoryImpl {
	return &AddressAreaRepositoryImpl{db: db}
}

func (repo *AddressAreaRepositoryImpl) FindAllProvince(ctx context.Context) ([]*entity.Province, error) {
	query := `SELECT id, name FROM provinces WHERE deleted_at IS NULL ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Province, 0)
	for rows.Next() {
		var province entity.Province
		if err := rows.Scan(
			&province.Id, &province.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, &province)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *AddressAreaRepositoryImpl) FindAllCities(ctx context.Context) ([]*entity.City, error) {
	query := `SELECT id, name, province_id FROM cities WHERE deleted_at IS NULL ORDER BY province_id, name`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.City, 0)
	for rows.Next() {
		var city entity.City
		if err := rows.Scan(
			&city.Id, &city.Name, &city.ProvinceId,
		); err != nil {
			return nil, err
		}
		items = append(items, &city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
