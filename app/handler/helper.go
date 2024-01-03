package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/util"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func WrapError(err error, customCode ...int) error {
	errWrapper := &apperror.Wrapper{}

	if ok := errors.As(err, &errWrapper); !ok {
		errWrapper.ErrorStored = err
	}

	if len(customCode) > 0 {
		errWrapper.Code = customCode[0]
		return errWrapper
	}

	var (
		errJsonSyntax     *json.SyntaxError
		errJsonUnmarshall *json.UnmarshalTypeError
		errTimeParse      *time.ParseError
		errValidation     validator.ValidationErrors
		errNotFound       *apperror.NotFound
		errAlreadyExist   *apperror.AlreadyExist
		errNotMatch       *apperror.NotMatch
		errForbidden      *apperror.Forbidden
		errAuth           *apperror.AuthError
	)

	switch {
	case errors.Is(errWrapper.ErrorStored, context.DeadlineExceeded):
		errWrapper.Code = http.StatusGatewayTimeout
		errWrapper.Message = "server timeout, took too long to process request"

	case errors.Is(errWrapper.ErrorStored, io.EOF):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, io.ErrUnexpectedEOF):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errJsonSyntax):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errJsonUnmarshall):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errTimeParse):
		errWrapper.Code = http.StatusBadRequest

	case errWrapper.ErrorStored.Error() == "invalid request":
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrForgotPasswordTokenInvalid), errors.Is(errWrapper.ErrorStored, apperror.ErrForgotPasswordTokenExpired):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrRegisterTokenInvalid), errors.Is(errWrapper.ErrorStored, apperror.ErrRegisterTokenExpired):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errValidation):
		errWrapper.Code = http.StatusBadRequest
		errWrapper.Message = handleErrValidation(errValidation)

	case errors.As(errWrapper.ErrorStored, &errForbidden):
		errWrapper.Code = http.StatusForbidden

	case errors.As(errWrapper.ErrorStored, &errNotFound):
		errWrapper.Code = http.StatusNotFound

	case errors.As(errWrapper.ErrorStored, &errAlreadyExist):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errNotMatch):
		errWrapper.Code = http.StatusBadRequest

	case errors.As(errWrapper.ErrorStored, &errAuth):
		errWrapper.Code = http.StatusUnauthorized

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidDecimal):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrProductUniqueConstraint):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrInvalidRegisterRole):
		errWrapper.Code = http.StatusBadRequest

	case errors.Is(errWrapper.ErrorStored, apperror.ErrWrongCredentials):
		errWrapper.Code = http.StatusBadRequest

	default:
		errWrapper.Code = http.StatusInternalServerError
	}

	return errWrapper
}

func handleErrValidation(ve validator.ValidationErrors) string {
	buff := bytes.NewBufferString("")

	for i, _ := range ve {
		buff.WriteString(createErrValidationMsgTag(ve[i]))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func createErrValidationMsgTag(fieldError validator.FieldError) string {
	fieldName := util.PascalToSnake(fieldError.Field())
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", fieldName)
	case "email":
		return fmt.Sprintf("field '%s' must be in the format of an email", fieldName)
	case "number":
		return fmt.Sprintf("field '%s' must be a number", fieldName)
	case "numeric":
		return fmt.Sprintf("field '%s' must be numeric", fieldName)
	case "len":
		return fmt.Sprintf("field '%s' must have exactly %s characters long", fieldName, fieldError.Param())
	case "min":
		switch fieldError.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("field '%s' must be at least %s characters long", fieldName, fieldError.Param())
		case reflect.Slice:
			return fmt.Sprintf("field '%s' must have at least %s item", fieldName, fieldError.Param())
		}
		return fmt.Sprintf("field '%s' have a minimum value of %s", fieldName, fieldError.Param())
	case "max":
		switch fieldError.Type().Kind() {
		case reflect.String:
			return fmt.Sprintf("field '%s' must be at maximum %s characters long", fieldName, fieldError.Param())
		case reflect.Slice:
			return fmt.Sprintf("field '%s' must have at maximum %s item", fieldName, fieldError.Param())
		}
		return fmt.Sprintf("field '%s' have maximum value of %s", fieldName, fieldError.Param())
	case "oneof":
		params := strings.ReplaceAll(fieldError.Param(), " ", ", ")
		return fmt.Sprintf("item '%s' on field '%s' must be one of %s", fieldError.Value(), fieldName, params)
	case "latitude":
		return fmt.Sprintf("field '%s' must be a valid latitude", fieldName)
	case "longitude":
		return fmt.Sprintf("field '%s' must be a valid longitude", fieldName)
	case "gtfield":
		otherFieldName := util.PascalToSnake(fieldError.Param())
		return fmt.Sprintf("field '%s' must be greater than '%s' field", fieldName, otherFieldName)
	default:
		msg := fmt.Sprintf("field '%s' failed on validation %s %s", fieldName, fieldError.Tag(), fieldError.Param())
		return strings.TrimSpace(msg)
	}
}
