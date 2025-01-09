package vehicle

import (
	"github.com/iondodon/go-vbs/domain"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

	"github.com/google/uuid"
)

type GetVehicle interface {
	ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error)
}

type getVehicle struct {
	vehicleRepository vehRepo.VehicleRepository
}

func NewGetVehicle(vehicleRepository vehRepo.VehicleRepository) GetVehicle {
	return &getVehicle{
		vehicleRepository: vehicleRepository,
	}
}

func (gvuc *getVehicle) ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error) {
	vehicle, err := gvuc.vehicleRepository.FindByUUID(vUUID)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
