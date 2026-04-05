package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(req interface{}) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	messages := make(map[string]string)

	for _, fe := range validationErrors {
		field := strings.ToLower(fe.Field())
		tag := fe.Tag()

		switch tag {
		case "required":
			messages[field] = "This field is required"
		case "email":
			messages[field] = "Must be a valid email address"
		case "min":
			param := fe.Param()
			if param == "" {
				messages[field] = "Value is too small"
			} else {
				messages[field] = "Must be at least " + param
			}
		case "max":
			param := fe.Param()
			if param == "" {
				messages[field] = "Value is too large"
			} else {
				messages[field] = "Must be no more than " + param
			}
		default:
			messages[field] = "Invalid value for " + field
		}
	}

	return messages
}
