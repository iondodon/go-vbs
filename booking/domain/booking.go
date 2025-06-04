package domain

import (
	"github.com/google/uuid"
	customerDomain "github.com/iondodon/go-vbs/customer/domain"
	vehicleDomain "github.com/iondodon/go-vbs/vehicle/domain"
)

type Booking struct {
	ID           int64
	UUID         uuid.UUID
	BookingDates []*BookingDate
	Vehicle      *vehicleDomain.Vehicle
	Customer     *customerDomain.Customer
}
