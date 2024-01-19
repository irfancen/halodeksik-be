package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ConsultationSessionUseCase interface {
	Add(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error)
	GetById(ctx context.Context, id int64) (*entity.ConsultationSession, error)
	GetAllByUserIdOrDoctorId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type ConsultationSessionUseCaseImpl struct {
	repo repository.ConsultationSessionRepository
}

func NewConsultationSessionUseCaseImpl(repo repository.ConsultationSessionRepository) *ConsultationSessionUseCaseImpl {
	return &ConsultationSessionUseCaseImpl{repo: repo}
}

func (uc *ConsultationSessionUseCaseImpl) Add(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.repo.FindByUserIdAndDoctorId(ctx, session.UserId, session.DoctorId)
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	if !errors.Is(err, apperror.ErrRecordNotFound) && sessionDb.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusOngoing {
		return sessionDb, apperror.ErrChatStillOngoing
	}

	session.ConsultationSessionStatusId = appconstant.ConsultationSessionStatusOngoing
	added, err := uc.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	return added, nil
}

func (uc *ConsultationSessionUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(sessionDb, "Id", id)
		}
		return nil, err
	}

	return sessionDb, nil
}

func (uc *ConsultationSessionUseCaseImpl) GetAllByUserIdOrDoctorId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userIdOrDoctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	sessions, err := uc.repo.FindAllByUserIdOrDoctorId(ctx, userIdOrDoctorId, param)
	if err != nil {
		return nil, err
	}

	paginatedItems := entity.NewPaginationInfo(int64(len(sessions)), 1, int64(len(sessions)), 1, sessions)
	return paginatedItems, nil
}
