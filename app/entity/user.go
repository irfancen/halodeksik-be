package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64        `json:"id"`
	Email      string       `json:"email"`
	Password   string       `json:"password"`
	UserRoleID int64        `json:"user_role_id"`
	IsVerified bool         `json:"is_verified"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}
