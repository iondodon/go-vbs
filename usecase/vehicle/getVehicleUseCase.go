package vehicle

import (
	"github.com/iondodon/go-vbs/domain"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

	"github.com/google/uuid"
)

type GetVehicleUseCase interface {
	ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error)
}

type getVehicleUseCase struct {
	vehicleRepository vehRepo.VehicleRepository
}

func NewGetVehicleUseCase(vehicleRepository vehRepo.VehicleRepository) GetVehicleUseCase {
	return &getVehicleUseCase{
		vehicleRepository: vehicleRepository,
	}
}

func (gvuc *getVehicleUseCase) ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
