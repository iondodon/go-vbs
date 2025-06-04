package getVehicle

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/usecase"

	"github.com/google/uuid"
)

type Interface interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

type UseCase struct {
	vehicleRepository usecase.VehicleRepositoryInterface
}

func New(vehicleRepo usecase.VehicleRepositoryInterface) Interface {
	return &UseCase{
		vehicleRepository: vehicleRepo,
	}
}

func (gvuc *UseCase) ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
