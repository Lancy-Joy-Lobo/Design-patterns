package states

import (
	"errors"
	"fmt"
	"vendingMachine/vendingmachine"
)

type ReturnChange struct{}

var _ vendingmachine.VendingMachineInterface = (*ReturnChange)(nil)

func (i *ReturnChange) InsertMoney(vm *vendingmachine.VendingMachine, money map[vendingmachine.Denominations]int) {
	fmt.Println("entering money is not allowed")
}

func (i *ReturnChange) SelectProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return errors.New("selecting product is not allowed")
}

func (i *ReturnChange) DispenseProduct(vm *vendingmachine.VendingMachine, productId string, quantity int) error {
	return errors.New("returning product is not allowed")
}

func (i *ReturnChange) ReturnChange(vm *vendingmachine.VendingMachine, productId string, quantity int) {
	if vm.InsertedAmount > 0 {
		change, err := vm.ReturnChange()
		if err != nil {
			fmt.Printf("Error returning change: %s\n", err.Error())
		} else {
			for coin, q := range change {
				if q > 0 {
					fmt.Printf("Change: %d x %d denomination\n", q, coin.GetMoney())
				}
			}
		}
		vm.InsertedAmount = 0
	}

	// Clear the denominations tracking for this transaction
	vm.Denominations = make(map[vendingmachine.Denominations]int)
	vm.State = &IdleState{}
}

func (i *ReturnChange) CancelProduct(vm *vendingmachine.VendingMachine, productId string) error {
	return errors.New("cannot cancel the product")
}
