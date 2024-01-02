package apperror

import "errors"

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrForgotPasswordTokenInvalid = errors.New("token is invalid")
	ErrForgotPasswordTokenExpired = errors.New("token is already expired")
	ErrMissingAuthorizationToken  = errors.New("missing authorization token")
	ErrParsingBearerToken         = errors.New("failed to parse bearer token")
)
