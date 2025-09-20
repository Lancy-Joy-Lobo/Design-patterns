package main

import (
	"fmt"
	"wallet/config"
	"wallet/repository"
	"wallet/service"
	"wallet/sort"
	"wallet/user"
	"wallet/validation"
)

func main() {
	fmt.Println("🏦 Flipkart Wallet System - Enterprise Edition")
	fmt.Println("==============================================")

	// Initialize configuration
	cfg := config.DefaultConfig()
	fmt.Printf("⚙️  System configured with max transaction amount: %.2f\n", cfg.MaxTransactionAmount)

	// Initialize repositories
	userRepo := user.CreateNewUserRepo()
	walletRepo := repository.CreateNewRepo()

	// Initialize service layer
	walletService := service.NewWalletService(userRepo, walletRepo)

	// Demonstrate user registration with validation
	fmt.Println("\n📝 User Registration:")
	user1, err := walletService.RegisterUser("lancy_joy_lobo", "lobolancy@gmail.com")
	if err != nil {
		fmt.Printf("❌ Registration failed: %s\n", err.Error())
		return
	}
	fmt.Printf("✅ User registered: %s (%s)\n", user1.UserID, user1.EmailID)

	user2, err := walletService.RegisterUser("john_doe", "john@example.com")
	if err != nil {
		fmt.Printf("❌ Registration failed: %s\n", err.Error())
		return
	}
	fmt.Printf("✅ User registered: %s (%s)\n", user2.UserID, user2.EmailID)

	// Demonstrate validation
	fmt.Println("\n🔍 Validation Examples:")
	if err := validation.ValidateEmailID("invalid-email"); err != nil {
		fmt.Printf("❌ Email validation: %s\n", err.Error())
	}
	if err := validation.ValidateTransactionAmount(-10); err != nil {
		fmt.Printf("❌ Amount validation: %s\n", err.Error())
	}

	// Demonstrate money loading with error handling
	fmt.Println("\n💰 Money Loading:")
	if err := walletService.LoadMoney(user1.UserID, 0, "UPI"); err != nil {
		fmt.Printf("❌ Load failed: %s\n", err.Error())
	}

	err = walletService.LoadMoney(user1.UserID, 1000, "UPI")
	if err != nil {
		fmt.Printf("❌ Load failed: %s\n", err.Error())
	} else {
		fmt.Printf("✅ Loaded ₹1000 to %s via UPI\n", user1.UserID)
	}

	err = walletService.LoadMoney(user1.UserID, 500, "CREDITCARD")
	if err != nil {
		fmt.Printf("❌ Load failed: %s\n", err.Error())
	} else {
		fmt.Printf("✅ Loaded ₹500 to %s via Credit Card\n", user1.UserID)
	}

	// Check balance
	balance, err := walletService.GetBalance(user1.UserID)
	if err != nil {
		fmt.Printf("❌ Balance check failed: %s\n", err.Error())
	} else {
		fmt.Printf("💳 Current balance for %s: ₹%.2f\n", user1.UserID, balance)
	}

	// Demonstrate money transfers
	fmt.Println("\n💸 Money Transfers:")
	transfers := []float64{100, 150, 75, 200}
	for _, amount := range transfers {
		err = walletService.SendMoney(user1.UserID, user2.UserID, amount)
		if err != nil {
			fmt.Printf("❌ Transfer failed: %s\n", err.Error())
		} else {
			fmt.Printf("✅ Transferred ₹%.2f from %s to %s\n", amount, user1.UserID, user2.UserID)
		}
	}

	// Check final balances
	fmt.Println("\n💰 Final Balances:")
	balance1, _ := walletService.GetBalance(user1.UserID)
	balance2, _ := walletService.GetBalance(user2.UserID)
	fmt.Printf("👤 %s: ₹%.2f\n", user1.UserID, balance1)
	fmt.Printf("👤 %s: ₹%.2f\n", user2.UserID, balance2)

	// Demonstrate insufficient balance scenario
	fmt.Println("\n⚠️  Testing Insufficient Balance:")
	err = walletService.SendMoney(user2.UserID, user1.UserID, 1000)
	if err != nil {
		fmt.Printf("❌ Transfer failed as expected: %s\n", err.Error())
	}

	// Load some money to user2
	walletService.LoadMoney(user2.UserID, 300, "DEBITCARD")
	fmt.Printf("✅ Loaded ₹300 to %s via Debit Card\n", user2.UserID)

	// Transfer back
	err = walletService.SendMoney(user2.UserID, user1.UserID, 100)
	if err != nil {
		fmt.Printf("❌ Transfer failed: %s\n", err.Error())
	} else {
		fmt.Printf("✅ Transferred ₹100 from %s to %s\n", user2.UserID, user1.UserID)
	}

	// Demonstrate improved sorting and filtering
	fmt.Println("\n📊 Transaction Analysis:")
	strategy := sort.GetSortStrategy("date")
	strategy.Sort(user1.UserID, walletRepo)

	strategy = sort.GetSortStrategy("amount")
	strategy.Sort(user1.UserID, walletRepo)

	// Filter transactions
	walletRepo.Filter(user1.UserID, "SEND")
	walletRepo.Filter(user2.UserID, "RECEIVE")

	// Final system status
	fmt.Println("\n📈 System Summary:")
	balance1, _ = walletService.GetBalance(user1.UserID)
	balance2, _ = walletService.GetBalance(user2.UserID)
	totalSystemBalance := balance1 + balance2
	fmt.Printf("🏦 Total system balance: ₹%.2f\n", totalSystemBalance)

	// Get transaction counts
	history1, _ := walletService.GetTransactionHistory(user1.UserID)
	history2, _ := walletService.GetTransactionHistory(user2.UserID)
	fmt.Printf("📝 Total transactions: %d\n", len(history1)+len(history2))

	fmt.Println("\n🎉 Flipkart Wallet System Demo Completed Successfully!")
}
