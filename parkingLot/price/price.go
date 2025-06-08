package price

import (
	"fmt"
	ticket "parkinglot/Ticket"
	"time"
)

// CalcualatePrice calculates and sets the total price for a ticket.
func CalcualatePrice(ticket *ticket.Ticket) {
	ticket.Exit = time.Now()
	ticket.TotalPrice = time.Since(ticket.Entry).Seconds() * 10
	fmt.Printf("total price is %.2f\n", ticket.TotalPrice)
}
