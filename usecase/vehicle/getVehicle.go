package vehicle

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

	"github.com/google/uuid"
)

type GetVehicleInterface interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

//gobok:constructor
//ctxboot:component
type GetVehicle struct {
	vehicleRepository vehRepo.VehicleRepositoryInterface `ctxboot:"inject"`
}

func (gvuc *GetVehicle) ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
