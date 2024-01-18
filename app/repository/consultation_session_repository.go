package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type ConsultationSessionRepository interface {
	Create(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error)
	FindById(ctx context.Context, id int64) (*entity.ConsultationSession, error)
}

type ConsultationSessionRepositoryImpl struct {
	db *sql.DB
}

func NewConsultationSessionRepositoryImpl(db *sql.DB) *ConsultationSessionRepositoryImpl {
	return &ConsultationSessionRepositoryImpl{db: db}
}

func (repo *ConsultationSessionRepositoryImpl) Create(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error) {
	const create = `INSERT INTO consultation_sessions(user_id, doctor_id, consultation_session_status_id)
	VALUES ($1, $2, $3) RETURNING
	id, user_id, doctor_id, consultation_session_status_id, created_at, updated_at`

	row := repo.db.QueryRowContext(ctx, create, session.UserId, session.DoctorId, session.ConsultationSessionStatusId)
	var created entity.ConsultationSession
	err := row.Scan(&created.Id, &created.UserId, &created.DoctorId, &created.ConsultationSessionStatusId, &created.CreatedAt, &created.UpdatedAt)

	return &created, err
}

func (repo *ConsultationSessionRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	const findById = `
	SELECT consultation_sessions.id, user_id, doctor_id, consultation_session_status_id, 
	       consultation_sessions.created_at, consultation_sessions.updated_at,
	       consultation_session_statuses.name
	FROM consultation_sessions
	INNER JOIN consultation_session_statuses ON consultation_sessions.consultation_session_status_id = consultation_session_statuses.id 
	WHERE consultation_sessions.id = $1`

	row := repo.db.QueryRowContext(ctx, findById, id)
	var session entity.ConsultationSession
	var sessionStatus entity.ConsultationSessionStatus
	err := row.Scan(
		&session.Id, &session.UserId, &session.DoctorId, &session.ConsultationSessionStatusId,
		&session.CreatedAt, &session.UpdatedAt,
		&sessionStatus.Name,
	)
	session.ConsultationSessionStatus = &sessionStatus

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &session, err
}
