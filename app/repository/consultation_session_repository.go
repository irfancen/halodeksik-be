package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type ConsultationSessionRepository interface {
	Create(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error)
	FindById(ctx context.Context, id int64) (*entity.ConsultationSession, error)
	FindByUserIdAndDoctorId(ctx context.Context, userId, doctorId int64) (*entity.ConsultationSession, error)
	FindAllByUserIdOrDoctorId(ctx context.Context, userIdOrDoctorId int64, param *queryparamdto.GetAllParams) ([]*entity.ConsultationSession, error)
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

func (repo *ConsultationSessionRepositoryImpl) FindByUserIdAndDoctorId(ctx context.Context, userId, doctorId int64) (*entity.ConsultationSession, error) {
	const findByUserIdAndDoctorId = `
	SELECT consultation_sessions.id, user_id, doctor_id, consultation_session_status_id, 
	       consultation_sessions.created_at, consultation_sessions.updated_at,
	       consultation_session_statuses.name
	FROM consultation_sessions
	INNER JOIN consultation_session_statuses ON consultation_sessions.consultation_session_status_id = consultation_session_statuses.id 
	WHERE consultation_sessions.user_id = $1 AND consultation_sessions.doctor_id = $2`

	row := repo.db.QueryRowContext(ctx, findByUserIdAndDoctorId, userId, doctorId)
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

func (repo *ConsultationSessionRepositoryImpl) FindAllByUserIdOrDoctorId(ctx context.Context, userIdOrDoctorId int64, param *queryparamdto.GetAllParams) ([]*entity.ConsultationSession, error) {
	initQuery := `
	SELECT consultation_sessions.id, consultation_sessions.user_id, doctor_id, consultation_session_status_id,
    consultation_sessions.created_at, consultation_sessions.updated_at,
    consultation_session_statuses.name AS session_status,
    user_profiles.user_id, user_profiles.name, user_profiles.profile_photo,
    doctor_profiles.user_id, doctor_profiles.name, doctor_profiles.profile_photo,
    cm.id, cm.sender_id, cm.message, cm.attachment, cm.created_at AS message_created_at,
    cm.updated_at AS message_updated_at
	FROM  consultation_sessions
	INNER JOIN consultation_session_statuses ON consultation_sessions.consultation_session_status_id = consultation_session_statuses.id
	INNER JOIN user_profiles ON consultation_sessions.user_id = user_profiles.user_id
	INNER JOIN doctor_profiles ON consultation_sessions.doctor_id = doctor_profiles.user_id
	LEFT JOIN LATERAL (
		SELECT id, session_id, sender_id, message, attachment, created_at, updated_at
		FROM consultation_messages
		WHERE session_id = consultation_sessions.id
		ORDER BY created_at DESC
		LIMIT 1
	) cm ON true
	WHERE consultation_sessions.deleted_at IS NULL AND (consultation_sessions.user_id = $1 OR consultation_sessions.doctor_id = $1) `
	indexPreparedStatement := 1

	query, values := buildQuery(
		initQuery, &entity.ConsultationSession{}, param, false, false, indexPreparedStatement,
	)
	values = util.AppendAtIndex(values, 0, interface{}(userIdOrDoctorId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*entity.ConsultationSession, 0)
	for rows.Next() {
		var (
			session       entity.ConsultationSession
			sessionStatus entity.ConsultationSessionStatus
			userProfile   entity.UserProfile
			doctorProfile entity.DoctorProfile
			message       entity.ConsultationMessage
		)
		if err := rows.Scan(
			&session.Id, &session.UserId, &session.DoctorId, &session.ConsultationSessionStatusId,
			&session.CreatedAt, &session.UpdatedAt,
			&sessionStatus.Name,
			&userProfile.UserId, &userProfile.Name, &userProfile.ProfilePhoto,
			&doctorProfile.UserId, &doctorProfile.Name, &doctorProfile.ProfilePhoto,
			&message.Id, &message.SenderId, &message.Message, &message.Attachment, &message.CreatedAt, &message.UpdatedAt,
		); err != nil {
			return nil, err
		}
		session.ConsultationSessionStatus = &sessionStatus
		session.UserProfile = &userProfile
		session.DoctorProfile = &doctorProfile
		session.Message = []*entity.ConsultationMessage{&message}
		sessions = append(sessions, &session)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
