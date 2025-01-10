package controller

import (
	"encoding/json"
	"log"
	"net/http"

	vehUCs "github.com/iondodon/go-vbs/usecase/vehicle"

	uuidLib "github.com/google/uuid"
)

type VehicleController interface {
	HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) error
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

func (vc *vehicleController) HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) error {
	uuid := r.PathValue("uuid")

	vUUID, err := uuidLib.Parse(uuid)
	if err != nil {
		return err
	}

	vehicle, err := vc.getVehicleUseCase.ByUUID(vUUID)
	if err != nil {
		return err
	}

	responseJSON, err := json.Marshal(vehicle)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	return nil
}
