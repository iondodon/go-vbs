package business

import (
	"context"

	"github.com/google/uuid"
	bookingController "github.com/iondodon/go-vbs/booking/controller"
	"github.com/iondodon/go-vbs/vehicle/domain"
)

// VehicleRepository defines what the vehicle business logic needs from data access
type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error)
}
