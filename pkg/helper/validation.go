package helper

import (
	"fashion-api/pkg/exception"

	"github.com/asaskevich/govalidator"
)

func ValidateStruct(s interface{}) exception.Exception {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return exception.NewBadRequestError(err.Error())
	}

	return nil
}
