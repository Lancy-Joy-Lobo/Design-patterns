package main

import (
	"fmt"
	"log"
	ticket "parkinglot/Ticket"
	parkingLot "parkinglot/parkingLot"
	payment "parkinglot/payment"
	price "parkinglot/price"
	vehicle "parkinglot/vehicle"
	"time"
)

func main() {
	pl := parkingLot.NewParkingLot()

	vehicles := []vehicle.Vehicle{
		{LicensePlate: "ABCD", Type: vehicle.COMPACT},
		{LicensePlate: "DEFG", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY1", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY2", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY3", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY4", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY5", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY6", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY7", Type: vehicle.COMPACT},
		{LicensePlate: "GTHY8", Type: vehicle.COMPACT},
	}

	for i := 0; i < len(vehicles); i++ {
		floorNo, spotId, err := pl.ParkVehicle(vehicle.COMPACT, &vehicles[i])
		if err != nil {
			log.Printf("Cannot park vehicle %s: %v\n", vehicles[i].LicensePlate, err)
			continue
		}
		fmt.Printf("Park your vehicle in floor %d spot %s\n", floorNo, spotId)
		t := ticket.GenerateTicket(spotId, vehicles[i])
		pl.Tickets[vehicles[i].LicensePlate] = t
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < len(pl.Floor); i++ {
		for j := 0; j < len(pl.Floor[i].ParkingSpots); j++ {
			spot := pl.Floor[i].ParkingSpots[j]
			if spot.IsOccupied {
				spot.UnParkVehicle(spot.SpotId)
				price.CalcualatePrice(pl.Tickets[spot.CurrentVehicle.LicensePlate])
				payment.GetPaymentMethod("CASH").Credit(pl.Tickets[spot.CurrentVehicle.LicensePlate].TotalPrice)
			}
		}
	}
}
