package domain

import "github.com/google/uuid"

type FuelType string

const (
	Petrol FuelType = "PETROL"
	Diesel FuelType = "DIESEL"
)

type Vehicle struct {
	ID                 int64
	UUID               uuid.UUID
	RegistrationNumber string
	Make               string
	Model              string
	FuelType           FuelType
	VehicleCategory    VehicleCategory
	Bookings           []Booking
}
