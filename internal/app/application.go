package app

import (
	"database/sql"
	"net/http"

	authController "github.com/iondodon/go-vbs/internal/auth/controller"
	bookingBusiness "github.com/iondodon/go-vbs/internal/booking/business"
	getAllBookingsService "github.com/iondodon/go-vbs/internal/booking/business/bookings/all/get"
	bookVehicleService "github.com/iondodon/go-vbs/internal/booking/business/bookings/vehicle"
	bookingController "github.com/iondodon/go-vbs/internal/booking/controller/booking"
	bookingRepository "github.com/iondodon/go-vbs/internal/booking/repository/booking"
	bookingDateRepository "github.com/iondodon/go-vbs/internal/booking/repository/booking_date"
	customerBusiness "github.com/iondodon/go-vbs/internal/customer/business"
	customerRepository "github.com/iondodon/go-vbs/internal/customer/repository/customer"
	"github.com/iondodon/go-vbs/internal/http/server"
	"github.com/iondodon/go-vbs/internal/repository"
	vehicleBusiness "github.com/iondodon/go-vbs/internal/vehicle/business"
	availabilityService "github.com/iondodon/go-vbs/internal/vehicle/business/availability"
	getVehicleService "github.com/iondodon/go-vbs/internal/vehicle/business/vehicle/get"
	vehicleController "github.com/iondodon/go-vbs/internal/vehicle/controller/vehicle"
	vehicleRepository "github.com/iondodon/go-vbs/internal/vehicle/repository/vehicle"
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

	var vehicleRepo vehicleBusiness.VehicleRepository = vehicleRepository.New(queries)
	var customerRepo customerBusiness.CustomerRepository = customerRepository.New(queries)
	var bookingRepo bookingBusiness.BookingRepository = bookingRepository.New(queries)
	var bookingDateRepo bookingBusiness.BookingDateRepository = bookingDateRepository.New(queries)
	var vehicleAvailabilityService bookingBusiness.VehicleAvailabilityService = availabilityService.New(vehicleRepo)

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
