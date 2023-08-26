package usecase

import (
	"github.com/google/uuid"
	"go-vbs/domain"
	"go-vbs/repository"
)

type GetVehicle interface {
	ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error)
}

type getVehicle struct {
	vehicleRepository repository.VehicleRepository
}

func NewGetVehicleUseCase(vehicleRepository repository.VehicleRepository) GetVehicle {
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
