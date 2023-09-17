package booking

import (
	"github.com/iondodon/go-vbs/dto"
	custRepo "github.com/iondodon/go-vbs/repository/customer"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

	"github.com/google/uuid"
)

type BookVehicleUseCase interface {
	ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

type bookVehicleUseCase struct {
	vehRepo  vehRepo.VehicleRepository
	custRepo custRepo.CustomerRepository
}

func NewBookVehicleUseCase(
	vrp vehRepo.VehicleRepository,
	crp custRepo.CustomerRepository,
) BookVehicleUseCase {
	return &bookVehicleUseCase{vehRepo: vrp, custRepo: crp}
}

func (uc *bookVehicleUseCase) ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
	_, err := uc.vehRepo.FindByUUID(vehicleUUID)
	if err != nil {
		return err
	}

	_, err = uc.custRepo.FindByUUID(customerUID)
	if err != nil {
		return err
	}

	return nil
}
