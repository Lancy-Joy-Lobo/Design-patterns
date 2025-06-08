package parkingSpace

import (
	"errors"
	"parkinglot/vehicle"
	"sync"
)

// ParkingSpot represents a single parking spot.
type ParkingSpot struct {
	SpotId         string
	IsOccupied     bool
	VehicleType    vehicle.VehicleType
	CurrentVehicle *vehicle.Vehicle
	mutex          sync.Mutex
}

// GetSpot checks if the spot is available for the given vehicle type.
func (p *ParkingSpot) GetSpot(vehicleType vehicle.VehicleType) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.IsOccupied && p.VehicleType == vehicleType {
		return p.SpotId, nil
	}
	return "", errors.New("no spot available in the current floor")
}

// ParkVehicle parks a vehicle in the spot.
func (p *ParkingSpot) ParkVehicle(vehicleType vehicle.VehicleType, vehicle *vehicle.Vehicle) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.IsOccupied = true
	p.CurrentVehicle = vehicle
}

// UnParkVehicle un-parks the vehicle from the spot.
func (p *ParkingSpot) UnParkVehicle(spotId string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.SpotId == spotId {
		p.IsOccupied = false
		p.CurrentVehicle = nil
		return true
	}
	return false
}
