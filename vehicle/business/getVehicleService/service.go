package getVehicleService

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/vehicle/business"
)

// Service handles getting vehicle by UUID
type Service struct {
	vehicleRepository business.Repository
}

// New creates a new GetVehicle service
func New(vehicleRepo business.Repository) *Service {
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
