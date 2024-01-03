package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type VerifyTokenRepository interface {
	CreateVerifyToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error)
	FindTokenByToken(ctx context.Context, token string) (*entity.VerificationToken, error)
	FindTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error)
	DeactivateToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error)
}

type VerifyTokenRepositoryImpl struct {
	db *sql.DB
}

func (repo *VerifyTokenRepositoryImpl) DeactivateToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error) {
	const deleteToken = `
	DELETE FROM verification_tokens WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, deleteToken,
		token.Token,
	)

	return &token, row.Err()
}

func (repo *VerifyTokenRepositoryImpl) FindTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error) {
	const getActiveVerifyTokenByEmail = `
	SELECT id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at FROM verification_tokens
	WHERE email = $1
	`

	row := repo.db.QueryRowContext(ctx, getActiveVerifyTokenByEmail,
		email,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return &createdToken, err

}

func (repo *VerifyTokenRepositoryImpl) FindTokenByToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	const getTokenByToken = `
	SELECT id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at FROM verification_tokens
	WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, getTokenByToken,
		token,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return &createdToken, err

}

func (repo *VerifyTokenRepositoryImpl) CreateVerifyToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error) {

	const createVerifyToken = `
	INSERT INTO verification_tokens(token, is_valid, expired_at, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at
	`
	row := repo.db.QueryRowContext(ctx, createVerifyToken,
		token.Token,
		token.IsValid,
		token.ExpiredAt,
		token.Email,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	return &createdToken, err

}

func NewVerifyTokenRepository(db *sql.DB) *VerifyTokenRepositoryImpl {
	repo := VerifyTokenRepositoryImpl{db: db}
	return &repo
}
