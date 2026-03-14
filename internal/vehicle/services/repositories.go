package vehicleServices

import (
	"context"

	"github.com/google/uuid"
	bookingController "github.com/iondodon/go-vbs/internal/booking/controller"
	vehicleDomain "github.com/iondodon/go-vbs/internal/vehicle/domain"
)

// VehicleRepository defines what the vehicle services layer needs from data access
type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*vehicleDomain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error)
}
