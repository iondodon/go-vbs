package isVehicleAvailable

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/usecase"
)

type Interface interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

type UseCase struct {
	vehRepo usecase.VehicleRepositoryInterface
}

func New(vehicleRepo usecase.VehicleRepositoryInterface) Interface {
	return &UseCase{
		vehRepo: vehicleRepo,
	}
}

func (us *UseCase) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
