package config

import "time"

// WalletConfig holds all configuration for the wallet system
type WalletConfig struct {
	// Transaction limits
	MinTransactionAmount float64
	MaxTransactionAmount float64
	MinLoadAmount        float64
	MaxLoadAmount        float64

	// User limits
	MaxUsersPerSystem      int
	MaxTransactionsPerUser int

	// Security settings
	RequireEmailVerification bool
	EnableAuditLogging       bool
	SessionTimeout           time.Duration

	// Business rules
	AllowSelfTransfer     bool
	EnableTransactionFees bool
	TransactionFeePercent float64

	// System settings
	DatabaseTimeout  time.Duration
	MaxConcurrentOps int
}

// DefaultConfig returns default configuration
func DefaultConfig() *WalletConfig {
	return &WalletConfig{
		MinTransactionAmount:     0.01,
		MaxTransactionAmount:     100000.00,
		MinLoadAmount:            1.00,
		MaxLoadAmount:            50000.00,
		MaxUsersPerSystem:        1000000,
		MaxTransactionsPerUser:   10000,
		RequireEmailVerification: false,
		EnableAuditLogging:       true,
		SessionTimeout:           time.Hour * 24,
		AllowSelfTransfer:        false,
		EnableTransactionFees:    false,
		TransactionFeePercent:    0.0,
		DatabaseTimeout:          time.Second * 30,
		MaxConcurrentOps:         1000,
	}
}

// PaymentMethods returns supported payment methods
func SupportedPaymentMethods() []string {
	return []string{"UPI", "CREDITCARD", "DEBITCARD"}
}

// TransactionTypes returns supported transaction types
func SupportedTransactionTypes() []string {
	return []string{"LOAD", "SEND", "RECEIVE"}
}

// SortOptions returns supported sort options
func SupportedSortOptions() []string {
	return []string{"date", "amount"}
}

// FilterOptions returns supported filter options
func SupportedFilterOptions() []string {
	return []string{"SEND", "RECEIVE", "LOAD"}
}
