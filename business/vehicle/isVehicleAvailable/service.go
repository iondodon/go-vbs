package isVehicleAvailable

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/dto"
)

type Service struct {
	vehRepo business.VehicleRepository
}

func New(vehicleRepo business.VehicleRepository) *Service {
	return &Service{
		vehRepo: vehicleRepo,
	}
}

func (us *Service) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
