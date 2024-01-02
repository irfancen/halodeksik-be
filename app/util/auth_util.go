package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"halodeksik-be/app/apperror"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUtil interface {
	ComparePassword(hashedPwd, plainPwd string) bool
	SignToken(token *jwt.Token) (string, error)
	HashAndSalt(pwd string) (string, error)
	GenerateSecureToken(length int) (string, error)
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

func (u *authUtil) GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "", apperror.ErrPasswordTooLong
	}
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
