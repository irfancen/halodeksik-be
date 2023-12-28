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

type DoctorProfile struct {
	UserID            int64        `json:"user_id"`
	Name              string       `json:"name"`
	ProfilePhoto      string       `json:"profile_photo"`
	StartingYear      int32        `json:"starting_year"`
	DoctorCertificate string       `json:"doctor_certificate"`
	Specialization    string       `json:"specialization"`
	ConsultationFee   string       `json:"consultation_fee"`
	IsOnline          bool         `json:"is_online"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
}

type UserProfile struct {
	UserID       int64        `json:"user_id"`
	Name         string       `json:"name"`
	ProfilePhoto string       `json:"profile_photo"`
	DateOfBirth  time.Time    `json:"date_of_birth"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
}

type UserRole struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
