package controller

import (
	uuidLib "github.com/google/uuid"
	"github.com/gorilla/mux"
	"go-vbs/usecase"
	"log"
	"net/http"
)

type VehicleController interface {
	HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request)
}

type vehicleController struct {
	infoLog, errorLog *log.Logger
	getVehicleUseCase usecase.GetVehicle
}

func NewVehicleController(
	infoLog, errorLog *log.Logger,
	getVehicleUseCase usecase.GetVehicle,
) VehicleController {
	return &vehicleController{
		infoLog:           infoLog,
		errorLog:          errorLog,
		getVehicleUseCase: getVehicleUseCase,
	}
}

func (vc *vehicleController) HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	vUUID, err := uuidLib.Parse(uuid)
	if err != nil {
		vc.infoLog.Println(err)
	}

	vehicle, err := vc.getVehicleUseCase.ByUUID(vUUID)
	if err != nil {
		vc.infoLog.Println(err)
	}

	vc.infoLog.Printf("%+v\n", vehicle)
	vc.infoLog.Printf("%+v\n", vehicle.Bookings[0])
	vc.infoLog.Printf("%+v\n", vehicle.Bookings[0].Customer)
}
