package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type ProductCategoryRepository interface {
	FindAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error)
}

type ProductCategoryRepositoryImpl struct {
	db *sql.DB
}

func NewProductCategoryRepositoryImpl(db *sql.DB) *ProductCategoryRepositoryImpl {
	return &ProductCategoryRepositoryImpl{db: db}
}

func (repo *ProductCategoryRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error) {
	const findAll = `
		SELECT id, name, created_at, updated_at, deleted_at FROM product_categories
		`

	rows, err := repo.db.QueryContext(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*entity.ProductCategory, 0)
	for rows.Next() {
		var category entity.ProductCategory
		if err := rows.Scan(
			&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
