package vehicle

import (
	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	vehRepoPkg "github.com/iondodon/go-vbs/repository/vehicle"
)

type IsAvailableForHireUseCase interface {
	CheckForPeriod(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

type isAvailableForHireUseCase struct {
	vehRepo vehRepoPkg.VehicleRepository
}

func NewIsAvailableForHireUseCase(vrp vehRepoPkg.VehicleRepository) IsAvailableForHireUseCase {
	return &isAvailableForHireUseCase{vehRepo: vrp}
}

func (us *isAvailableForHireUseCase) CheckForPeriod(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
