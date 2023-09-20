package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/iondodon/go-vbs/dto"
	bookingUCPkg "github.com/iondodon/go-vbs/usecase/booking"
)

type BookingController interface {
	HandleBookVehicle(w http.ResponseWriter, r *http.Request)
}

type bookingController struct {
	infoLog, errorLog  *log.Logger
	bookVehicleUseCase bookingUCPkg.BookVehicleUseCase
}

func NewBookingController(
	infoLog, errorLog *log.Logger,
	bookVehicleUseCase bookingUCPkg.BookVehicleUseCase,
) BookingController {
	return &bookingController{
		infoLog:            infoLog,
		errorLog:           errorLog,
		bookVehicleUseCase: bookVehicleUseCase,
	}
}

func (bc *bookingController) HandleBookVehicle(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cbr dto.CreateBookingRequestDTO
	err = json.Unmarshal(reqBody, &cbr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bc.bookVehicleUseCase.ForPeriod(cbr.CustomerUUID, cbr.VehicleUUID, cbr.DatePeriodD)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
