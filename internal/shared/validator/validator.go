package validator

import (
	"modular-fx-fiber/internal/shared/logger"
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       any
	}

	Validator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResponse struct {
		Success bool   `json:"success" default:"false"`
		Message string `json:"message" default:"Internal server error"`
		Status  int    `json:"status" default:"500"`
	}
)

// Create a custom validator instance
var validate = validator.New()

// ValidateVNPhone validates a phone number
func validateVNPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Vietnamese phone number patterns:
	// - 10 digits starting with 0
	// - 11 digits starting with 84 (country code without +)
	// - 12 digits starting with +84 (country code with +)

	// Check for valid prefixes and length
	vnPhoneRegex := regexp.MustCompile(`^(0|\+84|84)([3|5|7|8|9])([0-9]{8})$`)
	return vnPhoneRegex.MatchString(phone)
}

// NewValidator creates a new validator
func NewValidator(l *logger.ZapLogger) *Validator {
	err := validate.RegisterValidation("vn_phone", validateVNPhone)

	if err != nil {
		l.Warn("Failed to register custom validation", zap.Error(err))
		return nil
	}

	return &Validator{
		validator: validate,
	}
}

func (cv *Validator) Validate(data any) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func (cv *Validator) ParseErrorToString(errs []ErrorResponse) string {
	var result string
	if len(errs) == 0 {
		return result
	}
	result = "Validation failed for the following fields:\n"
	for _, err := range errs {
		result += err.FailedField + " " + err.Tag + " " + err.Value.(string) + "\n"
	}

	return result
}
