package appvalidator

func AddCustomValidators(validator AppValidator) error {
	if err := validator.AddNewCustomValidation("filesize", FileSizeValidation); err != nil {
		return err
	}
	if err := validator.AddNewCustomValidation("filetype", FileTypeValidation); err != nil {
		return err
	}
	if err := validator.AddNewCustomValidation("numericgt", StringNumericGreaterThanValidation); err != nil {
		return err
	}
	if err := validator.AddNewCustomValidation("datetime", StringDateTimeValidation); err != nil {
		return err
	}
	return nil
}
