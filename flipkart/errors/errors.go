package errors

import "fmt"

// Custom error types for better error handling

// ValidationError represents validation failures
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// InsufficientBalanceError represents insufficient balance scenarios
type InsufficientBalanceError struct {
	UserID    string
	Available float64
	Requested float64
}

func (e *InsufficientBalanceError) Error() string {
	return fmt.Sprintf("insufficient balance for user %s: available %.2f, requested %.2f",
		e.UserID, e.Available, e.Requested)
}

// UserNotFoundError represents cases where user doesn't exist
type UserNotFoundError struct {
	UserID string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.UserID)
}

// DuplicateUserError represents cases where user already exists
type DuplicateUserError struct {
	UserID string
}

func (e *DuplicateUserError) Error() string {
	return fmt.Sprintf("user %s already exists", e.UserID)
}

// InvalidAmountError represents invalid amount scenarios
type InvalidAmountError struct {
	Amount float64
	Reason string
}

func (e *InvalidAmountError) Error() string {
	return fmt.Sprintf("invalid amount %.2f: %s", e.Amount, e.Reason)
}

// PaymentMethodError represents payment method related errors
type PaymentMethodError struct {
	Method string
}

func (e *PaymentMethodError) Error() string {
	return fmt.Sprintf("invalid payment method: %s", e.Method)
}
