package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ConsultationSessionUseCase interface {
	Add(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error)
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
