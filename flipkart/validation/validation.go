package validation

import (
	"regexp"
	"strings"
	"wallet/errors"
)

const (
	// Business rules
	MinTransactionAmount = 0.01
	MaxTransactionAmount = 100000.00
	MinLoadAmount        = 1.00
	MaxLoadAmount        = 50000.00
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateUserID validates user ID format
func ValidateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return &errors.ValidationError{Field: "userID", Message: "cannot be empty"}
	}
	if len(userID) < 3 {
		return &errors.ValidationError{Field: "userID", Message: "must be at least 3 characters long"}
	}
	if len(userID) > 50 {
		return &errors.ValidationError{Field: "userID", Message: "cannot exceed 50 characters"}
	}
	return nil
}

// ValidateEmailID validates email format
func ValidateEmailID(emailID string) error {
	if strings.TrimSpace(emailID) == "" {
		return &errors.ValidationError{Field: "emailID", Message: "cannot be empty"}
	}
	if !emailRegex.MatchString(emailID) {
		return &errors.ValidationError{Field: "emailID", Message: "invalid email format"}
	}
	return nil
}

// ValidateTransactionAmount validates transaction amount
func ValidateTransactionAmount(amount float64) error {
	if amount <= 0 {
		return &errors.InvalidAmountError{Amount: amount, Reason: "must be greater than 0"}
	}
	if amount < MinTransactionAmount {
		return &errors.InvalidAmountError{Amount: amount, Reason: "below minimum transaction amount"}
	}
	if amount > MaxTransactionAmount {
		return &errors.InvalidAmountError{Amount: amount, Reason: "exceeds maximum transaction amount"}
	}
	return nil
}

// ValidateLoadAmount validates load amount
func ValidateLoadAmount(amount float64) error {
	if amount <= 0 {
		return &errors.InvalidAmountError{Amount: amount, Reason: "must be greater than 0"}
	}
	if amount < MinLoadAmount {
		return &errors.InvalidAmountError{Amount: amount, Reason: "below minimum load amount"}
	}
	if amount > MaxLoadAmount {
		return &errors.InvalidAmountError{Amount: amount, Reason: "exceeds maximum load amount"}
	}
	return nil
}

// ValidatePaymentMethod validates payment method
func ValidatePaymentMethod(method string) error {
	validMethods := map[string]bool{
		"UPI":        true,
		"CREDITCARD": true,
		"DEBITCARD":  true,
	}

	if !validMethods[method] {
		return &errors.PaymentMethodError{Method: method}
	}
	return nil
}
