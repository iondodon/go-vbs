package domain

import "github.com/google/uuid"

type FuelType string

const (
	Petrol FuelType = "PETROL"
	Diesel FuelType = "DIESEL"
)

type Vehicle struct {
	ID                 int64            `json:"id"`
	UUID               uuid.UUID        `json:"uuid"`
	RegistrationNumber string           `json:"registartion_number"`
	Make               string           `json:"make"`
	Model              string           `json:"model"`
	FuelType           FuelType         `json:"fuel_type"`
	VehicleCategory    *VehicleCategory `json:"vehicle_category"`
	Bookings           []*Booking       `json:"-"`
}
