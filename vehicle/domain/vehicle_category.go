package vehicleDomain

type VehicleType string

const (
	SmallCar  VehicleType = "SMALL_CAR"
	EstateCar VehicleType = "ESTATE_CAR"
	Van       VehicleType = "VAN"
)

type VehicleCategory struct {
	ID          int64       `json:"id"`
	VehicleType VehicleType `json:"vehicle_type"`
	PricePerDay float32     `json:"price_per_day"`
}
