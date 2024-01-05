package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"

	"github.com/jackc/pgx/v5/pgconn"
)

type PharmacyProductRepository interface {
	Create(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
	FindAllJoinProducts(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
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

func (repo *PharmacyProductRepositoryImpl) FindAllJoinProducts(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error) {
	initQuery := `
	SELECT pharmacy_products.id, pharmacy_products.pharmacy_id, pharmacy_products.product_id, pharmacy_products.is_active, pharmacy_products.price, pharmacy_products.stock,
		products.id, products.name, products.generic_name, products.content, products.manufacturer_id, products.description, products.drug_classification_id, products.product_category_id, products.drug_form, products.unit_in_pack, products.selling_unit, products.weight, products.length, products.width, products.height, products.image,
		product_categories.name, manufacturers.name, drug_classifications.name
	FROM pharmacy_products
			INNER JOIN products ON pharmacy_products.product_id = products.id
			INNER JOIN product_categories ON products.product_category_id = product_categories.id
	        INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	        INNER JOIN drug_classifications ON products.drug_classification_id = drug_classifications.id
	WHERE pharmacy_products.deleted_at IS NULL AND products.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.PharmacyProduct, 0)
	for rows.Next() {
		var (
			pharmacyProduct    entity.PharmacyProduct
			product            entity.Product
			productCategory    entity.ProductCategory
			manufacturer       entity.Manufacturer
			drugClassification entity.DrugClassification
		)
		if err := rows.Scan(
			&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId, &pharmacyProduct.IsActive, &pharmacyProduct.Price, &pharmacyProduct.Stock,
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm, &product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image,
			&productCategory.Name, &manufacturer.Name, &drugClassification.Name,
		); err != nil {
			return nil, err
		}
		product.ProductCategory = &productCategory
		product.Manufacturer = &manufacturer
		product.DrugClassification = &drugClassification
		pharmacyProduct.Product = &product
		items = append(items, &pharmacyProduct)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *PharmacyProductRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(pharmacy_products.id) FROM pharmacy_products
	INNER JOIN products ON pharmacy_products.product_id = products.id
	INNER JOIN product_categories ON products.product_category_id = product_categories.id
	INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	INNER JOIN drug_classifications ON products.drug_classification_id = drug_classifications.id
	WHERE pharmacy_products.deleted_at IS NULL AND products.deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, false)

	var totalItems int64

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&totalItems,
		); err != nil {
			return totalItems, err
		}
	}

	if err := rows.Err(); err != nil {
		return totalItems, err
	}
	return totalItems, nil
}
