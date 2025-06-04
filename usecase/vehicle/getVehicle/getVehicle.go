package getVehicle

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository/vehicleRepository"

	"github.com/google/uuid"
)

type GetVehicleInterface interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

type GetVehicle struct {
	vehicleRepository vehicleRepository.VehicleRepositoryInterface
}

func New(vehicleRepo vehicleRepository.VehicleRepositoryInterface) GetVehicleInterface {
	return &GetVehicle{
		vehicleRepository: vehicleRepo,
	}
}

func (gvuc *GetVehicle) ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
