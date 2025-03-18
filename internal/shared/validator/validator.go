package validator

import (
	"github.com/go-playground/validator/v10"
)

// Create a custom validator instance
var validate = validator.New()

// Validator is a struct that holds our validator instance
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		validator: validate,
	}
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	return "validation failed"
}

// ValidateStruct validates a struct and returns formatted errors
func (cv *Validator) ValidateStruct(s any) error {
	if err := cv.validator.Struct(s); err != nil {
		// If validation fails, return validation errors
		var errors ValidationErrors

		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}

		return errors
	}
	return nil
}
