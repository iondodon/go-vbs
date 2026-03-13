package main

import (
	"database/sql"
	"net/http"

	authController "github.com/iondodon/go-vbs/auth/controller"
	bookingBusiness "github.com/iondodon/go-vbs/booking/business"
	getAllBookingsService "github.com/iondodon/go-vbs/booking/business/bookings/all/get"
	bookVehicleService "github.com/iondodon/go-vbs/booking/business/bookings/vehicle"
	bookingController "github.com/iondodon/go-vbs/booking/controller/booking"
	bookingRepository "github.com/iondodon/go-vbs/booking/repository/booking"
	bookingDateRepository "github.com/iondodon/go-vbs/booking/repository/booking_date"
	customerBusiness "github.com/iondodon/go-vbs/customer/business"
	customerRepository "github.com/iondodon/go-vbs/customer/repository/customer"
	"github.com/iondodon/go-vbs/repository"
	"github.com/iondodon/go-vbs/server"
	vehicleBusiness "github.com/iondodon/go-vbs/vehicle/business"
	availabilityService "github.com/iondodon/go-vbs/vehicle/business/availability"
	getVehicleService "github.com/iondodon/go-vbs/vehicle/business/vehicle/get"
	vehicleController "github.com/iondodon/go-vbs/vehicle/controller/vehicle"
	vehicleRepository "github.com/iondodon/go-vbs/vehicle/repository/vehicle"
)

type application struct {
	server   *http.Server
	database *sql.DB
}

func InitializeApplication() (*application, error) {
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

	return &application{
		server:   server.NewServer(authCtrl, vehicleCtrl, bookingCtrl),
		database: db,
	}, nil
}
