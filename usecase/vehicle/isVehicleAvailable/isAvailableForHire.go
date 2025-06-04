package isVehicleAvailable

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/repository/vehicleRepository"
)

type IsAvailableForHireInterface interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

type IsAvailableForHire struct {
	vehRepo vehicleRepository.VehicleRepositoryInterface
}

func New(vehicleRepo vehicleRepository.VehicleRepositoryInterface) IsAvailableForHireInterface {
	return &IsAvailableForHire{
		vehRepo: vehicleRepo,
	}
}

func (us *IsAvailableForHire) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
