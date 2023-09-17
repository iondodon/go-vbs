package usecase

import (
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"

	"github.com/google/uuid"
)

type GetVehicleUseCase interface {
	ByUUID(vUUID uuid.UUID) (*domain.Vehicle, error)
}

type getVehicleUseCase struct {
	vehicleRepository repository.VehicleRepository
}

func NewGetVehicleUseCase(vehicleRepository repository.VehicleRepository) GetVehicleUseCase {
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
