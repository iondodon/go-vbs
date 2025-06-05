package bookingDomain

import (
	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/customer/customerDomain"
	"github.com/iondodon/go-vbs/vehicle/vehicleDomain"
)

type Booking struct {
	ID           int64
	UUID         uuid.UUID
	BookingDates []*BookingDate
	Vehicle      *vehicleDomain.Vehicle
	Customer     *customerDomain.Customer
}
