package lib

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidateStructErrMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate a struct field and then return a human friendly error message
func ValidateStruct(s any) (errorMessages []ValidateStructErrMessage) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(s)
	if err != nil {
		log.Print("Invalid input data")
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				var msg string

				switch e.Tag() {
				case "required":
					msg = fmt.Sprintf("%s is required", e.Field())
				case "email":
					msg = "Invalid email format"
				case "min":
					msg = fmt.Sprintf("%s is too short, at least minimum of %s characters", e.Field(), e.Param())
				case "max":
					msg = fmt.Sprintf("%s is too long, at least maximum of %s characters", e.Field(), e.Param())
				case "len":
					msg = fmt.Sprintf("%s must be at least %s long", e.Field(), e.Param())
				case "eq":
					msg = fmt.Sprintf("%s must be equal to %s", e.Field(), e.Param())
				case "ne":
					msg = fmt.Sprintf("%s must not be equal to %s", e.Field(), e.Param())
				case "lt":
					msg = fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
				case "lte":
					msg = fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
				case "gt":
					msg = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
				case "gte":
					msg = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
				case "alphanum":
					msg = fmt.Sprintf("%s must contain only alphanumeric characters", e.Field())
				case "url":
					msg = "Invalid URL format"
				case "uuid":
					msg = "Invalid UUID format"
				case "numeric":
					msg = "Only numeric values are allowed"
				case "boolean":
					msg = fmt.Sprintf("%s must be a boolean value", e.Field())
				default:
					msg = "Unsupported error message, please contact the developer :P"
				}

				errorMessages = append(errorMessages, ValidateStructErrMessage{Field: strings.ToLower(e.Field()), Message: msg})
			}

			return
		}
	}

	return nil
}
