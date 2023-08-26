package usecase

import (
	"github.com/google/uuid"
	"go-vbs/dto"
	"go-vbs/repository"
)

type BookVehicle interface {
	ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO)
}

type bookVehicle struct {
	vehicleRepository repository.VehicleRepository
}

func (*bookVehicle) ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) {
	return
}
