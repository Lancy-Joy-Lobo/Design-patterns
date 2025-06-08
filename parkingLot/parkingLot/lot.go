package parkinglot

import (
	"errors"
	ticket "parkinglot/Ticket"
	floor "parkinglot/floor"
	"parkinglot/vehicle"
)

type ParkingLot struct {
	Floor   []floor.Floor
	Tickets map[string]*ticket.Ticket
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		Floor: []floor.Floor{
			*floor.NewParkingFloor(1, 2, 2, 5),
			*floor.NewParkingFloor(2, 5, 2, 5),
		},
		Tickets: make(map[string]*ticket.Ticket),
	}
}

func (p *ParkingLot) ParkVehicle(vehicleType vehicle.VehicleType, vehicle *vehicle.Vehicle) (int, string, error) {
	for i := 0; i < len(p.Floor); i++ {
		for j := 0; j < len(p.Floor[i].ParkingSpots); j++ {
			spot := p.Floor[i].ParkingSpots[j]
			if !spot.IsOccupied && spot.VehicleType == vehicleType {
				spot.ParkVehicle(vehicleType, vehicle)
				return i, spot.SpotId, nil
			}
		}
	}
	return 0, "", errors.New("no spots are available, cannot park vehicle here")
}

func (p *ParkingLot) UnParkVehicle(floorNo int, spotId string) error {
	for j := 0; j < len(p.Floor[floorNo].ParkingSpots); j++ {
		if p.Floor[floorNo].ParkingSpots[j].SpotId == spotId {
			success := p.Floor[floorNo].ParkingSpots[j].UnParkVehicle(spotId)
			if success {
				return nil
			}
			return errors.New("failed to unpark vehicle")
		}
	}
	return errors.New("spot not found")
}
