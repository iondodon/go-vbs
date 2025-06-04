package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/controller"
	"github.com/iondodon/go-vbs/vehicle/domain"
)

// GetVehicleUseCase defines the interface for getting vehicle by UUID
type GetVehicleUseCase interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

// AvailabilityUseCase defines the interface for checking vehicle availability
type AvailabilityUseCase interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period controller.DatePeriodDTO) (bool, error)
}
