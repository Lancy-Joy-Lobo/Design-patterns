package floor

import (
	"fmt"
	parkingSpace "parkinglot/parkingSpace"
	"parkinglot/vehicle"
)

// Floor represents a parking floor with multiple parking spots.
type Floor struct {
	FloorNo      int
	ParkingSpots []*parkingSpace.ParkingSpot
}

// NewParkingFloor creates a new floor with the specified number of spots per type.
func NewParkingFloor(floorNo int, largeSpots, compactSpots, motorCycleSpots int) *Floor {
	parkingSpots := []*parkingSpace.ParkingSpot{}

	for i := 0; i < largeSpots; i++ {
		parkingSpots = append(parkingSpots, &parkingSpace.ParkingSpot{
			SpotId:      fmt.Sprintf("LARGE-%d-%d", floorNo, i),
			VehicleType: vehicle.LARGE,
		})
	}

	for i := 0; i < compactSpots; i++ {
		parkingSpots = append(parkingSpots, &parkingSpace.ParkingSpot{
			SpotId:      fmt.Sprintf("COMPACT-%d-%d", floorNo, i),
			VehicleType: vehicle.COMPACT,
		})
	}

	for i := 0; i < motorCycleSpots; i++ {
		parkingSpots = append(parkingSpots, &parkingSpace.ParkingSpot{
			SpotId:      fmt.Sprintf("MOTORCYCLE-%d-%d", floorNo, i),
			VehicleType: vehicle.MOTORCYCLE,
		})
	}

	return &Floor{
		FloorNo:      floorNo,
		ParkingSpots: parkingSpots,
	}
}
