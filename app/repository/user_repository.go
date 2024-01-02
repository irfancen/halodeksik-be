package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, int64, error)
	Update(ctx context.Context, user entity.User) (*entity.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	repo := UserRepositoryImpl{db: db}
	return &repo
}

func (repo *UserRepositoryImpl) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	const create = `INSERT INTO users(email, password, user_role_id, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
	)

	var createdUser entity.User

	err := row.Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.UserRoleId,
		&createdUser.IsVerified,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.DeletedAt,
	)

	return &createdUser, err
}

func (repo *UserRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.User, error) {
	const getById = `SELECT id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at
	FROM users WHERE id = $1
`

	row := repo.db.QueryRowContext(ctx, getById, id)
	var user entity.User
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.UserRoleId,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (repo *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	const getUserByEmail = `SELECT id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at FROM users
WHERE email = $1
`

	row := repo.db.QueryRowContext(ctx, getUserByEmail, email)
	var user entity.User
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.UserRoleId,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (repo *UserRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error) {
	initQuery := `SELECT id, email, user_role_id, is_verified FROM users WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, param)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.Id, &user.Email, &user.UserRoleId, &user.IsVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *UserRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, int64, error) {
	initQuery := `SELECT count(id) FROM users WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, param, false)

	var (
		totalItems int64
		totalPages int64
	)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, totalPages, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&totalItems,
		); err != nil {
			return totalItems, totalPages, err
		}
	}
	totalPages = totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	if err := rows.Err(); err != nil {
		return totalItems, totalPages, err
	}

	return totalItems, totalPages, nil
}

func (repo *UserRepositoryImpl) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	const updateById = `UPDATE users
SET email = $1, password = $2, user_role_id = $3, is_verified = $4, updated_at = now()
WHERE id = $5
RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at
`

	row := repo.db.QueryRowContext(ctx, updateById,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
		user.Id,
	)
	var updated entity.User
	err := row.Scan(
		&updated.Id,
		&updated.Email,
		&updated.Password,
		&updated.UserRoleId,
		&updated.IsVerified,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)
	return &updated, err
}
