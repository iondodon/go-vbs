// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"time"
)

type Booking struct {
	ID         interface{}
	Uuid       interface{}
	VehicleID  int64
	CustomerID int64
}

type BookingDate struct {
	ID   interface{}
	Time time.Time
}

type BookingsBookingdate struct {
	BookingID     int64
	BookingdateID int64
}

type Customer struct {
	ID       interface{}
	Uuid     interface{}
	Username string
}

type Vehicle struct {
	ID                 interface{}
	Uuid               interface{}
	RegistrationNumber string
	Make               string
	Model              string
	FuelType           string
	CategoryID         int64
}

type VehicleCategory struct {
	ID          interface{}
	Category    string
	PricePerDay float64
}
