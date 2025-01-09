package controller

import (
	"encoding/json"
	"log"
	"net/http"

	vehUCs "github.com/iondodon/go-vbs/usecase/vehicle"

	uuidLib "github.com/google/uuid"
	"github.com/gorilla/mux"
)

type VehicleController interface {
	HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request)
}

type vehicleController struct {
	infoLog, errorLog *log.Logger
	getVehicleUseCase vehUCs.GetVehicle
}

func NewVehicleController(
	infoLog, errorLog *log.Logger,
	getVehicleUseCase vehUCs.GetVehicle,
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

	responseJSON, err := json.Marshal(vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
