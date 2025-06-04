package business

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
)

// BookVehicleUseCase defines the interface for booking a vehicle
type BookVehicleUseCase interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

// GetAllBookingsUseCase defines the interface for getting all bookings
type GetAllBookingsUseCase interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}
