package service

import (
	"fmt"
	"time"
	"wallet/payment"
	"wallet/repository"
	"wallet/user"
)

// WalletService handles all wallet operations
type WalletService struct {
	userRepo   *user.UserRepo
	walletRepo *repository.Repository
}

// NewWalletService creates a new wallet service
func NewWalletService(userRepo *user.UserRepo, walletRepo *repository.Repository) *WalletService {
	return &WalletService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

// RegisterUser registers a new user
func (w *WalletService) RegisterUser(userID, emailID string) (*user.User, error) {
	if userID == "" || emailID == "" {
		return nil, fmt.Errorf("userID and emailID cannot be empty")
	}

	// Check if user already exists
	w.userRepo.Mu.Lock()
	if _, exists := w.userRepo.UserRepo[userID]; exists {
		w.userRepo.Mu.Unlock()
		return nil, fmt.Errorf("user %s already exists", userID)
	}
	w.userRepo.Mu.Unlock()

	newUser := user.CreateNewUser(userID, emailID)
	w.userRepo.Register(newUser)

	// Initialize wallet balance
	w.walletRepo.Mu.Lock()
	w.walletRepo.Bank[userID] = 0
	w.walletRepo.Mu.Unlock()

	return newUser, nil
}

// LoadMoney loads money into user's wallet
func (w *WalletService) LoadMoney(userID string, amount float64, paymentMethod string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Check if user exists
	w.userRepo.Mu.Lock()
	_, exists := w.userRepo.UserRepo[userID]
	w.userRepo.Mu.Unlock()

	if !exists {
		return fmt.Errorf("user %s not found", userID)
	}

	paymentProcessor, err := payment.GetPayment(paymentMethod)
	if err != nil {
		return err
	}

	return paymentProcessor.LoadMoney(userID, amount, w.walletRepo)
}

// SendMoney transfers money between users
func (w *WalletService) SendMoney(senderID, receiverID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if senderID == receiverID {
		return fmt.Errorf("cannot send money to yourself")
	}

	// Check if both users exist
	w.userRepo.Mu.Lock()
	_, senderExists := w.userRepo.UserRepo[senderID]
	_, receiverExists := w.userRepo.UserRepo[receiverID]
	w.userRepo.Mu.Unlock()

	if !senderExists {
		return fmt.Errorf("sender %s not found", senderID)
	}
	if !receiverExists {
		return fmt.Errorf("receiver %s not found", receiverID)
	}

	// Perform transaction with proper locking
	w.walletRepo.Mu.Lock()
	defer w.walletRepo.Mu.Unlock()

	senderBalance := w.walletRepo.Bank[senderID]
	if amount > senderBalance {
		return fmt.Errorf("insufficient balance: available %.2f, requested %.2f", senderBalance, amount)
	}

	// Execute the transaction
	w.walletRepo.Bank[senderID] -= amount
	w.walletRepo.Bank[receiverID] += amount

	// Record transaction history
	now := time.Now()
	w.walletRepo.History[senderID] = append(w.walletRepo.History[senderID], &repository.Transaction{
		Amount:  amount,
		Action:  "SEND",
		Date:    now,
		Payment: "TRANSFER",
	})

	w.walletRepo.History[receiverID] = append(w.walletRepo.History[receiverID], &repository.Transaction{
		Amount:  amount,
		Action:  "RECEIVE",
		Date:    now,
		Payment: "TRANSFER",
	})

	return nil
}

// GetBalance returns user's current balance
func (w *WalletService) GetBalance(userID string) (float64, error) {
	w.userRepo.Mu.Lock()
	_, exists := w.userRepo.UserRepo[userID]
	w.userRepo.Mu.Unlock()

	if !exists {
		return 0, fmt.Errorf("user %s not found", userID)
	}

	w.walletRepo.Mu.RLock()
	defer w.walletRepo.Mu.RUnlock()
	return w.walletRepo.Bank[userID], nil
}

// GetTransactionHistory returns user's transaction history
func (w *WalletService) GetTransactionHistory(userID string) ([]*repository.Transaction, error) {
	w.userRepo.Mu.Lock()
	_, exists := w.userRepo.UserRepo[userID]
	w.userRepo.Mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("user %s not found", userID)
	}

	w.walletRepo.Mu.RLock()
	defer w.walletRepo.Mu.RUnlock()

	// Return a copy to avoid concurrent modification
	history := make([]*repository.Transaction, len(w.walletRepo.History[userID]))
	copy(history, w.walletRepo.History[userID])

	return history, nil
}
