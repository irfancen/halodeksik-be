package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
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
