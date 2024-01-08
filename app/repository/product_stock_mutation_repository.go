package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/entity"
)

type ProductStockMutationRepository interface {
	Create(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error)
}

type ProductStockMutationRepositoryImpl struct {
	db *sql.DB
}

func NewProductStockMutationRepositoryImpl(db *sql.DB) *ProductStockMutationRepositoryImpl {
	return &ProductStockMutationRepositoryImpl{db: db}
}

func (repo *ProductStockMutationRepositoryImpl) Create(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	create := `INSERT INTO product_stock_mutations(pharmacy_product_id, product_stock_mutation_type_id, stock)
VALUES ($1, $2, $3)
RETURNING id, pharmacy_product_id, product_stock_mutation_type_id, stock
`
	row1 := tx.QueryRowContext(ctx, create,
		stockMutation.PharmacyProductId, stockMutation.ProductStockMutationTypeId, stockMutation.Stock,
	)

	var created entity.ProductStockMutation
	if err := row1.Scan(
		&created.Id, &created.PharmacyProductId, &created.ProductStockMutationTypeId, &created.Stock,
	); err != nil {
		return nil, err
	}

	stock := stockMutation.Stock
	if stockMutation.ProductStockMutationTypeId == appconstant.StockMutationTypeReduction {
		stock = 0 - stock
	}

	updateStock := `UPDATE pharmacy_products
SET stock = stock + $1
WHERE id = $2
	`

	if _, err := tx.ExecContext(ctx, updateStock,
		stock,
		stockMutation.PharmacyProductId,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &created, err
}
