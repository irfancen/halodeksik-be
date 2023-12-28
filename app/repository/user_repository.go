package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	const create = `INSERT INTO users(email, password, user_role_id, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row := u.db.QueryRowContext(ctx, create,
		user.Email,
		user.Password,
		user.UserRoleID,
		user.IsVerified,
	)

	var createdUser entity.User

	err := row.Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.UserRoleID,
		&createdUser.IsVerified,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.DeletedAt,
	)

	return &createdUser, err
}

func NewUserRepository(db *sql.DB) UserRepository {
	repo := userRepository{db: db}
	return &repo
}
