package vehicleServices

import (
	"context"

	"github.com/google/uuid"
	bookingController "github.com/iondodon/go-vbs/internal/booking/controller"
	vehicleDomain "github.com/iondodon/go-vbs/internal/vehicle/domain"
)

// GetVehicleUseCase defines the interface for getting vehicle by UUID
type GetVehicleUseCase interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*vehicleDomain.Vehicle, error)
}

// AvailabilityUseCase defines the interface for checking vehicle availability
type AvailabilityUseCase interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error)
}
