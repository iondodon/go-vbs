package vehicleController

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/iondodon/go-vbs/vehicle/business"

	uuidLib "github.com/google/uuid"
)

type Controller struct {
	infoLog           *log.Logger
	errorLog          *log.Logger
	getVehicleUseCase business.GetVehicleUseCase
}

func New(infoLog *log.Logger, errorLog *log.Logger, getVehicleUseCase business.GetVehicleUseCase) *Controller {
	return &Controller{
		infoLog:           infoLog,
		errorLog:          errorLog,
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
