package responsedto

import "time"

type WsConsultationMessage struct {
	IsTyping   bool      `json:"is_typing"`
	Message    string    `json:"message"`
	Attachment string    `json:"attachment"`
	CreatedAt  time.Time `json:"created_at"`
	SenderId   int64     `json:"sender_id"`
	SessionId  int64     `json:"session_id"`
}
