package entity

import (
	"database/sql"
	"time"
)

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
