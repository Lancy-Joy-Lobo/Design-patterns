package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	ID       string
	Amount   float64
	Action   string
	Date     time.Time
	Payment  string
	FromUser string
	ToUser   string
}

type Repository struct {
	Mu      sync.RWMutex
	History map[string][]*Transaction
	Bank    map[string]float64
}

func CreateNewRepo() *Repository {
	return &Repository{
		History: make(map[string][]*Transaction),
		Bank:    make(map[string]float64),
	}
}

// GenerateTransactionID generates a unique transaction ID
func GenerateTransactionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (repo *Repository) Filter(userId, action string) {
	repo.Mu.RLock()
	history := repo.History[userId]
	repo.Mu.RUnlock()

	fmt.Printf("\n=== %s Transactions for User %s ===\n", action, userId)
	for _, data := range history {
		if data.Action == action {
			fmt.Printf("ID: %s | Date: %v | Amount: %.2f | Action: %s | Payment: %s\n",
				data.ID, data.Date.Format("2006-01-02 15:04:05"), data.Amount, data.Action, data.Payment)
		}
	}
}
