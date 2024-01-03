package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type DrugClassification struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *DrugClassification) ToResponse() *responsedto.DrugClassificationResponse {
	if e == nil {
		return nil
	}
	return &responsedto.DrugClassificationResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}
