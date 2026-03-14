package app

import (
	"database/sql"
	"net/http"

	authController "github.com/iondodon/go-vbs/internal/auth/controller"
	bookingController "github.com/iondodon/go-vbs/internal/booking/controller/booking"
	bookingRepository "github.com/iondodon/go-vbs/internal/booking/repository/booking"
	bookingDateRepository "github.com/iondodon/go-vbs/internal/booking/repository/booking_date"
	bookingServices "github.com/iondodon/go-vbs/internal/booking/services"
	getAllBookingsService "github.com/iondodon/go-vbs/internal/booking/services/bookings/all/get"
	bookVehicleService "github.com/iondodon/go-vbs/internal/booking/services/bookings/vehicle"
	customerRepository "github.com/iondodon/go-vbs/internal/customer/repository/customer"
	customerServices "github.com/iondodon/go-vbs/internal/customer/services"
	"github.com/iondodon/go-vbs/internal/http/server"
	"github.com/iondodon/go-vbs/internal/repository"
	vehicleController "github.com/iondodon/go-vbs/internal/vehicle/controller/vehicle"
	vehicleRepository "github.com/iondodon/go-vbs/internal/vehicle/repository/vehicle"
	vehicleServices "github.com/iondodon/go-vbs/internal/vehicle/services"
	availabilityService "github.com/iondodon/go-vbs/internal/vehicle/services/availability"
	getVehicleService "github.com/iondodon/go-vbs/internal/vehicle/services/vehicle/get"
)

type Application struct {
	Server   *http.Server
	Database *sql.DB
}

func InitializeApplication() (*Application, error) {
	db, err := repository.NewInMemDBConn()
	if err != nil {
		return nil, err
	}

	queries := repository.New(db)

	var vehicleRepo vehicleServices.VehicleRepository = vehicleRepository.New(queries)
	var customerRepo customerServices.CustomerRepository = customerRepository.New(queries)
	var bookingRepo bookingServices.BookingRepository = bookingRepository.New(queries)
	var bookingDateRepo bookingServices.BookingDateRepository = bookingDateRepository.New(queries)
	var vehicleAvailabilityService bookingServices.VehicleAvailabilityService = availabilityService.New(vehicleRepo)

	vehicleCtrl := vehicleController.New(getVehicleService.New(vehicleRepo))
	bookingCtrl := bookingController.New(
		db,
		bookVehicleService.New(vehicleRepo, customerRepo, bookingRepo, bookingDateRepo, vehicleAvailabilityService),
		getAllBookingsService.New(bookingRepo),
	)
	authCtrl := authController.New()

	return &Application{
		Server:   server.NewServer(authCtrl, vehicleCtrl, bookingCtrl),
		Database: db,
	}, nil
}
