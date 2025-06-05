package getVehicleService

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness"
	"github.com/iondodon/go-vbs/vehicle/vehicleDomain"
)

// Service handles getting vehicle by UUID
type Service struct {
	vehicleRepository vehicleBusiness.VehicleRepository
}

// New creates a new get vehicle service
func New(vehicleRepo vehicleBusiness.VehicleRepository) *Service {
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
