package getVehicle

import (
	"context"

	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/domain"

	"github.com/google/uuid"
)

type UseCase interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}

type Service struct {
	vehicleRepository business.VehicleRepository
}

func New(vehicleRepo business.VehicleRepository) *Service {
	return &Service{
		vehicleRepository: vehicleRepo,
	}
}

func (gvuc *Service) ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
