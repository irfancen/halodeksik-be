package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type CartItemRepository interface {
	Create(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error)
	FindByUserIdAndProductId(ctx context.Context, userId int64, productId int64) (*entity.CartItem, error)
	FindAllByUserId(ctx context.Context, userId int64) ([]*entity.CartItem, error)
	Update(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error)
	Delete(ctx context.Context, userId int64, productIds []int64) error
}

type CartItemRepositoryImpl struct {
	db *sql.DB
}

func NewCartItemRepositoryImpl(db *sql.DB) *CartItemRepositoryImpl {
	return &CartItemRepositoryImpl{db: db}
}

func (repo *CartItemRepositoryImpl) Create(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error) {
	const create = `
	INSERT INTO cart_items(user_id, product_id, quantity)
	VALUES ($1, $2, $3) RETURNING id, user_id, product_id, quantity, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create, cartItem.UserId, cartItem.ProductId, cartItem.Quantity)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var created entity.CartItem
	err := row.Scan(
		&created.Id,
		&created.UserId,
		&created.ProductId,
		&created.Quantity,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.DeletedAt,
	)
	return &created, err
}

func (repo *CartItemRepositoryImpl) FindByUserIdAndProductId(ctx context.Context, userId int64, productId int64) (*entity.CartItem, error) {
	const findByUserIdAndProductId = `SELECT id, user_id, product_id, quantity, created_at, updated_at, deleted_at
		FROM cart_items
		WHERE user_id = $1 AND product_id = $2 AND deleted_at IS NULL LIMIT 1`

	row := repo.db.QueryRowContext(ctx, findByUserIdAndProductId, userId, productId)
	var found entity.CartItem
	err := row.Scan(
		&found.Id, &found.UserId, &found.ProductId, &found.Quantity, &found.CreatedAt, &found.UpdatedAt, &found.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &found, err
}

func (repo *CartItemRepositoryImpl) FindAllByUserId(ctx context.Context, userId int64) ([]*entity.CartItem, error) {
	const findAll = `
	SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, 
	       p.id, p.name, p.generic_name, p.content, p.manufacturer_id, 
	       p.description, p.drug_classification_id, p.product_category_id, p.drug_form, 
	       p.unit_in_pack, p.selling_unit, p.weight, p.length, p.width, p.height, p.image
	FROM cart_items ci INNER JOIN products p ON ci.product_id = p.id
	WHERE ci.user_id = $1 AND ci.deleted_at IS NULL AND p.deleted_at IS NULL ORDER BY ci.updated_at DESC`

	rows, err := repo.db.QueryContext(ctx, findAll, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cartItems := make([]*entity.CartItem, 0)
	for rows.Next() {
		var (
			cartItem entity.CartItem
			product  entity.Product
		)
		if err := rows.Scan(
			&cartItem.Id, &cartItem.UserId, &cartItem.ProductId, &cartItem.Quantity,
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId,
			&product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm,
			&product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width,
			&product.Height, &product.Image,
		); err != nil {
			return nil, err
		}
		cartItem.Product = &product
		cartItems = append(cartItems, &cartItem)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (repo *CartItemRepositoryImpl) Update(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error) {
	const update = `UPDATE cart_items SET quantity = $1, updated_at = now()
		WHERE user_id = $2 AND product_id = $3 AND deleted_at IS NULL
		RETURNING id, user_id, product_id, quantity, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, update, cartItem.Quantity, cartItem.UserId, cartItem.ProductId)
	var updated entity.CartItem
	err := row.Scan(
		&updated.Id,
		&updated.UserId,
		&updated.ProductId,
		&updated.Quantity,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)
	return &updated, err
}

func (repo *CartItemRepositoryImpl) Delete(ctx context.Context, userId int64, productIds []int64) error {
	const deleteCartItems = `UPDATE cart_items SET deleted_at = now() WHERE deleted_at IS NULL AND user_id = $1 AND product_id = ANY ($2::int[])`
	_, err := repo.db.ExecContext(ctx, deleteCartItems, userId, pq.Array(productIds))
	return err
}
