package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type PharmacyRepository interface {
	Create(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
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
