package controller

import (
	"encoding/json"
	"go-vbs/usecase"
	"log"
	"net/http"

	uuidLib "github.com/google/uuid"
	"github.com/gorilla/mux"
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

	if err := json.NewEncoder(w).Encode(vehicle); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
