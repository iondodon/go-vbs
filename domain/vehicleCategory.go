package domain

type VehicleType string

const (
	SmallCar  VehicleType = "SMALL_CAR"
	EstateCar VehicleType = "ESTATE_CAR"
	Van       VehicleType = "VAN"
)

type VehicleCategory struct {
	ID          int64
	VehicleType VehicleType
	PricePerDay float32
}
