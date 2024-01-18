package usecase

import (
	"context"
	"halodeksik-be/app/appconstant"
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
	session.ConsultationSessionStatusId = appconstant.ConsultationSessionStatusOngoing
	added, err := uc.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	return added, nil
}
