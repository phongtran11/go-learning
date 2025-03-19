package validator

import (
	"modular-fx-fiber/internal/shared/logger"
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Create a custom validator instance
var validate = validator.New()

// Validator is a struct that holds our validator instance
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator
func NewValidator(l *logger.ZapLogger) *Validator {
	err := validate.RegisterValidation("vn_phone", ValidateVNPhone)
	if err != nil {
		l.Warn("Failed to register custom validation", zap.Error(err))
		return nil
	}

	return &Validator{
		validator: validate,
	}
}

// ValidateVNPhone validates a phone number
func ValidateVNPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Vietnamese phone number patterns:
	// - 10 digits starting with 0
	// - 11 digits starting with 84 (country code without +)
	// - 12 digits starting with +84 (country code with +)

	// Check for valid prefixes and length
	vnPhoneRegex := regexp.MustCompile(`^(0|\+84|84)([3|5|7|8|9])([0-9]{8})$`)
	return vnPhoneRegex.MatchString(phone)
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
