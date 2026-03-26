package vehicle

import (
	"encoding/json"
	"net/http"

	uuidLib "github.com/google/uuid"
	vehicleServices "github.com/iondodon/go-vbs/internal/vehicle/services"
)

type Controller struct {
	getVehicleUseCase vehicleServices.GetVehicleUseCase
}

func New(getVehicleUseCase vehicleServices.GetVehicleUseCase) *Controller {
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
