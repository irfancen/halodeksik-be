package apperror

import "errors"

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrForgotPasswordTokenInvalid = errors.New("token is invalid")
	ErrForgotPasswordTokenExpired = errors.New("token is already expired")
	ErrMissingAuthorizationToken  = errors.New("missing authorization token")
	ErrParsingBearerToken         = errors.New("failed to parse bearer token")

	ErrPasswordTooLong = errors.New("password too long")

	ErrInvalidDecimal = errors.New("invalid decimal")

	ErrProductUniqueConstraint = errors.New("name, generic_name, content, and manufacturer_id combinations violate unique constraint")
)
