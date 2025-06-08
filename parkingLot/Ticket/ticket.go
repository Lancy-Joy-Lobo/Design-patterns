package ticket

import (
	"fmt"
	vehicle "parkinglot/vehicle"
	"time"
)

// Ticket represents a parking ticket.
type Ticket struct {
	TicketId   string
	Entry      time.Time
	Exit       time.Time
	SpotId     string
	Vehicle    vehicle.Vehicle
	TotalPrice float64
}

// GenerateTicket creates a new ticket for a vehicle and spot.
func GenerateTicket(spotId string, vehicle vehicle.Vehicle) *Ticket {
	return &Ticket{
		TicketId: fmt.Sprintf("%d", time.Now().UnixNano()),
		Entry:    time.Now(),
		SpotId:   spotId,
		Vehicle:  vehicle,
	}
}
