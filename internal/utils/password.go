package utils

import (
	"regexp"

	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
)

func ValidatePasswordStrength(password string) []*common.ValidationError {
	var errors []*common.ValidationError

	if len(password) < 8 {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
	}

	if len(password) > 32 {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password must not exceed 32 characters",
		})
	}

	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter (a-z)",
		})
	}

	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter (A-Z)",
		})
	}

	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password) {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one number (0-9)",
		})
	}

	validCharsRegex := regexp.MustCompile(`^[A-Za-z\d]+$`)
	if !validCharsRegex.MatchString(password) {
		errors = append(errors, &common.ValidationError{
			Field:   "password",
			Message: "Password contains invalid characters. Only letters and numbers are allowed",
		})
	}

	return errors
}

// ValidatePasswordMatch checks if password and confirmation match
func ValidatePasswordMatch(password, confirmation string) *common.ValidationError {
	if password != confirmation {
		return &common.ValidationError{
			Field:   "password_confirmation",
			Message: "Password and confirmation do not match",
		}
	}
	return nil
}
