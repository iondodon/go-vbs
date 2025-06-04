package getVehicle

import (
	"context"

	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/domain"

	"github.com/google/uuid"
)

type Service struct {
	vehicleRepository business.VehicleRepository
}

func New(vehicleRepo business.VehicleRepository) *Service {
	return &Service{
		vehicleRepository: vehicleRepo,
	}
}

func (s *Service) ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
