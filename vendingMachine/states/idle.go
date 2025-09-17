package states

import (
	"vendingMachine/vendingmachine"
)

type IdleState struct{}

var _ vendingmachine.VendingMachineInterface = (*IdleState)(nil)

func (i *IdleState) InsertMoney(vm *vendingmachine.VendingMachine, money map[vendingmachine.Denominations]int) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	for coin, q := range money {
		vm.Denominations[coin] += q
		vm.TotalAmount[coin] += q // Also update TotalAmount for change calculations
		vm.InsertedAmount += coin.GetMoney() * q
	}
	vm.State = &HasMoney{}
}

func (i *IdleState) SelectProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return vendingmachine.NewVendingMachineError(vendingmachine.ErrInvalidState, "insert money first")
}

func (i *IdleState) DispenseProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return vendingmachine.NewVendingMachineError(vendingmachine.ErrInvalidState, "insert money first")
}

func (i *IdleState) ReturnChange(vm *vendingmachine.VendingMachine, productId string, quantity int) {
	// No action needed in idle state
}

func (i *IdleState) CancelProduct(vm *vendingmachine.VendingMachine, productIt string) error {
	return vendingmachine.NewVendingMachineError(vendingmachine.ErrInvalidState, "insert money first")
}
