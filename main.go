package main

import (
	"fmt"
	"github.com/google/uuid"
	"go-vbs/repository"
	"go-vbs/usecase"
)

func main() {
	vrp := repository.NewVehicleRepository()
	gvuc := usecase.NewGetVehicleUseCase(vrp)

	vUUID, err := uuid.Parse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	if err != nil {
		fmt.Println(err)
	}

	vehicle, err := gvuc.ByUUID(vUUID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vehicle)
}
