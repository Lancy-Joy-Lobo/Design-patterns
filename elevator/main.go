package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Direction represents the direction of elevator movement
type Direction int

const (
	IDLE Direction = iota
	UP
	DOWN
)

func (d Direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		return "IDLE"
	}
}

// ElevatorState represents the current state of an elevator
type ElevatorState int

const (
	MOVING ElevatorState = iota
	STOPPED
	MAINTENANCE
)

func (s ElevatorState) String() string {
	switch s {
	case MOVING:
		return "MOVING"
	case STOPPED:
		return "STOPPED"
	default:
		return "MAINTENANCE"
	}
}

// Request represents a floor request (internal or external)
type Request struct {
	Floor     int
	Direction Direction
	Timestamp time.Time
}

// Elevator represents a single elevator unit
type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    Direction
	State        ElevatorState
	Capacity     int
	CurrentLoad  int
	UpRequests   []int // floors requested going up
	DownRequests []int // floors requested going down
	mutex        sync.RWMutex
}

// NewElevator creates a new elevator instance
func NewElevator(id, capacity int) *Elevator {
	return &Elevator{
		ID:           id,
		CurrentFloor: 1,
		Direction:    IDLE,
		State:        STOPPED,
		Capacity:     capacity,
		CurrentLoad:  0,
		UpRequests:   make([]int, 0),
		DownRequests: make([]int, 0),
	}
}

// AddInternalRequest adds an internal floor request
func (e *Elevator) AddInternalRequest(floor int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if floor > e.CurrentFloor {
		e.UpRequests = append(e.UpRequests, floor)
		sort.Ints(e.UpRequests)
	} else if floor < e.CurrentFloor {
		e.DownRequests = append(e.DownRequests, floor)
		sort.Sort(sort.Reverse(sort.IntSlice(e.DownRequests)))
	}
}

// HasRequests checks if elevator has any pending requests
func (e *Elevator) HasRequests() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return len(e.UpRequests) > 0 || len(e.DownRequests) > 0
}

// GetNextFloor returns the next floor to visit
func (e *Elevator) GetNextFloor() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	switch e.Direction {
	case UP:
		if len(e.UpRequests) > 0 {
			return e.UpRequests[0]
		}
		if len(e.DownRequests) > 0 {
			return e.DownRequests[0]
		}
	case DOWN:
		if len(e.DownRequests) > 0 {
			return e.DownRequests[0]
		}
		if len(e.UpRequests) > 0 {
			return e.UpRequests[0]
		}
	case IDLE:
		if len(e.UpRequests) > 0 {
			return e.UpRequests[0]
		}
		if len(e.DownRequests) > 0 {
			return e.DownRequests[0]
		}
	}
	return e.CurrentFloor
}

// RemoveRequest removes a floor request when reached
func (e *Elevator) RemoveRequest(floor int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Remove from up requests
	for i, f := range e.UpRequests {
		if f == floor {
			e.UpRequests = append(e.UpRequests[:i], e.UpRequests[i+1:]...)
			break
		}
	}

	// Remove from down requests
	for i, f := range e.DownRequests {
		if f == floor {
			e.DownRequests = append(e.DownRequests[:i], e.DownRequests[i+1:]...)
			break
		}
	}
}

// Move simulates elevator movement
func (e *Elevator) Move() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !e.HasRequests() {
		e.Direction = IDLE
		e.State = STOPPED
		return
	}

	nextFloor := e.GetNextFloor()
	if nextFloor == e.CurrentFloor {
		e.State = STOPPED
		e.RemoveRequest(e.CurrentFloor)
		return
	}

	e.State = MOVING
	if nextFloor > e.CurrentFloor {
		e.Direction = UP
		e.CurrentFloor++
	} else {
		e.Direction = DOWN
		e.CurrentFloor--
	}

	// Check if we've reached the requested floor
	if e.CurrentFloor == nextFloor {
		e.State = STOPPED
		e.RemoveRequest(e.CurrentFloor)
		fmt.Printf("Elevator %d reached floor %d\n", e.ID, e.CurrentFloor)
	}
}

// GetStatus returns current elevator status
func (e *Elevator) GetStatus() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return fmt.Sprintf("Elevator %d: Floor %d, Direction: %s, State: %s, Load: %d/%d",
		e.ID, e.CurrentFloor, e.Direction, e.State, e.CurrentLoad, e.Capacity)
}

// ElevatorController manages multiple elevators
type ElevatorController struct {
	Elevators        []*Elevator
	TotalFloors      int
	ExternalRequests []Request
	mutex            sync.RWMutex
}

