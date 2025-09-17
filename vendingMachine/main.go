package main

import (
	"fmt"
	"vendingMachine/states"
	"vendingMachine/vendingmachine"
)

func main() {
	fmt.Println("Vending Machine - State Pattern Implementation")

	// Initialize vending machine
	vm := vendingmachine.NewVendingMachine()
	vm.SetState(&states.IdleState{})

	// Setup machine
	setupMachine(vm)

	// Test scenarios
	fmt.Println("\n1. Successful Purchase")
	demonstrateSuccessfulPurchase(vm)

	fmt.Println("\n2. Purchase with Change")
	demonstrateReturnChange(vm)

	fmt.Println("\n3. Insufficient Funds")
	demonstrateInsufficientFunds(vm)

	fmt.Println("\n4. Product Not Found")
	demonstrateProductNotFound(vm)

	fmt.Println("\n5. Insufficient Stock")
	demonstrateInsufficientStock(vm)

	fmt.Println("\n6. Admin Operations")
	demonstrateAdminOperations(vm)

	fmt.Println("\n7. Current Status")
	displayMachineStatus(vm)
}

func setupMachine(vm *vendingmachine.VendingMachine) {
	// Add change denominations
	vm.AddChangeToMachine(map[vendingmachine.Denominations]int{
		vendingmachine.FIVE:   20,
		vendingmachine.TEN:    15,
		vendingmachine.TWENTY: 10,
	})

	// Add products
	vm.AddProduct("COKE", "Coca Cola", 25, 10)
	vm.AddProduct("PEPSI", "Pepsi Cola", 30, 8)
	vm.AddProduct("WATER", "Bottled Water", 15, 15)

	fmt.Println("Machine initialized with products and change")
}

func demonstrateSuccessfulPurchase(vm *vendingmachine.VendingMachine) {
	vm.SetState(&states.IdleState{})

	money := map[vendingmachine.Denominations]int{
		vendingmachine.TWENTY: 1,
		vendingmachine.FIVE:   1,
	}
	vm.State.InsertMoney(vm, money)

	err := vm.ProcessTransaction("COKE", 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully purchased Coca Cola")
	}
}

func demonstrateReturnChange(vm *vendingmachine.VendingMachine) {
	vm.SetState(&states.IdleState{})

	// Insert $20 for a $15 water bottle (should return $5 change)
	money := map[vendingmachine.Denominations]int{
		vendingmachine.TWENTY: 1,
	}
	vm.State.InsertMoney(vm, money)
	fmt.Printf("Inserted: $%d\n", vm.InsertedAmount)

	err := vm.ProcessTransaction("WATER", 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully purchased Water with change returned")
	}
}

func demonstrateInsufficientFunds(vm *vendingmachine.VendingMachine) {
	vm.SetState(&states.IdleState{})

	money := map[vendingmachine.Denominations]int{
		vendingmachine.TWENTY: 1, // $20 inserted, Pepsi costs $30
	}
	vm.State.InsertMoney(vm, money)

	err := vm.ProcessTransaction("PEPSI", 1)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}

func demonstrateProductNotFound(vm *vendingmachine.VendingMachine) {
	vm.SetState(&states.IdleState{})

	money := map[vendingmachine.Denominations]int{
		vendingmachine.TWENTY: 2,
	}
	vm.State.InsertMoney(vm, money)

	err := vm.ProcessTransaction("SPRITE", 1)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}

func demonstrateInsufficientStock(vm *vendingmachine.VendingMachine) {
	vm.SetState(&states.IdleState{})

	money := map[vendingmachine.Denominations]int{
		vendingmachine.TWENTY: 10,
	}
	vm.State.InsertMoney(vm, money)

	err := vm.ProcessTransaction("WATER", 20) // Try to buy more than available
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}

func demonstrateAdminOperations(vm *vendingmachine.VendingMachine) {
	// Restock
	err := vm.RestockProduct("COKE", 5)
	if err != nil {
		fmt.Printf("Restock failed: %v\n", err)
	} else {
		fmt.Println("Restocked COKE (+5 units)")
	}

	// Update price
	err = vm.UpdateProductPrice("PEPSI", 35)
	if err != nil {
		fmt.Printf("Price update failed: %v\n", err)
	} else {
		fmt.Println("Updated PEPSI price to $35")
	}
}

func displayMachineStatus(vm *vendingmachine.VendingMachine) {
	inventory := vm.GetInventoryStatus()
	fmt.Println("Inventory:")
	for _, product := range inventory {
		fmt.Printf("  %s: %d units @ $%d\n", product.Name, product.Quantity, product.Price)
	}

	totalValue, _ := vm.GetMachineStatus()
	fmt.Printf("Total cash: $%d\n", totalValue)
}
