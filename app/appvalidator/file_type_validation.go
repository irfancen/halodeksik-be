package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"halodeksik-be/app/util"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func FileTypeValidation(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	fileType := fl.Param()
	extension := filepath.Ext(fileHeader.Filename)
	extension = strings.Replace(extension, ".", "", 1)

	if util.IsEmptyString(extension) {
		return false
	}
	if !strings.Contains(fileType, extension) {
		return false
	}
	return true
}
