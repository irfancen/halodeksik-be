package entity

import (
	"database/sql"
	"time"
)

type ProductStockMutationRequestStatus struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
