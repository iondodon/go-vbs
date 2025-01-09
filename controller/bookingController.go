package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/iondodon/go-vbs/dto"
	bookingUCs "github.com/iondodon/go-vbs/usecase/booking"
)

type BookingController interface {
	HandleBookVehicle(w http.ResponseWriter, r *http.Request)
	HandleGetAllBookings(w http.ResponseWriter, r *http.Request)
}

type bookingController struct {
	infoLog, errorLog  *log.Logger
	bookVehicleUseCase bookingUCs.BookVehicle
	getAllBookings     bookingUCs.GetAllBookings
}

func NewBookingController(
	infoLog, errorLog *log.Logger,
	bookVehicleUseCase bookingUCs.BookVehicle,
	getAllBookings bookingUCs.GetAllBookings,
) BookingController {
	return &bookingController{
		infoLog:            infoLog,
		errorLog:           errorLog,
		bookVehicleUseCase: bookVehicleUseCase,
		getAllBookings:     getAllBookings,
	}
}

func (c *bookingController) HandleBookVehicle(w http.ResponseWriter, r *http.Request) {
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

	err = c.bookVehicleUseCase.ForPeriod(cbr.CustomerUUID, cbr.VehicleUUID, cbr.DatePeriodD)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *bookingController) HandleGetAllBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := c.getAllBookings.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(bookings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
