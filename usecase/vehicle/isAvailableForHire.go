package vehicle

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	vehRepoPkg "github.com/iondodon/go-vbs/repository/vehicle"
)

//gobok:constructor
//ctxboot:component
type IsAvailableForHire struct {
	vehRepo vehRepoPkg.VehicleRepository `ctxboot:"inject"`
}

func (us *IsAvailableForHire) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := us.vehRepo.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
