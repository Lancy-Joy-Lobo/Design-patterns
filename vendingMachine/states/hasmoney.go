package states

import (
	"fmt"
	"vendingMachine/vendingmachine"
)

type HasMoney struct{}

var _ vendingmachine.VendingMachineInterface = (*HasMoney)(nil)

func (i *HasMoney) InsertMoney(vm *vendingmachine.VendingMachine, money map[vendingmachine.Denominations]int) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	for coin, q := range money {
		vm.Denominations[coin] += q
		vm.TotalAmount[coin] += q // Also update TotalAmount for change calculations
		vm.InsertedAmount += coin.GetMoney() * q
	}
}

func (i *HasMoney) SelectProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if err := vm.ValidateProductId(productId); err != nil {
		return err
	}

	if err := vm.ValidateQuantity(quantity); err != nil {
		return err
	}

	product, exist := vm.Inventory[productId]
	if !exist {
		return vendingmachine.NewVendingMachineError(vendingmachine.ErrProductNotFound, "product doesn't exist")
	}

	if quantity > product.Quantity {
		return vendingmachine.NewVendingMachineError(vendingmachine.ErrInsufficientStock, "selected quantity is more than what is available")
	}

	totalCost := quantity * product.Price
	if vm.InsertedAmount < totalCost {
		diff := totalCost - vm.InsertedAmount
		return vendingmachine.NewVendingMachineError(vendingmachine.ErrInsufficientFunds, fmt.Sprintf("insufficient balance: insert %d more", diff))
	}

	// Check if machine can provide change
	changeAmount := vm.InsertedAmount - totalCost
	if changeAmount > 0 && !vm.CanProvideChange(changeAmount) {
		return vendingmachine.NewVendingMachineError(vendingmachine.ErrInsufficientChange, "exact change required - machine cannot provide sufficient change")
	}

	vm.State = &Dispense{}
	return nil
}

func (i *HasMoney) DispenseProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return vendingmachine.NewVendingMachineError(vendingmachine.ErrInvalidState, "select the product first")
}

func (i *HasMoney) ReturnChange(vm *vendingmachine.VendingMachine, productId string, quantity int) {
	// No action needed
}

func (i *HasMoney) CancelProduct(vm *vendingmachine.VendingMachine, productId string) error {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if vm.InsertedAmount > 0 {
		// Return all inserted money without deducting from TotalAmount
		for denom, count := range vm.Denominations {
			vm.TotalAmount[denom] -= count
		}
		// Clear denominations tracking
		vm.Denominations = make(map[vendingmachine.Denominations]int)
	}
	vm.InsertedAmount = 0
	vm.State = &IdleState{}
	return nil
}
