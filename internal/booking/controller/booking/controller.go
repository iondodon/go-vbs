package booking

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	bookingController "github.com/iondodon/go-vbs/internal/booking/controller"
	bookingServices "github.com/iondodon/go-vbs/internal/booking/services"
)

type Controller struct {
	db                 *sql.DB
	bookVehicleUseCase bookingServices.BookVehicleUseCase
	getAllBookings     bookingServices.GetAllBookingsUseCase
}

func New(db *sql.DB, bookVehicleUseCase bookingServices.BookVehicleUseCase, getAllBookings bookingServices.GetAllBookingsUseCase) *Controller {
	return &Controller{
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
