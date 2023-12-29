package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type ProductCategory struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *ProductCategory) ToResponse() *responsedto.ProductCategoryResponse {
	return &responsedto.ProductCategoryResponse{
		Id:   e.ID,
		Name: e.Name,
	}
}
