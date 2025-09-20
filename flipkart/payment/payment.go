package payment

import (
	"fmt"
	"time"
	"wallet/repository"
)

type Payment interface {
	LoadMoney(id string, amount float64, repo *repository.Repository) error
}

func GetPayment(paymentType string) (Payment, error) {
	switch paymentType {
	case "UPI":
		return &Upi{}, nil
	case "CREDITCARD":
		return &CreditCard{}, nil
	case "DEBITCARD":
		return &DebitCard{}, nil
	default:
		return nil, fmt.Errorf("invalid payment mode")
	}
}

type CreditCard struct{}

func (c *CreditCard) LoadMoney(id string, amount float64, repo *repository.Repository) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	repo.Mu.Lock()
	defer repo.Mu.Unlock()

	repo.Bank[id] += amount

	// Add transaction to history
	transaction := &repository.Transaction{
		ID:       repository.GenerateTransactionID(),
		Amount:   amount,
		Action:   "LOAD",
		Date:     time.Now(),
		Payment:  "CREDITCARD",
		FromUser: "",
		ToUser:   id,
	}
	repo.History[id] = append(repo.History[id], transaction)

	return nil
}

type DebitCard struct{}

func (c *DebitCard) LoadMoney(id string, amount float64, repo *repository.Repository) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	repo.Mu.Lock()
	defer repo.Mu.Unlock()

	repo.Bank[id] += amount

	// Add transaction to history
	transaction := &repository.Transaction{
		ID:       repository.GenerateTransactionID(),
		Amount:   amount,
		Action:   "LOAD",
		Date:     time.Now(),
		Payment:  "DEBITCARD",
		FromUser: "",
		ToUser:   id,
	}
	repo.History[id] = append(repo.History[id], transaction)

	return nil
}

type Upi struct{}

func (c *Upi) LoadMoney(id string, amount float64, repo *repository.Repository) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	repo.Mu.Lock()
	defer repo.Mu.Unlock()

	repo.Bank[id] += amount

	// Add transaction to history
	transaction := &repository.Transaction{
		ID:       repository.GenerateTransactionID(),
		Amount:   amount,
		Action:   "LOAD",
		Date:     time.Now(),
		Payment:  "UPI",
		FromUser: "",
		ToUser:   id,
	}
	repo.History[id] = append(repo.History[id], transaction)

	return nil
}
