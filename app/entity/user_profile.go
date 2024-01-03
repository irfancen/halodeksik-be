package entity

import (
	"database/sql"
	"time"
)

type UserProfile struct {
	UserId       int64        `json:"user_id"`
	Name         string       `json:"name"`
	ProfilePhoto string       `json:"profile_photo"`
	DateOfBirth  time.Time    `json:"date_of_birth"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
}
