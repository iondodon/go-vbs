package usecase

import (
	"go-vbs/dto"
	"go-vbs/repository"

	"github.com/google/uuid"
)

type BookVehicleUseCase interface {
	ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO)
}

type bookVehicleUseCase struct {
	vehicleRepository repository.VehicleRepository
}

func (*bookVehicleUseCase) ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) {
	return
}
