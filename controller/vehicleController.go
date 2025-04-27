package controller

import (
	"encoding/json"
	"log"
	"net/http"

	vehUCs "github.com/iondodon/go-vbs/usecase/vehicle"

	uuidLib "github.com/google/uuid"
)

//gobok:constructor
//ctxboot:component
type VehicleController struct {
	infoLog           *log.Logger       `ctxboot:"inject"`
	errorLog          *log.Logger       `ctxboot:"inject"`
	getVehicleUseCase vehUCs.GetVehicle `ctxboot:"inject"`
}

func (vc *VehicleController) HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) error {
	uuid := r.PathValue("uuid")

	vUUID, err := uuidLib.Parse(uuid)
	if err != nil {
		return err
	}

	vehicle, err := vc.getVehicleUseCase.ByUUID(r.Context(), vUUID)
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
