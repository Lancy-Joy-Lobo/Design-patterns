package states

import (
	"errors"
	"fmt"
	"vendingMachine/vendingmachine"
)

type Dispense struct{}

var _ vendingmachine.VendingMachineInterface = (*Dispense)(nil)

func (i *Dispense) InsertMoney(vm *vendingmachine.VendingMachine, money map[vendingmachine.Denominations]int) {
	fmt.Println("entering money is not allowed")
}

func (i *Dispense) SelectProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return errors.New("selecting product is not allowed")
}

func (i *Dispense) DispenseProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	product, exists := vm.Inventory[productId]
	if !exists {
		return errors.New("product not found")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if product.Quantity < quantity {
		return errors.New("insufficient product quantity")
	}

	// Deduct cost from inserted amount
	totalCost := quantity * product.Price
	if vm.InsertedAmount < totalCost {
		return errors.New("insufficient funds")
	}

	vm.InsertedAmount -= totalCost

	// Update inventory with bounds checking
	newQuantity := vm.Inventory[productId].Quantity - quantity
	if newQuantity < 0 {
		return errors.New("inventory error: quantity would become negative")
	}

	vm.Inventory[productId].Quantity = newQuantity
	fmt.Printf("Dispensed %d units of product %s (remaining: %d)\n", quantity, productId, newQuantity)

	vm.State = &ReturnChange{}
	return nil
}

func (i *Dispense) ReturnChange(vm *vendingmachine.VendingMachine, productId string, quantity int) {
	fmt.Println("cannot return before dispensing")
}

func (i *Dispense) CancelProduct(vm *vendingmachine.VendingMachine, productId string) error {
	return errors.New("cannot cancel the product")
}
