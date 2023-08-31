package domain

import "github.com/google/uuid"

type Booking struct {
	ID           int64
	UUID         uuid.UUID
	BookingDates []BookingDate
	Vehicle      Vehicle
	Customer     Customer
}
