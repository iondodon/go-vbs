package vehicleController

import (
	"encoding/json"
	"log"
	"net/http"

	uuidLib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness"
)

type Controller struct {
	getVehicleUseCase vehicleBusiness.GetVehicleUseCase
}

func New(infoLog *log.Logger, errorLog *log.Logger, getVehicleUseCase vehicleBusiness.GetVehicleUseCase) *Controller {
	return &Controller{
		getVehicleUseCase: getVehicleUseCase,
	}
}

func (c *Controller) HandleGetVehicleByUUID(w http.ResponseWriter, r *http.Request) error {
	uuid := r.PathValue("uuid")

	vUUID, err := uuidLib.Parse(uuid)
	if err != nil {
		return err
	}

	vehicle, err := c.getVehicleUseCase.ByUUID(r.Context(), vUUID)
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
