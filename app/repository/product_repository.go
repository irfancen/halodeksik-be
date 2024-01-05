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

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (*entity.Product, error)
	FindById(ctx context.Context, id int64) (*entity.Product, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Product, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, product entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepositoryImpl(db *sql.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (repo *ProductRepositoryImpl) Create(ctx context.Context, product entity.Product) (*entity.Product, error) {
	const create = `INSERT INTO products
		(name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form,
 		unit_in_pack, selling_unit, weight, length, width, height, image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create,
		product.Name, product.GenericName, product.Content, product.ManufacturerId, product.Description, product.DrugClassificationId, product.ProductCategoryId,
		product.DrugForm, product.UnitInPack, product.SellingUnit, product.Weight, product.Length, product.Width, product.Height, product.Image,
	)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrProductUniqueConstraint
		}
		return nil, row.Err()
	}

	var created entity.Product
	err := row.Scan(
		&created.Id, &created.Name, &created.GenericName, &created.Content, &created.ManufacturerId, &created.Description, &created.DrugClassificationId, &created.ProductCategoryId, &created.DrugForm,
		&created.UnitInPack, &created.SellingUnit, &created.Weight, &created.Length, &created.Width, &created.Height, &created.Image, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)
	return &created, err
}

func (repo *ProductRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Product, error) {
	const getById = `
SELECT p.id, p.name, p.generic_name, p.content, p.manufacturer_id, p.description, p.drug_classification_id, p.product_category_id, p.drug_form, p.unit_in_pack, p.selling_unit, p.weight, p.length, p.width, p.height, p.image, p.created_at, p.updated_at, p.deleted_at, pc.name, m.name, dc.name
FROM products p
         INNER JOIN product_categories pc ON p.product_category_id = pc.id
         INNER JOIN manufacturers m ON p.manufacturer_id = m.id
         INNER JOIN drug_classifications dc ON p.drug_classification_id = dc.id
WHERE p.id = $1 AND p.deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		product            entity.Product
		productCategory    entity.ProductCategory
		manufacturer       entity.Manufacturer
		drugClassification entity.DrugClassification
	)
	err := row.Scan(
		&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm,
		&product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
		&productCategory.Name, &manufacturer.Name, &drugClassification.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	product.ProductCategory = &productCategory
	product.Manufacturer = &manufacturer
	product.DrugClassification = &drugClassification
	return &product, err
}

func (repo *ProductRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Product, error) {
	initQuery := `SELECT id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image FROM products WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Product{}, param)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Product, 0)
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm,
			&product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ProductRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM products WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Product{}, param, false)

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

func (repo *ProductRepositoryImpl) Update(ctx context.Context, product entity.Product) (*entity.Product, error) {
	const updateById = `
		UPDATE products
		SET name=$1, generic_name=$2, content=$3, manufacturer_id=$4, description=$5, drug_classification_id=$6, product_category_id=$7, drug_form=$8, 
			unit_in_pack=$9, selling_unit=$10, weight=$11, length=$12, width=$13, height=$14, image=$15, updated_at = now()
		WHERE id = $16
		RETURNING id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, created_at, updated_at, deleted_at
		`

	row := repo.db.QueryRowContext(ctx, updateById,
		product.Name, product.GenericName, product.Content, product.ManufacturerId, product.Description, product.DrugClassificationId, product.ProductCategoryId, product.DrugForm,
		product.UnitInPack, product.SellingUnit, product.Weight, product.Length, product.Width, product.Height, product.Image, product.Id,
	)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrProductUniqueConstraint
		}
		return nil, row.Err()
	}

	var updated entity.Product
	err := row.Scan(
		&updated.Id, &updated.Name, &updated.GenericName, &updated.Content, &updated.ManufacturerId, &updated.Description, &updated.DrugClassificationId, &updated.ProductCategoryId, &updated.DrugForm,
		&updated.UnitInPack, &updated.SellingUnit, &updated.Weight, &updated.Length, &updated.Width, &updated.Height, &updated.Image, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt,
	)
	return &updated, err
}

func (repo *ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteById = `
		UPDATE products
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		`

	_, err := repo.db.ExecContext(ctx, deleteById, id)
	return err
}
