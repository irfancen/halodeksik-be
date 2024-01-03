package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
	"os"
	"strconv"
	"time"
)

type AuthUsecase interface {
	SendRegisterToken(ctx context.Context, email string) (string, error)
	VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error)
	Register(ctx context.Context, user entity.User, token string) (*entity.User, error)
	Login(ctx context.Context, req requestdto.LoginRequest) (*entity.User, string, error)
}

type AuthUseCaseImpl struct {
	userRepository        repository.UserRepository
	verifyTokenRepository repository.VerifyTokenRepository
	authUtil              util.AuthUtil
	mailUtil              util.EmailUtil
	registerTokenExpired  int
	loginTokenExpired     int
}

func (uc *AuthUseCaseImpl) VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	existedToken, err := uc.verifyTokenRepository.FindTokenByToken(ctx, token)
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

func (uc *AuthUseCaseImpl) SendRegisterToken(ctx context.Context, email string) (string, error) {
	var userVerify entity.VerificationToken

	existedUser, err := uc.userRepository.FindByEmail(ctx, email)
	if existedUser != nil {
		return "", apperror.NewAlreadyExist(existedUser, "Email", email)
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	token, err := uc.authUtil.GenerateSecureToken()
	if err != nil {
		return "", err
	}

	tokenFound, err := uc.verifyTokenRepository.FindTokenByToken(ctx, token)
	if tokenFound != nil {
		return "", &apperror.AlreadyExist{
			Resource:        tokenFound,
			FieldInResource: "Token",
			Value:           tokenFound.Token,
		}
	}

	activeToken, err := uc.verifyTokenRepository.FindTokenByEmail(ctx, email)
	if activeToken != nil {
		_, err2 := uc.verifyTokenRepository.DeactivateToken(ctx, *activeToken)
		if err2 != nil {
			return "", err
		}
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	userVerify.Token = token
	userVerify.Email = email
	userVerify.ExpiredAt = time.Now().Add(time.Duration(uc.registerTokenExpired) * time.Minute)
	userVerify.IsValid = true

	_, err = uc.verifyTokenRepository.CreateVerifyToken(ctx, userVerify)
	if err != nil {
		return "", err
	}

	to := []string{email}
	subject := "Email Verification"
	message := fmt.Sprintf("Verification link:\n%s/verify-register?token=%s", os.Getenv("FRONTEND_URL"), token)

	err = uc.mailUtil.SendEmail(to, []string{}, subject, message)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *AuthUseCaseImpl) Register(ctx context.Context, user entity.User, token string) (*entity.User, error) {
	verifiedToken, err := uc.VerifyRegisterToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if verifiedToken.Email != user.Email {
		return nil, apperror.ErrRegisterTokenInvalid
	}

	existedUser, err := uc.userRepository.FindByEmail(ctx, user.Email)
	if err == nil {
		return nil, apperror.NewAlreadyExist(existedUser, "Email", user.Email)
	}
	if !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	if user.UserRoleId == appconstant.UserRoleIdAdmin || user.UserRoleId == appconstant.UserRoleIdPharmacyAdmin {
		return nil, apperror.ErrInvalidRegisterRole
	}

	// todo: upload doctor certificate
	//if user.UserRoleId == appconstant.UserRoleIdDoctor {
	//}

	// todo: handle hash errors
	hashedPw, err := uc.authUtil.HashAndSalt(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPw
	user.IsVerified = true
	createdUser, err := uc.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	_, err = uc.verifyTokenRepository.DeactivateToken(ctx, *verifiedToken)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, req requestdto.LoginRequest) (*entity.User, string, error) {
	user, err := uc.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}

	if !uc.authUtil.ComparePassword(user.Password, req.Password) {
		return nil, "", apperror.ErrWrongCredentials
	}

	expirationTime := time.Now().Add(time.Duration(uc.loginTokenExpired) * time.Minute)
	claims := &entity.Claims{
		UserId: user.Id,
		Email:  user.Email,
		RoleId: user.UserRoleId,
		Image:  "",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ByeByeSick Healthcare",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := uc.authUtil.SignToken(token)

	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

func NewAuthUsecase(a repository.UserRepository, b repository.VerifyTokenRepository, u util.AuthUtil, m util.EmailUtil) AuthUsecase {

	expiryRegister, err := strconv.Atoi(os.Getenv("REGISTER_TOKEN_EXPIRED_MINUTE"))
	if err != nil {
		return nil
	}

	expiryLogin, err := strconv.Atoi(os.Getenv("LOGIN_TOKEN_EXPIRED_MINUTE"))
	if err != nil {
		return nil
	}

	return &AuthUseCaseImpl{
		userRepository:        a,
		verifyTokenRepository: b,
		authUtil:              u,
		mailUtil:              m,
		registerTokenExpired:  expiryRegister,
		loginTokenExpired:     expiryLogin,
	}
}
