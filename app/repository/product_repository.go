package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"log"
	"strings"
)

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (*entity.Product, error)
	FindById(ctx context.Context, id int64) (*entity.Product, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Product, error)
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
 		unit_in_pack, selling_unit, weight, length, width, height, image, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, price, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create,
		product.Name, product.GenericName, product.Content, product.ManufacturerId, product.Description, product.DrugClassificationId, product.ProductCategoryId,
		product.DrugForm, product.UnitInPack, product.SellingUnit, product.Weight, product.Length, product.Width, product.Height, product.Image, product.Price,
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
		&created.UnitInPack, &created.SellingUnit, &created.Weight, &created.Length, &created.Width, &created.Height, &created.Image, &created.Price, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)
	return &created, err
}

func (repo *ProductRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Product, error) {
	var getById = `SELECT id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, price, created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err()!= nil {
		return nil, row.Err()
	}

	var product entity.Product
	err := row.Scan(
		&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm,
		&product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &product, err
}

func buildQuery(initQuery string, param *queryparamdto.GetAllParams) (string, []interface{}) {
	var query strings.Builder
	var values []interface{}

	query.WriteString(initQuery)

	if len(param.WhereClauses) > 0 {
		query.WriteString(appdb.AND + " ")
	}

	indexPreparedStatement := 0

	for index, whereClause := range param.WhereClauses {
		if whereClause.Condition == appdb.In {
			query.WriteString(fmt.Sprintf("%s %s (", whereClause.Column, whereClause.Condition))
			val := strings.Split(whereClause.Value.(string), ",")
			for idx, v := range val {
				indexPreparedStatement++
				query.WriteString(fmt.Sprintf("$%d", indexPreparedStatement))
				if idx != len(val) - 1 {
					query.WriteString(",")
				}
				values = append(values, v)
			}
			query.WriteString(string(") " + whereClause.Logic))
			continue
		}

		indexPreparedStatement++
		query.WriteString(fmt.Sprintf("%s %s $%d %s ", whereClause.Column, whereClause.Condition, indexPreparedStatement, whereClause.Logic))

		if index != len(param.WhereClauses)-1 && util.IsEmptyString(string(whereClause.Logic)) {
			query.WriteString(appdb.AND + " ")
		}

		values = append(values, whereClause.Value)
	}

	query.WriteString(" ORDER BY ")

	for _, sortClause := range param.SortClauses {
		query.WriteString(fmt.Sprintf("%s %s,", sortClause.Column, sortClause.Order))
	}

	query.WriteString(` id ASC `)

	if param.PageId != nil && param.PageSize != nil {
		size := *param.PageSize
		offset := (*param.PageId - 1) * size

		query.WriteString(fmt.Sprintf("LIMIT $%d OFFSET $%d", indexPreparedStatement+1, indexPreparedStatement+2))
		values = append(values, size, offset)
	}
	log.Println(query.String())
	return query.String(), values
}

func (repo *ProductRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Product, error) {
	initQuery := `SELECT id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, price FROM products WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, param)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Product, 0)
	for rows.Next() {
		var product *entity.Product
		if err := rows.Scan(
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm,
			&product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ProductRepositoryImpl) Update(ctx context.Context, product entity.Product) (*entity.Product, error) {
	const updateById = `
		UPDATE products
		SET name=$1, generic_name=$2, content=$3, manufacturer_id=$4, description=$5, drug_classification_id=$6, product_category_id=$7, drug_form=$8, 
			unit_in_pack=$9, selling_unit=$10, weight=$11, length=$12, width=$13, height=$14, image=$15, price=$16, updated_at = now()
		WHERE id = $17
		RETURNING id, name, generic_name, content, manufacturer_id, description, drug_classification_id, product_category_id, drug_form, unit_in_pack, selling_unit, weight, length, width, height, image, price, created_at, updated_at, deleted_at
		`

	row := repo.db.QueryRowContext(ctx, updateById,
		product.Name, product.GenericName, product.Content, product.ManufacturerId, product.Description, product.DrugClassificationId, product.ProductCategoryId, product.DrugForm,
		product.UnitInPack, product.SellingUnit, product.Weight, product.Length, product.Width, product.Height, product.Image, product.Price, product.Id,
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
		&updated.UnitInPack, &updated.SellingUnit, &updated.Weight, &updated.Length, &updated.Width, &updated.Height, &updated.Image, &updated.Price, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt,
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
