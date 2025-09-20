package user

import (
	"fmt"
	"time"
	payment "wallet/payment"
	"wallet/repository"
)

type User struct {
	EmailID     string
	UserID      string
	TotalAmount float64
}

func CreateNewUser(userId, emailId string) *User {
	return &User{
		EmailID: emailId,
		UserID:  userId,
	}
}

func (user *User) LoadMoney(amount float64, payment payment.Payment, repo *repository.Repository, userRepo *UserRepo) error {
	_, exists := userRepo.UserRepo[user.UserID]
	if !exists {
		return fmt.Errorf("only registered users are allowed to load money")
	}
	return payment.LoadMoney(user.UserID, amount, repo)
}

func (user *User) SendMoney(receiver string, amount float64, repo *repository.Repository, userRepo *UserRepo) error {
	// Validate receiver exists first (can be done with user repo lock)
	userRepo.Mu.Lock()
	_, exists := userRepo.UserRepo[receiver]
	userRepo.Mu.Unlock()

	if !exists {
		return fmt.Errorf("receiver id does not exist")
	}

	// Use write lock for the entire transaction to avoid race conditions
	repo.Mu.Lock()
	defer repo.Mu.Unlock()

	// Check balance under the same lock that will be used for the transaction
	if amount > repo.Bank[user.UserID] {
		return fmt.Errorf("insufficient balance")
	}

	// Perform the transaction
	repo.Bank[user.UserID] -= amount
	repo.Bank[receiver] += amount

	repo.History[user.UserID] = append(repo.History[user.UserID], &repository.Transaction{
		ID:       repository.GenerateTransactionID(),
		Amount:   amount,
		Action:   "SEND",
		Date:     time.Now(),
		Payment:  "TRANSFER",
		FromUser: user.UserID,
		ToUser:   receiver,
	})

	repo.History[receiver] = append(repo.History[receiver], &repository.Transaction{
		ID:       repository.GenerateTransactionID(),
		Amount:   amount,
		Action:   "RECEIVE",
		Date:     time.Now(),
		Payment:  "TRANSFER",
		FromUser: user.UserID,
		ToUser:   receiver,
	})

	return nil
}

func (user *User) GetBalance(repo *repository.Repository) float64 {
	repo.Mu.RLock()
	defer repo.Mu.RUnlock()
	return repo.Bank[user.UserID]
}
