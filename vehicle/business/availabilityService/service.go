package availabilityService

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/vehicle/business"
)

// Service handles vehicle availability checks
type Service struct {
	vehicleRepository business.VehicleRepository
}

// New creates a new availability service
func New(vehicleRepo business.VehicleRepository) *Service {
	return &Service{
		vehicleRepository: vehicleRepo,
	}
}

func (s *Service) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := s.vehicleRepository.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