// NewElevatorController creates a new elevator controller
func NewElevatorController(numElevators, totalFloors, capacity int) *ElevatorController {
	controller := &ElevatorController{
		Elevators:        make([]*Elevator, numElevators),
		TotalFloors:      totalFloors,
		ExternalRequests: make([]Request, 0),
	}

	for i := 0; i < numElevators; i++ {
		controller.Elevators[i] = NewElevator(i+1, capacity)
	}

	return controller
}

// RequestElevator handles external elevator requests
func (ec *ElevatorController) RequestElevator(floor int, direction Direction) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	request := Request{
		Floor:     floor,
		Direction: direction,
		Timestamp: time.Now(),
	}
	ec.ExternalRequests = append(ec.ExternalRequests, request)

	// Find the best elevator for this request
	bestElevator := ec.findBestElevator(request)
	if bestElevator != nil {
		bestElevator.AddInternalRequest(floor)
		fmt.Printf("Assigned request (Floor: %d, Direction: %s) to Elevator %d\n",
			floor, direction, bestElevator.ID)
	}
}

// findBestElevator implements SCAN algorithm to find optimal elevator
func (ec *ElevatorController) findBestElevator(request Request) *Elevator {
	var bestElevator *Elevator
	minCost := int(^uint(0) >> 1) // Max int

	for _, elevator := range ec.Elevators {
		if elevator.State == MAINTENANCE {
			continue
		}

		cost := ec.calculateCost(elevator, request)
		if cost < minCost {
			minCost = cost
			bestElevator = elevator
		}
	}

	return bestElevator
}

// calculateCost calculates the cost of assigning a request to an elevator
func (ec *ElevatorController) calculateCost(elevator *Elevator, request Request) int {
	distance := abs(elevator.CurrentFloor - request.Floor)

	// Add penalty for direction mismatch
	penalty := 0
	if elevator.Direction != IDLE && elevator.Direction != request.Direction {
		penalty = 10
	}

	// Add penalty for high load
	loadPenalty := elevator.CurrentLoad * 2

	return distance + penalty + loadPenalty
}

// ProcessRequests handles internal elevator button presses
func (ec *ElevatorController) ProcessInternalRequest(elevatorID, floor int) {
	if elevatorID < 1 || elevatorID > len(ec.Elevators) {
		fmt.Printf("Invalid elevator ID: %d\n", elevatorID)
		return
	}

	elevator := ec.Elevators[elevatorID-1]
	elevator.AddInternalRequest(floor)
	fmt.Printf("Internal request added: Elevator %d to Floor %d\n", elevatorID, floor)
}

// Start begins the elevator system operation
func (ec *ElevatorController) Start() {
	fmt.Println("Starting Elevator System...")

	for {
		// Move all elevators
		for _, elevator := range ec.Elevators {
			elevator.Move()
		}

		// Display status
		ec.DisplayStatus()

		// Check if all elevators are idle
		allIdle := true
		for _, elevator := range ec.Elevators {
			if elevator.HasRequests() {
				allIdle = false
				break
			}
		}

		if allIdle && len(ec.ExternalRequests) == 0 {
			fmt.Println("All elevators are idle. System ready for new requests.")
			break
		}

		time.Sleep(1 * time.Second)
	}
}

// DisplayStatus shows current system status
func (ec *ElevatorController) DisplayStatus() {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	fmt.Println("\n--- Elevator System Status ---")
	for _, elevator := range ec.Elevators {
		fmt.Println(elevator.GetStatus())
	}
	fmt.Printf("Pending External Requests: %d\n", len(ec.ExternalRequests))
	fmt.Println("-------------------------------")
}

// Utility function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Building represents the building with elevator system
type Building struct {
	Controller *ElevatorController
	Floors     int
}

// NewBuilding creates a new building with elevator system
func NewBuilding(floors, numElevators, elevatorCapacity int) *Building {
	return &Building{
		Controller: NewElevatorController(numElevators, floors, elevatorCapacity),
		Floors:     floors,
	}
}

// Example usage and testing
func main() {
	// Create a building with 10 floors, 3 elevators, each with capacity of 8
	building := NewBuilding(10, 3, 8)

	// Simulate some requests
	fmt.Println("=== Elevator System Demo ===")

	// External requests (someone calling elevator)
	building.Controller.RequestElevator(5, UP)
	building.Controller.RequestElevator(3, DOWN)
	building.Controller.RequestElevator(8, UP)

	// Internal requests (someone inside elevator pressing buttons)
	building.Controller.ProcessInternalRequest(1, 7)
	building.Controller.ProcessInternalRequest(2, 1)
	building.Controller.ProcessInternalRequest(3, 9)

	// Start the system
	building.Controller.Start()

	fmt.Println("\n=== Demo Completed ===")
}
