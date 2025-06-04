package business

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	bookingIn "github.com/iondodon/go-vbs/booking/in"
	"github.com/iondodon/go-vbs/domain"
)

// BookVehicleUseCase defines the interface for booking a vehicle
type BookVehicleUseCase interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period bookingIn.DatePeriodDTO) error
}

// GetAllBookingsUseCase defines the interface for getting all bookings
type GetAllBookingsUseCase interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}
