package get

import (
	"context"

	"github.com/google/uuid"
	vehicleDomain "github.com/iondodon/go-vbs/internal/vehicle/domain"
	vehicleServices "github.com/iondodon/go-vbs/internal/vehicle/services"
)

// Service handles getting vehicle by UUID
type Service struct {
	vehicleRepository vehicleServices.VehicleRepository
}

// New creates a new get vehicle service
func New(vehicleRepo vehicleServices.VehicleRepository) *Service {
	return &Service{
		vehicleRepository: vehicleRepo,
	}
}

func (s *Service) ByUUID(ctx context.Context, vUUID uuid.UUID) (*vehicleDomain.Vehicle, error) {
	vehicle, err := s.vehicleRepository.FindByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
