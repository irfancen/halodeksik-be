package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/util"
	"strings"
)

func CommaSeparatedValidation(fl validator.FieldLevel) bool {
	valuesInStr := fl.Field().Interface().(string)
	valuesInStr = strings.TrimSpace(valuesInStr)

	if util.IsEmptyString(valuesInStr) {
		return false
	}

	valuesWithoutComma := strings.ReplaceAll(valuesInStr, ",", "")
	if len(valuesWithoutComma) == 0 {
		return false
	}

	return true
}
