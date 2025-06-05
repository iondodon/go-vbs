package bookingController

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/iondodon/go-vbs/booking/bookingBusiness"
	"github.com/iondodon/go-vbs/booking/bookingController"
)

type Controller struct {
	infoLog            *log.Logger
	errorLog           *log.Logger
	db                 *sql.DB
	bookVehicleUseCase bookingBusiness.BookVehicleUseCase
	getAllBookings     bookingBusiness.GetAllBookingsUseCase
}

func New(infoLog *log.Logger, errorLog *log.Logger, db *sql.DB, bookVehicleUseCase bookingBusiness.BookVehicleUseCase, getAllBookings bookingBusiness.GetAllBookingsUseCase) *Controller {
	return &Controller{
		infoLog:            infoLog,
		errorLog:           errorLog,
		db:                 db,
		bookVehicleUseCase: bookVehicleUseCase,
		getAllBookings:     getAllBookings,
	}
}

func (c *Controller) HandleBookVehicle(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var cbr bookingController.CreateBookingRequestDTO
	if err = json.Unmarshal(reqBody, &cbr); err != nil {
		return err
	}

	tx, err := c.db.BeginTx(r.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	if err = c.bookVehicleUseCase.ForPeriod(r.Context(), tx, cbr.CustomerUUID, cbr.VehicleUUID, cbr.DatePeriodD); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c *Controller) HandleGetAllBookings(w http.ResponseWriter, r *http.Request) error {
	bookings, err := c.getAllBookings.Execute(r.Context())
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
