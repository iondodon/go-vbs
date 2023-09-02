package controller

import (
	"fmt"
	uuidLib "github.com/google/uuid"
	"github.com/gorilla/mux"
	"go-vbs/usecase"
	"net/http"
)

type VehicleController interface {
	HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request)
}

type vehicleController struct {
	getVehicleUseCase usecase.GetVehicle
}

func NewVehicleController(getVehicleUseCase usecase.GetVehicle) VehicleController {
	return &vehicleController{
		getVehicleUseCase: getVehicleUseCase,
	}
}

func (vc *vehicleController) HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	vUUID, err := uuidLib.Parse(uuid)
	if err != nil {
		fmt.Println(err)
	}

	vehicle, err := vc.getVehicleUseCase.ByUUID(vUUID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", vehicle)
	fmt.Printf("%+v\n", vehicle.Bookings[0])
	fmt.Printf("%+v\n", vehicle.Bookings[0].Customer)
}
