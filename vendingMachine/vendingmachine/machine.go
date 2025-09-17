package vendingmachine

import (
	"errors"
	"fmt"
	"sync"
)

type Product struct {
	Quantity  int
	ProductId string
	Name      string
	Price     int
}

type Denominations int

const (
	ONE Denominations = iota
	TWO
	FIVE
	TEN
	TWENTY
	HUNDRED
	FIVEHUNDRED
	THOUSAND
)

func (d Denominations) GetMoney() int {
	return [...]int{1, 2, 5, 10, 20, 100, 500, 1000}[d]
}

type VendingMachine struct {
	Mu             sync.Mutex
	Inventory      map[string]*Product
	TotalAmount    map[Denominations]int
	Denominations  map[Denominations]int
	State          VendingMachineInterface
	InsertedAmount int
}

func (vm *VendingMachine) SetState(state VendingMachineInterface) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()
	vm.State = state
}

// Error types for better error handling
type VendingMachineError struct {
	Code    string
	Message string
}

func (e *VendingMachineError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewVendingMachineError(code, message string) *VendingMachineError {
	return &VendingMachineError{Code: code, Message: message}
}

// Error codes
const (
	ErrInsufficientFunds  = "INSUFFICIENT_FUNDS"
	ErrProductNotFound    = "PRODUCT_NOT_FOUND"
	ErrInsufficientStock  = "INSUFFICIENT_STOCK"
	ErrInsufficientChange = "INSUFFICIENT_CHANGE"
	ErrInvalidQuantity    = "INVALID_QUANTITY"
	ErrInvalidPrice       = "INVALID_PRICE"
	ErrProductExists      = "PRODUCT_EXISTS"
	ErrInvalidState       = "INVALID_STATE"
)

// Validation helpers
func (vm *VendingMachine) ValidateProductId(productId string) error {
	if productId == "" {
		return NewVendingMachineError(ErrProductNotFound, "product ID cannot be empty")
	}
	return nil
}

func (vm *VendingMachine) ValidateQuantity(quantity int) error {
	if quantity <= 0 {
		return NewVendingMachineError(ErrInvalidQuantity, "quantity must be greater than 0")
	}
	return nil
}

func (vm *VendingMachine) ValidatePrice(price int) error {
	if price <= 0 {
		return NewVendingMachineError(ErrInvalidPrice, "price must be greater than 0")
	}
	return nil
}

// Transaction workflow methods
func (vm *VendingMachine) ProcessTransaction(productId string, quantity int) error {
	// Validate transaction
	if err := vm.State.SelectProduct(vm, productId, quantity); err != nil {
		return err
	}

	// Dispense product
	if err := vm.State.DispenseProduct(vm, productId, quantity); err != nil {
		return err
	}

	// Return change
	vm.State.ReturnChange(vm, productId, quantity)

	return nil
}

func (vm *VendingMachine) CanProvideChange(amount int) bool {
	if amount <= 0 {
		return true
	}

	denominations := []Denominations{THOUSAND, FIVEHUNDRED, HUNDRED, TWENTY, TEN, FIVE, TWO, ONE}
	tempAmount := amount

	for _, d := range denominations {
		if d.GetMoney() > tempAmount {
			continue
		}

		count := tempAmount / d.GetMoney()
		if available, exists := vm.TotalAmount[d]; exists && available >= count {
			tempAmount -= count * d.GetMoney()
		} else if exists && available > 0 {
			tempAmount -= available * d.GetMoney()
		}

		if tempAmount == 0 {
			return true
		}
	}

	return tempAmount == 0
}

func (vm *VendingMachine) ReturnChange() (map[Denominations]int, error) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	amount := vm.InsertedAmount

	denominations := []Denominations{THOUSAND, FIVEHUNDRED, HUNDRED, TWENTY, TEN, FIVE, TWO, ONE}
	change := make(map[Denominations]int)

	for _, d := range denominations {
		if d.GetMoney() > amount {
			continue
		}

		// Calculate how many of this denomination we need
		count := amount / d.GetMoney()

		// Check if we have enough of this denomination in the machine
		if available, exists := vm.TotalAmount[d]; exists && available >= count {
			change[d] = count
			vm.TotalAmount[d] -= count // Deduct from machine's available change
			amount -= count * d.GetMoney()
		} else if exists && available > 0 {
			// Use what we have available
			change[d] = available
			vm.TotalAmount[d] = 0
			amount -= available * d.GetMoney()
		}

		if amount == 0 {
			break
		}
	}

	if amount > 0 {
		// Revert the changes if we can't provide exact change
		for denom, count := range change {
			vm.TotalAmount[denom] += count
		}
		return nil, errors.New("insufficient change available in machine")
	}

	return change, nil
}

