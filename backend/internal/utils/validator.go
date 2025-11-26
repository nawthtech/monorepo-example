package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() error {
	validate = validator.New()
	
	// تسجيل تحقق مخصص للأسماء العربية
	validate.RegisterValidation("arabicname", func(fl validator.FieldLevel) bool {
		name := fl.Field().String()
		// نمط للأسماء العربية (تسمح بالحروف العربية والفراغات)
		pattern := `^[\p{Arabic}\s]+$`
		matched, _ := regexp.MatchString(pattern, name)
		return matched
	})

	// تسجيل تحقق مخصص لأرقام الهواتف السعودية
	validate.RegisterValidation("saudimobile", func(fl validator.FieldLevel) bool {
		mobile := fl.Field().String()
		// نمط لأرقام الهواتف السعودية
		pattern := `^(009665|9665|\+9665|05)(5|0|3|6|4|9|1|8|7)([0-9]{7})$`
		matched, _ := regexp.MatchString(pattern, mobile)
		return matched
	})

	return nil
}

func GetValidator() *validator.Validate {
	return validate
}
