package apperror

import "errors"

var (
	ErrRecordNotFound             = errors.New("record not found")
	ErrForgotPasswordTokenInvalid = errors.New("forgot token is invalid")
	ErrForgotPasswordTokenExpired = errors.New("forgot token is already expired")
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

	ErrInvalidCityProvinceCombi = errors.New("invalid city and province combination")

	ErrPasswordTooLong       = errors.New("password too long")
	ErrStartDateAfterEndDate = errors.New("start date cannot be after end date")
	ErrForbiddenViewEntity   = errors.New("you are not allowed to view this entity")
	ErrForbiddenModifyEntity = errors.New("you are not allowed to modify this entity")

	ErrDeleteAlreadyAssignedAdmin = errors.New("cannot delete already assigned pharmacy admin")

	ErrInvalidDecimal      = errors.New("invalid decimal")
	ErrInvalidIntInString  = errors.New("invalid integer in string")
	ErrInvalidLatLong      = errors.New("invalid lat and long")
	ErrMainAddressNotFound = errors.New("this user does not have main address")

	ErrPharmacyProductUniqueConstraint = errors.New("pharmacy_id and product_id combinations violate unique constraint")

	ErrProductUniqueConstraint           = errors.New("name, generic_name, content, and manufacturer_id combinations violate unique constraint")
	ErrProductImageDoesNotExistInContext = errors.New("product image does not exist in context")

	ErrInsufficientProductStock             = errors.New("insufficient product stock")
	ErrProductStockNotEnoughToAddToCart     = errors.New("product stock is not enough to add to cart")
	ErrProductAddedToCartMustHaveAtLeastOne = errors.New("product added to cart must have at least one item")

	ErrRequestStockMutationFromOwnPharmacy  = errors.New("cannot request stock mutation from own pharmacy")
	ErrRequestStockMutationDifferentProduct = errors.New("requested product from destination pharmacy is not the same as origin pharmacy")
	ErrAlreadyFinishedRequest               = errors.New("request already finished")

	ErrGetShipmentMethodNoItems           = errors.New("there are no items to ship")
	ErrGetShipmentMethodDifferentPharmacy = errors.New("cannot get shipment method for different pharmacy")
	ErrGetShipmentCost                    = errors.New("failed to retrieve shipment cost")

	ErrPaymentSent      = errors.New("payment has already been sent")
	ErrPaymentConfirmed = errors.New("transaction has already been paid")
	ErrPaymentNotSent   = errors.New("transaction has not been paid")
	ErrBadConfirmStatus = errors.New("can only confirm waiting order")
	ErrBadRejectStatus  = errors.New("can only reject waiting order")
	ErrBadShipStatus    = errors.New("can only ship processed order")
	ErrBadReceiveStatus = errors.New("can only receive shipped order")
)
