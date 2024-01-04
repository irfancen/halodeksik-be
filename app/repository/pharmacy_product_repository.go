package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"

	"github.com/jackc/pgx/v5/pgconn"
)

type PharmacyProductRepository interface {
	Create(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
}

type PharmacyProductRepositoryImpl struct {
	db *sql.DB
}

func NewPharmacyProductRepository(db *sql.DB) *PharmacyProductRepositoryImpl {
	return &PharmacyProductRepositoryImpl{db: db}
}

func (repo *PharmacyProductRepositoryImpl) Create(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	const create = `INSERT INTO pharmacy_products(pharmacy_id, product_id, is_active, price, stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, pharmacy_id, product_id, is_active, price, stock
`

	row := repo.db.QueryRowContext(ctx, create,
		pharmacyProduct.PharmacyId,
		pharmacyProduct.ProductId,
		pharmacyProduct.IsActive,
		pharmacyProduct.Price.String(),
		pharmacyProduct.Stock,
	)

	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrPharmacyProductUniqueConstraint
		}
		return nil, row.Err()
	}

	var created entity.PharmacyProduct
	err := row.Scan(
		&created.Id, &created.PharmacyId, &created.ProductId,
		&created.IsActive, &created.Price, &created.Stock,
	)
	return &created, err
}
