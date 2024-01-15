package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type ProductStockMutationRequestRepository interface {
	Create(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
}

type ProductStockMutationRequestRepositoryImpl struct {
	db *sql.DB
}

func NewProductStockMutationRequestRepositoryImpl(db *sql.DB) *ProductStockMutationRequestRepositoryImpl {
	return &ProductStockMutationRequestRepositoryImpl{db: db}
}

func (repo *ProductStockMutationRequestRepositoryImpl) Create(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	create := `INSERT INTO product_stock_mutation_requests(pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id)
VALUES ($1, $2, $3, $4)
RETURNING id, pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id, created_at`

	row := repo.db.QueryRowContext(ctx, create,
		mutationRequest.PharmacyProductOriginId,
		mutationRequest.PharmacyProductDestId,
		mutationRequest.Stock,
		mutationRequest.ProductStockMutationRequestStatusId,
	)

	var created entity.ProductStockMutationRequest
	if err := row.Scan(
		&created.Id,
		&created.PharmacyProductOriginId,
		&created.PharmacyProductDestId,
		&created.Stock,
		&created.ProductStockMutationRequestStatusId,
		&created.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &created, nil
}
