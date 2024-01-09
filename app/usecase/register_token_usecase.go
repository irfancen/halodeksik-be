package usecase

import (
	"context"
	"errors"
	"fmt"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
	"os"
	"strconv"
	"time"
)

type RegisterTokenUseCase interface {
	SendRegisterToken(ctx context.Context, email string) (string, error)
	VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error)
}

type RegisterTokenUseCaseImpl struct {
	registerTokenRepository repository.RegisterTokenRepository
	userRepository          repository.UserRepository
	authUtil                util.AuthUtil
	mailUtil                util.EmailUtil
	registerTokenExpired    int
	loginTokenExpired       int
}

func NewRegisterTokenUseCase(uRepo repository.UserRepository, tRegisterRepo repository.RegisterTokenRepository, aUtil util.AuthUtil, eUtil util.EmailUtil) RegisterTokenUseCase {

	expiryRegister, err := strconv.Atoi(os.Getenv("REGISTER_TOKEN_EXPIRED_MINUTE"))
	if err != nil {
		return nil
	}

	return &RegisterTokenUseCaseImpl{
		userRepository:          uRepo,
		registerTokenRepository: tRegisterRepo,
		authUtil:                aUtil,
		mailUtil:                eUtil,
		registerTokenExpired:    expiryRegister,
	}
}

func (uc *RegisterTokenUseCaseImpl) VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	existedToken, err := uc.registerTokenRepository.FindRegisterTokenByToken(ctx, token)
	if existedToken == nil {
		return nil, apperror.NewNotFound(existedToken, "Token", token)
	}
	if err != nil {
		return nil, err
	}

	if existedToken.IsValid == false {
		return nil, apperror.ErrRegisterTokenInvalid
	}

	if existedToken.ExpiredAt.Before(time.Now()) {
		return nil, apperror.ErrRegisterTokenExpired
	}
	return existedToken, nil
}

func (uc *RegisterTokenUseCaseImpl) SendRegisterToken(ctx context.Context, email string) (string, error) {
	var userVerify entity.VerificationToken

	existedUser, err := uc.userRepository.FindByEmail(ctx, email)
	if existedUser != nil {
		return "", apperror.NewAlreadyExist(existedUser, "Email", email)
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	uid, err := uc.authUtil.GenerateSecureToken()
	if err != nil {
		return "", err
	}

	tokenFound, err := uc.registerTokenRepository.FindRegisterTokenByToken(ctx, uid)
	if tokenFound != nil {
		return "", &apperror.AlreadyExist{
			Resource:        tokenFound,
			FieldInResource: "Token",
			Value:           tokenFound.Token,
		}
	}

	activeToken, err := uc.registerTokenRepository.FindRegisterTokenByEmail(ctx, email)
	if activeToken != nil {
		_, err2 := uc.registerTokenRepository.DeactivateRegisterToken(ctx, *activeToken)
		if err2 != nil {
			return "", err2
		}
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	userVerify.Token = uid
	userVerify.Email = email
	userVerify.ExpiredAt = time.Now().Add(time.Duration(uc.registerTokenExpired) * time.Minute)
	userVerify.IsValid = true

	_, err = uc.registerTokenRepository.CreateRegisterToken(ctx, userVerify)
	if err != nil {
		return "", err
	}

	to := []string{email}
	subject := "Email Verification"
	message := fmt.Sprintf("Verification link:\n%s/verify-register?token=%s", os.Getenv("FRONTEND_URL"), uid)

	err = uc.mailUtil.SendEmail(to, []string{}, subject, message)
	if err != nil {
		return "", err
	}

	return uid, nil
}
