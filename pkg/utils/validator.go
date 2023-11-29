package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := make(map[string]string) // or map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		//fmt.Println("field", err.Field(), err.Tag(), err.Param(),err.Error())
		fields[err.Field()] = err.Error()
	}
	fmt.Println(fields)
	return fields
}
