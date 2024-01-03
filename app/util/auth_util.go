package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type AuthUtil interface {
	ComparePassword(hashedPwd, plainPwd string) bool
	SignToken(token *jwt.Token) (string, error)
	HashAndSalt(pwd string) (string, error)
	GenerateSecureToken() (string, error)
}

func NewAuthUtil() AuthUtil {
	return &authUtil{}
}

type authUtil struct{}

func (u *authUtil) ComparePassword(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	password := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		return false
	}
	return true
}

func (u *authUtil) GenerateSecureToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (u *authUtil) SignToken(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *authUtil) HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
