package payment

import "fmt"

// PaymentMethod defines the interface for payment methods.
type PaymentMethod interface {
	Credit(amount float64) error
	Debit(amount float64) error
	GetType() string
}

type CreditCard struct{}
type Cash struct{}

func (c CreditCard) Credit(amount float64) error {
	fmt.Println("Payment received")
	return nil
}

func (c Cash) Credit(amount float64) error {
	fmt.Println("Payment received")
	return nil
}

func (c CreditCard) Debit(amount float64) error {
	fmt.Println("Payment received")
	return nil
}

func (c Cash) Debit(amount float64) error {
	fmt.Println("Payment received")
	return nil
}

func (c CreditCard) GetType() string {
	return "CREDIT_CARD"
}

func (c Cash) GetType() string {
	return "CASH"
}

// GetPaymentMethod returns a payment method by type.
func GetPaymentMethod(method string) PaymentMethod {
	switch method {
	case "CASH":
		return Cash{}
	case "CREDIT_CARD":
		return CreditCard{}
	default:
		return nil
	}
}
