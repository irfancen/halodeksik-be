package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type ConsultationSession struct {
	Id                          int64        `json:"id"`
	UserId                      int64        `json:"user_id"`
	DoctorId                    int64        `json:"doctor_id"`
	ConsultationSessionStatusId int64        `json:"consultation_session_status_id"`
	CreatedAt                   time.Time    `json:"created_at"`
	UpdatedAt                   time.Time    `json:"updated_at"`
	DeletedAt                   sql.NullTime `json:"deleted_at"`
	ConsultationSessionStatus   *ConsultationSessionStatus
}

func (e *ConsultationSession) ToResponse() *responsedto.ConsultationSessionResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ConsultationSessionResponse{
		Id:                          e.Id,
		UserId:                      e.UserId,
		DoctorId:                    e.DoctorId,
		ConsultationSessionStatusId: e.ConsultationSessionStatusId,
		CreatedAt:                   e.CreatedAt,
		UpdatedAt:                   e.UpdatedAt,
		ConsultationSessionStatus:   e.ConsultationSessionStatus.ToResponse(),
	}
}
