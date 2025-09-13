package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Direction of elevator
type Direction int

const (
	Up Direction = iota
	Down
	Idle
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	default:
		return "Idle"
	}
}

// Elevator struct
type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    Direction
	UpRequests   []int
	DownRequests []int
	mutex        sync.Mutex
}

// AddRequest adds a floor request
func (e *Elevator) AddRequest(floor int, direction Direction) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if direction == Up {
		for _, f := range e.UpRequests {
			if f == floor {
				return
			}
		}
		e.UpRequests = append(e.UpRequests, floor)
	} else {
		for _, f := range e.DownRequests {
			if f == floor {
				return
			}
		}
		e.DownRequests = append(e.DownRequests, floor)
	}
}

// MoveOneStep moves elevator one floor at a time
func (e *Elevator) MoveOneStep() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Decide direction dynamically
	if len(e.UpRequests) > 0 {
		e.Direction = Up
		sort.Ints(e.UpRequests)
		if e.CurrentFloor < e.UpRequests[0] {
			e.CurrentFloor++
		} else if e.CurrentFloor == e.UpRequests[0] {
			fmt.Printf("Elevator %d arrived at floor %d (Up)\n", e.ID, e.CurrentFloor)
			e.UpRequests = e.UpRequests[1:]
		}
	} else if len(e.DownRequests) > 0 {
		e.Direction = Down
		sort.Sort(sort.Reverse(sort.IntSlice(e.DownRequests)))
		if e.CurrentFloor > e.DownRequests[0] {
			e.CurrentFloor--
		} else if e.CurrentFloor == e.DownRequests[0] {
			fmt.Printf("Elevator %d arrived at floor %d (Down)\n", e.ID, e.CurrentFloor)
			e.DownRequests = e.DownRequests[1:]
		}
	} else {
		e.Direction = Idle
	}
}

// Start runs the elevator in its own goroutine
func (e *Elevator) Start(wg *sync.WaitGroup, stopChan chan struct{}) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				e.MoveOneStep()
			}
		}
	}()
}

// ElevatorSystem manages multiple elevators
type ElevatorSystem struct {
	Elevators []*Elevator
	mutex     sync.Mutex
}

// NewElevatorSystem initializes elevators
func NewElevatorSystem(num int) *ElevatorSystem {
	elevators := make([]*Elevator, num)
	for i := 0; i < num; i++ {
		elevators[i] = &Elevator{ID: i + 1}
	}
	return &ElevatorSystem{Elevators: elevators}
}

// RequestElevator assigns the best elevator
func (es *ElevatorSystem) RequestElevator(floor int, direction Direction) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	var best *Elevator
	minDist := int(^uint(0) >> 1) // Max int

	for _, e := range es.Elevators {
		e.mutex.Lock()
		dist := abs(e.CurrentFloor - floor)

		// Prefer idle elevators
		if e.Direction == Idle {
			dist = dist / 2
		}

		// Prefer elevators moving in same direction and can pick up
		if e.Direction == direction &&
			((direction == Up && e.CurrentFloor <= floor) ||
				(direction == Down && e.CurrentFloor >= floor)) {
			dist = dist / 3
		}

		if dist < minDist {
			minDist = dist
			best = e
		}
		e.mutex.Unlock()
	}

	if best != nil {
		best.AddRequest(floor, direction)
		fmt.Printf("Assigned floor %d (%s) to Elevator %d\n", floor, direction, best.ID)
	}
}

// Helper
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// ---------------- Demo ----------------
func main() {
	system := NewElevatorSystem(2)
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	// Start all elevators
	for _, e := range system.Elevators {
		e.Start(&wg, stopChan)
	}

	// Simulate requests
	system.RequestElevator(3, Up)
	time.Sleep(300 * time.Millisecond)
	system.RequestElevator(7, Down)
	time.Sleep(300 * time.Millisecond)
	system.RequestElevator(2, Up)
	time.Sleep(300 * time.Millisecond)

	// Simulate internal request
	system.Elevators[0].AddRequest(9, Up)
	system.Elevators[1].AddRequest(1, Down)

	// Let elevators process for a few seconds
	time.Sleep(5 * time.Second)
	close(stopChan)
	wg.Wait()

	fmt.Println("Simulation finished")
}
