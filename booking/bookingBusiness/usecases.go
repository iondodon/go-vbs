package bookingBusiness

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/bookingController"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
)

// BookVehicleUseCase defines the interface for booking a vehicle
type BookVehicleUseCase interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period bookingController.DatePeriodDTO) error
}

// GetAllBookingsUseCase defines the interface for getting all bookings
type GetAllBookingsUseCase interface {
	Execute(ctx context.Context) ([]bookingDomain.Booking, error)
}

// Cross-domain service dependencies (services from other domains that booking uses)
type VehicleAvailabilityService interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error)
}
