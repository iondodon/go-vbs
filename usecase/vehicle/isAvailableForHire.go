package vehicle

import (
	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	vehRepoPkg "github.com/iondodon/go-vbs/repository/vehicle"
)

type IsAvailableForHire interface {
	CheckForPeriod(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

type isAvailableForHire struct {
	vehRepo vehRepoPkg.VehicleRepository
}

func NewIsAvailableForHire(vrp vehRepoPkg.VehicleRepository) IsAvailableForHire {
	return &isAvailableForHire{vehRepo: vrp}
}

func (us *isAvailableForHire) CheckForPeriod(vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
