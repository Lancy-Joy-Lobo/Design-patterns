package validation

import (
	"testing"
)

func TestValidateUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"Valid userID", "test_user_123", false},
		{"Empty userID", "", true},
		{"Short userID", "ab", true},
		{"Long userID", "this_is_a_very_long_user_id_that_exceeds_fifty_characters_limit", true},
		{"Valid minimum length", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailID(t *testing.T) {
	tests := []struct {
		name    string
		emailID string
		wantErr bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "user@sub.example.com", false},
		{"Empty email", "", true},
		{"Invalid format - no @", "testexample.com", true},
		{"Invalid format - no domain", "test@", true},
		{"Invalid format - no TLD", "test@example", true},
		{"Valid complex email", "test.user+tag@example.co.uk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailID(tt.emailID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTransactionAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"Valid amount", 50.0, false},
		{"Zero amount", 0, true},
		{"Negative amount", -10.0, true},
		{"Below minimum", 0.005, true},
		{"Above maximum", 150000.0, true},
		{"Minimum valid", 0.01, false},
		{"Maximum valid", 100000.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransactionAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransactionAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLoadAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"Valid load amount", 100.0, false},
		{"Zero amount", 0, true},
		{"Negative amount", -50.0, true},
		{"Below minimum", 0.5, true},
		{"Above maximum", 60000.0, true},
		{"Minimum valid", 1.0, false},
		{"Maximum valid", 50000.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoadAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoadAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePaymentMethod(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		wantErr bool
	}{
		{"Valid UPI", "UPI", false},
		{"Valid Credit Card", "CREDITCARD", false},
		{"Valid Debit Card", "DEBITCARD", false},
		{"Invalid method", "PAYPAL", true},
		{"Empty method", "", true},
		{"Lowercase method", "upi", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentMethod(tt.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePaymentMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
