package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type ProductStockMutationRequestRepository interface {
	Create(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
	FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutationRequest, error)
	CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
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

func (repo *ProductStockMutationRequestRepositoryImpl) FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutationRequest, error) {
	initQuery := `
SELECT product_stock_mutation_requests.id,
       product_stock_mutation_requests.pharmacy_product_origin_id,
       product_stock_mutation_requests.pharmacy_product_dest_id,
       product_stock_mutation_requests.stock,
       product_stock_mutation_requests.product_stock_mutation_request_status_id,
       product_stock_mutation_requests.created_at AS request_date,
       p.name AS pharmacy_name,
       pr.name, pr.generic_name, pr.content, pr.manufacturer_id,
       m.name,
       psmrs.name AS status
FROM product_stock_mutation_requests
         INNER JOIN pharmacy_products ppo ON ppo.id = product_stock_mutation_requests.pharmacy_product_origin_id
         INNER JOIN pharmacy_products ppd ON ppd.id = product_stock_mutation_requests.pharmacy_product_dest_id
         INNER JOIN pharmacies p ON ppd.pharmacy_id = p.id
         INNER JOIN products pr ON ppo.product_id = pr.id
    	 INNER JOIN manufacturers m ON pr.manufacturer_id = m.id
         INNER JOIN product_stock_mutation_request_statuses psmrs
                    ON product_stock_mutation_requests.product_stock_mutation_request_status_id = psmrs.id
WHERE product_stock_mutation_requests.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.ProductStockMutationRequest{}, param, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.ProductStockMutationRequest, 0)
	for rows.Next() {
		var (
			mutationRequest       entity.ProductStockMutationRequest
			mutationRequestStatus entity.ProductStockMutationRequestStatus
			pharmacyProduct       entity.PharmacyProduct
			pharmacy              entity.Pharmacy
			product               entity.Product
			manufacturer          entity.Manufacturer
		)
		if err := rows.Scan(
			&mutationRequest.Id, &mutationRequest.PharmacyProductOriginId, &mutationRequest.PharmacyProductDestId, &mutationRequest.Stock, &mutationRequest.ProductStockMutationRequestStatusId, &mutationRequest.CreatedAt,
			&pharmacy.Name,
			&product.Name, &product.GenericName, &product.Content, &product.ManufacturerId,
			&manufacturer.Name,
			&mutationRequestStatus.Name,
		); err != nil {
			return nil, err
		}
		product.Manufacturer = &manufacturer
		pharmacyProduct.Pharmacy = &pharmacy
		pharmacyProduct.Product = &product
		mutationRequest.PharmacyProductDest = &pharmacyProduct
		mutationRequest.ProductStockMutationRequestStatus = &mutationRequestStatus
		items = append(items, &mutationRequest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ProductStockMutationRequestRepositoryImpl) CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `
	SELECT count(product_stock_mutation_requests.id)
FROM product_stock_mutation_requests
         INNER JOIN pharmacy_products ppo ON ppo.id = product_stock_mutation_requests.pharmacy_product_origin_id
         INNER JOIN pharmacy_products ppd ON ppd.id = product_stock_mutation_requests.pharmacy_product_dest_id
         INNER JOIN pharmacies p ON ppd.pharmacy_id = p.id
         INNER JOIN products pr ON ppo.product_id = pr.id
         INNER JOIN product_stock_mutation_request_statuses psmrs
                    ON product_stock_mutation_requests.product_stock_mutation_request_status_id = psmrs.id
WHERE product_stock_mutation_requests.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.ProductStockMutationRequest{}, param, false)

	var (
		totalItems int64
		temp       int64
	)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&temp,
		); err != nil {
			return totalItems, err
		}
		totalItems++
	}

	if err := rows.Err(); err != nil {
		return totalItems, err
	}
	return totalItems, nil
}
