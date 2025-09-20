package sort

import (
	"fmt"
	"sort"
	"wallet/repository"
)

type Sort interface {
	Sort(userId string, repo *repository.Repository)
}

type AmountSort struct{}
type DateSort struct{}

func GetSortStrategy(sortType string) Sort {
	if sortType == "date" {
		return &DateSort{}
	}
	return &AmountSort{}
}

func (h *AmountSort) Sort(userId string, repo *repository.Repository) {
	repo.Mu.RLock()
	history := make([]*repository.Transaction, len(repo.History[userId]))
	copy(history, repo.History[userId])
	repo.Mu.RUnlock()

	sort.Slice(history, func(i int, j int) bool {
		return history[i].Amount < history[j].Amount
	})

	fmt.Printf("\n=== Transactions Sorted by Amount (Low to High) ===\n")
	for _, data := range history {
		fmt.Printf("ID: %s | Date: %v | Amount: %.2f | Action: %s | Payment: %s\n",
			data.ID, data.Date.Format("2006-01-02 15:04:05"), data.Amount, data.Action, data.Payment)
	}
}

func (h *DateSort) Sort(userId string, repo *repository.Repository) {
	repo.Mu.RLock()
	history := make([]*repository.Transaction, len(repo.History[userId]))
	copy(history, repo.History[userId])
	repo.Mu.RUnlock()

	sort.Slice(history, func(i int, j int) bool {
		return history[i].Date.Before(history[j].Date)
	})

	fmt.Printf("\n=== Transactions Sorted by Date (Oldest First) ===\n")
	for _, data := range history {
		fmt.Printf("ID: %s | Date: %v | Amount: %.2f | Action: %s | Payment: %s\n",
			data.ID, data.Date.Format("2006-01-02 15:04:05"), data.Amount, data.Action, data.Payment)
	}
}
