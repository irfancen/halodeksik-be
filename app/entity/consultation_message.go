package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
)

type ConsultationMessage struct {
	Id         sql.NullInt64  `json:"id"`
	SessionId  sql.NullInt64  `json:"session_id"`
	SenderId   sql.NullInt64  `json:"sender_id"`
	Message    sql.NullString `json:"message"`
	Attachment sql.NullString `json:"attachment"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
}

func (e *ConsultationMessage) ToResponse() *responsedto.ConsultationMessageResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ConsultationMessageResponse{
		Id:         e.Id.Int64,
		SessionId:  e.SessionId.Int64,
		SenderId:   e.SenderId.Int64,
		Message:    e.Message.String,
		Attachment: e.Attachment.String,
		CreatedAt:  e.CreatedAt.Time,
		UpdatedAt:  e.UpdatedAt.Time,
	}
}
