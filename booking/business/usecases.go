package business

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/controller"
	"github.com/iondodon/go-vbs/booking/domain"
)

// BookVehicleUseCase defines the interface for booking a vehicle
type BookVehicleUseCase interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period controller.DatePeriodDTO) error
}

// GetAllBookingsUseCase defines the interface for getting all bookings
type GetAllBookingsUseCase interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

// Cross-domain service dependencies (services from other domains that booking uses)
type VehicleAvailabilityService interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period controller.DatePeriodDTO) (bool, error)
}