type VendingMachineInterface interface {
	InsertMoney(vm *VendingMachine, money map[Denominations]int)
	SelectProduct(vm *VendingMachine, productId string, quantity int) error
	DispenseProduct(vm *VendingMachine, productId string, quantity int) error
	ReturnChange(vm *VendingMachine, productId string, quantity int)
	CancelProduct(vm *VendingMachine, productIt string) error
}

// Administrative interface for maintenance operations
type VendingMachineAdmin interface {
	AddProduct(productId, name string, price, quantity int) error
	RestockProduct(productId string, quantity int) error
	UpdateProductPrice(productId string, newPrice int) error
	RemoveProduct(productId string) error
	CollectMoney() map[Denominations]int
	AddChangeToMachine(money map[Denominations]int)
	GetInventoryStatus() map[string]*Product
	GetMachineStatus() (int, map[Denominations]int) // returns total value and denomination breakdown
}

// Administrative functions implementation
func (vm *VendingMachine) AddProduct(productId, name string, price, quantity int) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	if _, exists := vm.Inventory[productId]; exists {
		return errors.New("product already exists, use RestockProduct to add quantity")
	}

	vm.Inventory[productId] = &Product{
		ProductId: productId,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
	}

	return nil
}

func (vm *VendingMachine) RestockProduct(productId string, quantity int) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	product, exists := vm.Inventory[productId]
	if !exists {
		return errors.New("product not found, use AddProduct to add new products")
	}

	vm.Inventory[productId].Quantity = product.Quantity + quantity
	return nil
}

func (vm *VendingMachine) UpdateProductPrice(productId string, newPrice int) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if newPrice <= 0 {
		return errors.New("price must be greater than 0")
	}

	_, exists := vm.Inventory[productId]
	if !exists {
		return errors.New("product not found")
	}

	vm.Inventory[productId].Price = newPrice
	return nil
}

func (vm *VendingMachine) RemoveProduct(productId string) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if _, exists := vm.Inventory[productId]; !exists {
		return errors.New("product not found")
	}

	delete(vm.Inventory, productId)
	return nil
}

func (vm *VendingMachine) CollectMoney() map[Denominations]int {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	collected := make(map[Denominations]int)
	for denom, count := range vm.TotalAmount {
		collected[denom] = count
		vm.TotalAmount[denom] = 0
	}

	return collected
}

func (vm *VendingMachine) AddChangeToMachine(money map[Denominations]int) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	for denom, count := range money {
		if count > 0 {
			vm.TotalAmount[denom] += count
		}
	}
}

func (vm *VendingMachine) GetInventoryStatus() map[string]*Product {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	// Return a copy to prevent external modification
	inventory := make(map[string]*Product)
	for id, product := range vm.Inventory {
		inventory[id] = &Product{
			ProductId: product.ProductId,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  product.Quantity,
		}
	}

	return inventory
}

func (vm *VendingMachine) GetMachineStatus() (int, map[Denominations]int) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	totalValue := 0
	denominations := make(map[Denominations]int)

	for denom, count := range vm.TotalAmount {
		totalValue += denom.GetMoney() * count
		denominations[denom] = count
	}

	return totalValue, denominations
}

// Constructor function - will be updated to use proper state import
func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		Inventory:      make(map[string]*Product),
		TotalAmount:    make(map[Denominations]int),
		Denominations:  make(map[Denominations]int),
		State:          nil, // Will be set after states are properly imported
		InsertedAmount: 0,
	}
}
