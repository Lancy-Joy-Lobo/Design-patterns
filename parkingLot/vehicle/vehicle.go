package vehicle

// VehicleType represents the type of a vehicle.
type VehicleType int

const (
	// MOTORCYCLE represents a motorcycle type vehicle.
	MOTORCYCLE VehicleType = iota
	// COMPACT represents a compact car type vehicle.
	COMPACT
	// TRUCK represents a truck type vehicle.
	TRUCK
	// LARGE represents a large vehicle, like a bus or RV.
	LARGE
)

// Vehicle represents a vehicle in the parking lot.
type Vehicle struct {
	// LicensePlate is the vehicle's license plate number.
	LicensePlate string
	// Type indicates the type of the vehicle.
	Type VehicleType
}

// GetVehicleDetails returns the details of a vehicle.
func GetVehicleDetails(licensePlate string, vehicleType VehicleType) *Vehicle {
	return &Vehicle{
		LicensePlate: licensePlate,
		Type:         vehicleType,
	}
}
