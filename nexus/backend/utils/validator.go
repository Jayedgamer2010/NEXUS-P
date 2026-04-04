package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateRequest validates a request struct and returns a map of field to error message.
// Returns nil if valid.
func ValidateRequest(req interface{}) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		var msg string
		switch fieldErr.Tag() {
		case "required":
			msg = "This field is required"
		case "min":
			msg = "Minimum value is " + fieldErr.Param()
		case "max":
			msg = "Maximum value is " + fieldErr.Param()
		case "email":
			msg = "Must be a valid email address"
		case "alphanum":
			msg = "Only alphanumeric characters allowed"
		case "hostname", "ip":
			msg = "Must be a valid hostname or IP address"
		case "oneof":
			msg = "Must be one of: " + fieldErr.Param()
		default:
			msg = "Invalid value"
		}
		errors[fieldErr.Field()] = msg
	}
	return errors
}
