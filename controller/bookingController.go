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
	HandleBookVehicle(w http.ResponseWriter, r *http.Request) error
	HandleGetAllBookings(w http.ResponseWriter, r *http.Request) error
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

func (c *bookingController) HandleBookVehicle(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var cbr dto.CreateBookingRequestDTO
	err = json.Unmarshal(reqBody, &cbr)
	if err != nil {
		return err
	}

	err = c.bookVehicleUseCase.ForPeriod(cbr.CustomerUUID, cbr.VehicleUUID, cbr.DatePeriodD)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c *bookingController) HandleGetAllBookings(w http.ResponseWriter, r *http.Request) error {
	bookings, err := c.getAllBookings.Execute()
	if err != nil {
		return err
	}

	jsonResp, err := json.Marshal(bookings)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	return nil
}
