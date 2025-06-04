package availabilityService

import (
	"context"

	"github.com/google/uuid"
	bookingController "github.com/iondodon/go-vbs/booking/controller"
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

func (s *Service) CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error) {
	hasBookedDates, err := s.vehicleRepository.VehicleHasBookedDatesOnPeriod(ctx, vUUID, period)

	if err != nil {
		return false, err
	}

	return !hasBookedDates, nil
}
