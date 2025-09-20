package interfaces

import "wallet/repository"

// WalletOperations defines the core wallet operations interface
type WalletOperations interface {
	RegisterUser(userID, emailID string) error
	LoadMoney(userID string, amount float64, paymentMethod string) error
	SendMoney(senderID, receiverID string, amount float64) error
	GetBalance(userID string) (float64, error)
	GetTransactionHistory(userID string) ([]*repository.Transaction, error)
}

// PaymentProcessor defines payment processing interface
type PaymentProcessor interface {
	LoadMoney(userID string, amount float64, repo *repository.Repository) error
	GetPaymentMethod() string
}

// TransactionSorter defines sorting interface
type TransactionSorter interface {
	Sort(transactions []*repository.Transaction) []*repository.Transaction
	GetSortType() string
}

// TransactionFilter defines filtering interface
type TransactionFilter interface {
	Filter(transactions []*repository.Transaction, criteria string) []*repository.Transaction
}

// UserRepository defines user storage interface
type UserRepository interface {
	Store(userID string, user interface{}) error
	Get(userID string) (interface{}, error)
	Exists(userID string) bool
	Delete(userID string) error
}

// WalletRepository defines wallet storage interface
type WalletRepository interface {
	UpdateBalance(userID string, amount float64) error
	GetBalance(userID string) (float64, error)
	AddTransaction(userID string, transaction *repository.Transaction) error
	GetTransactions(userID string) ([]*repository.Transaction, error)
}
