package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type PharmacyRepository interface {
	Create(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Pharmacy, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, int64, error)
}

type PharmacyRepositoryImpl struct {
	db *sql.DB
}

func NewPharmacyRepository(db *sql.DB) *PharmacyRepositoryImpl {
	return &PharmacyRepositoryImpl{db: db}
}

func (repo *PharmacyRepositoryImpl) Create(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	const create = `INSERT INTO pharmacies(name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id, created_at, updated_at, deleted_at
`

	row := repo.db.QueryRowContext(ctx, create,
		pharmacy.Name,
		pharmacy.Address, pharmacy.SubDistrict, pharmacy.District, pharmacy.City, pharmacy.Province, pharmacy.PostalCode,
		pharmacy.Latitude, pharmacy.Longitude,
		pharmacy.PharmacistName, pharmacy.PharmacistLicenseNo, pharmacy.PharmacistPhoneNo,
		pharmacy.OperationalHours, pharmacy.OperationalDays, pharmacy.PharmacyAdminId,
	)
	var created entity.Pharmacy
	err := row.Scan(
		&created.Id, &created.Name,
		&created.Address, &created.SubDistrict, &created.District, &created.City, &created.Province, &created.PostalCode,
		&created.Latitude, &created.Longitude,
		&created.PharmacistName, &created.PharmacistLicenseNo, &created.PharmacistPhoneNo,
		&created.OperationalHours, &created.OperationalDays, &created.PharmacyAdminId,
		&created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)
	return &created, err
}

func (repo *PharmacyRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Pharmacy, error) {
	initQuery := `SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id FROM pharmacies WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, param)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Pharmacy, 0)
	for rows.Next() {
		var pharmacy entity.Pharmacy
		if err := rows.Scan(
			&pharmacy.Id, &pharmacy.Name,
			&pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.City, &pharmacy.Province, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude,
			&pharmacy.PharmacistName, &pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo,
			&pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
		); err != nil {
			return nil, err
		}
		items = append(items, &pharmacy)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *PharmacyRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, int64, error) {
	initQuery := `SELECT count(id) FROM pharmacies WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, param, false)

	var (
		totalItems int64
		totalPages int64
	)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, totalPages, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&totalItems,
		); err != nil {
			return totalItems, totalPages, err
		}
	}
	totalPages = totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	if err := rows.Err(); err != nil {
		return totalItems, totalPages, err
	}

	return totalItems, totalPages, nil
}
