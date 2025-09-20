package service

import (
	"testing"
	"wallet/repository"
	"wallet/user"
)

func TestWalletService_RegisterUser(t *testing.T) {
	userRepo := user.CreateNewUserRepo()
	walletRepo := repository.CreateNewRepo()
	service := NewWalletService(userRepo, walletRepo)

	tests := []struct {
		name    string
		userID  string
		emailID string
		wantErr bool
	}{
		{
			name:    "Valid user registration",
			userID:  "test_user",
			emailID: "test@example.com",
			wantErr: false,
		},
		{
			name:    "Empty userID",
			userID:  "",
			emailID: "test@example.com",
			wantErr: true,
		},
		{
			name:    "Empty emailID",
			userID:  "test_user2",
			emailID: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.RegisterUser(tt.userID, tt.emailID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_LoadMoney(t *testing.T) {
	userRepo := user.CreateNewUserRepo()
	walletRepo := repository.CreateNewRepo()
	service := NewWalletService(userRepo, walletRepo)

	// Register a test user first
	_, err := service.RegisterUser("test_user", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	tests := []struct {
		name          string
		userID        string
		amount        float64
		paymentMethod string
		wantErr       bool
	}{
		{
			name:          "Valid load money",
			userID:        "test_user",
			amount:        100.0,
			paymentMethod: "UPI",
			wantErr:       false,
		},
		{
			name:          "Invalid amount (zero)",
			userID:        "test_user",
			amount:        0,
			paymentMethod: "UPI",
			wantErr:       true,
		},
		{
			name:          "Invalid amount (negative)",
			userID:        "test_user",
			amount:        -50,
			paymentMethod: "UPI",
			wantErr:       true,
		},
		{
			name:          "Invalid payment method",
			userID:        "test_user",
			amount:        100,
			paymentMethod: "INVALID",
			wantErr:       true,
		},
		{
			name:          "Non-existent user",
			userID:        "non_existent",
			amount:        100,
			paymentMethod: "UPI",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.LoadMoney(tt.userID, tt.amount, tt.paymentMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_SendMoney(t *testing.T) {
	userRepo := user.CreateNewUserRepo()
	walletRepo := repository.CreateNewRepo()
	service := NewWalletService(userRepo, walletRepo)

	// Register test users
	_, err := service.RegisterUser("sender", "sender@example.com")
	if err != nil {
		t.Fatalf("Failed to register sender: %v", err)
	}
	_, err = service.RegisterUser("receiver", "receiver@example.com")
	if err != nil {
		t.Fatalf("Failed to register receiver: %v", err)
	}

	// Load money for sender
	err = service.LoadMoney("sender", 100.0, "UPI")
	if err != nil {
		t.Fatalf("Failed to load money: %v", err)
	}

	tests := []struct {
		name       string
		senderID   string
		receiverID string
		amount     float64
		wantErr    bool
	}{
		{
			name:       "Valid money transfer",
			senderID:   "sender",
			receiverID: "receiver",
			amount:     50.0,
			wantErr:    false,
		},
		{
			name:       "Insufficient balance",
			senderID:   "sender",
			receiverID: "receiver",
			amount:     200.0,
			wantErr:    true,
		},
		{
			name:       "Invalid amount (zero)",
			senderID:   "sender",
			receiverID: "receiver",
			amount:     0,
			wantErr:    true,
		},
		{
			name:       "Self transfer",
			senderID:   "sender",
			receiverID: "sender",
			amount:     10.0,
			wantErr:    true,
		},
		{
			name:       "Non-existent sender",
			senderID:   "non_existent",
			receiverID: "receiver",
			amount:     10.0,
			wantErr:    true,
		},
		{
			name:       "Non-existent receiver",
			senderID:   "sender",
			receiverID: "non_existent",
			amount:     10.0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SendMoney(tt.senderID, tt.receiverID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	userRepo := user.CreateNewUserRepo()
	walletRepo := repository.CreateNewRepo()
	service := NewWalletService(userRepo, walletRepo)

	// Register a test user
	_, err := service.RegisterUser("test_user", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test initial balance
	balance, err := service.GetBalance("test_user")
	if err != nil {
		t.Errorf("GetBalance() error = %v", err)
	}
	if balance != 0 {
		t.Errorf("Expected initial balance 0, got %f", balance)
	}

	// Load money and test balance
	err = service.LoadMoney("test_user", 100.0, "UPI")
	if err != nil {
		t.Fatalf("Failed to load money: %v", err)
	}

	balance, err = service.GetBalance("test_user")
	if err != nil {
		t.Errorf("GetBalance() error = %v", err)
	}
	if balance != 100.0 {
		t.Errorf("Expected balance 100.0, got %f", balance)
	}

	// Test non-existent user
	_, err = service.GetBalance("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}
