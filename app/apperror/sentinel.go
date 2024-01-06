package apperror

import "errors"

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrForgotPasswordTokenInvalid = errors.New("token is invalid")
	ErrForgotPasswordTokenExpired = errors.New("token is already expired")
	ErrMissingAuthorizationToken  = errors.New("missing authorization token")
	ErrParsingBearerToken         = errors.New("failed to parse bearer token")
	ErrRegisterTokenInvalid       = errors.New("register token is invalid")
	ErrRegisterTokenExpired       = errors.New("register token is already expired")
	ErrInvalidRegisterRole        = errors.New("invalid register role, only doctor and user are allowed")
	ErrWrongCredentials           = errors.New("wrong password or email")

	ErrLoginNoToken          = errors.New("login token is not provided")
	ErrLoginTokenInvalidSign = errors.New("invalid signature")
	ErrLoginTokenNotValid    = errors.New("login token is invalid")
	ErrUnauthorized          = errors.New("you don't have permission to access this endpoint")

	ErrPasswordTooLong = errors.New("password too long")

	ErrInvalidDecimal     = errors.New("invalid decimal")
	ErrInvalidIntInString = errors.New("invalid integer in string")

	ErrPharmacyProductUniqueConstraint = errors.New("pharmacy_id and product_id combinations violate unique constraint")

	ErrProductUniqueConstraint           = errors.New("name, generic_name, content, and manufacturer_id combinations violate unique constraint")
	ErrProductImageDoesNotExistInContext = errors.New("product image does not exist in context")

	ErrProductStockNotEnoughToAddToCart = errors.New("product stock is not enough to add to cart")
)
