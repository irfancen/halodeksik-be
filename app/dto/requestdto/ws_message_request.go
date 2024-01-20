package requestdto

import (
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
)

type WsConsultationMessage struct {
	IsTyping   bool   `json:"is_typing"`
	Message    string `json:"message"`
	Attachment string `json:"attachment"`
	SenderId   int64  `json:"sender_id"`
	SessionId  int64  `json:"session_id"`
}

func (r *WsConsultationMessage) ToConsultationMessage() *entity.ConsultationMessage {
	return &entity.ConsultationMessage{
		Message:    appdb.NewSqlNullString(r.Message),
		Attachment: appdb.NewSqlNullString(r.Attachment),
		SenderId:   appdb.NewSqlNullInt64(r.SenderId),
		SessionId:  appdb.NewSqlNullInt64(r.SessionId),
	}
}
