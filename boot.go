package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/business/booking/bookVehicle"
	"github.com/iondodon/go-vbs/business/booking/getAllBookings"
	"github.com/iondodon/go-vbs/business/bookingdate/getBookingDate"
	"github.com/iondodon/go-vbs/business/vehicle/getVehicle"
	"github.com/iondodon/go-vbs/business/vehicle/isVehicleAvailable"
	"github.com/iondodon/go-vbs/controller/bookingController"
	"github.com/iondodon/go-vbs/controller/tokenController"
	"github.com/iondodon/go-vbs/controller/vehicleController"
	"github.com/iondodon/go-vbs/repository"
	"github.com/iondodon/go-vbs/repository/bookingDateRepository"
	"github.com/iondodon/go-vbs/repository/bookingRepository"
	"github.com/iondodon/go-vbs/repository/customerRepository"
	"github.com/iondodon/go-vbs/repository/vehicleRepository"
)

type Dependencies struct {
	TokenController   *tokenController.Controller
	VehicleController *vehicleController.Controller
	BookingController *bookingController.Controller
}

func BootstrapApplication(db *sql.DB) *Dependencies {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create repository layer
	queries := repository.New(db)

	var vehicleRepo business.VehicleRepository = vehicleRepository.New(queries)
	var customerRepo business.CustomerRepository = customerRepository.New(queries)
	var bookingRepo business.BookingRepository = bookingRepository.New(queries)
	var bookingDateRepo business.BookingDateRepository = bookingDateRepository.New(queries)

	// Create use case layer
	var getVehicleUC getVehicle.UseCase = getVehicle.New(vehicleRepo)
	isAvailableForHireUC := isVehicleAvailable.New(vehicleRepo)
	getBookingDatesUC := getBookingDate.New(bookingDateRepo)
	var getAllBookingsUC getAllBookings.UseCase = getAllBookings.New(bookingRepo)
	var bookVehicleUC bookVehicle.UseCase = bookVehicle.New(
		infoLog,
		errorLog,
		vehicleRepo,
		customerRepo,
		bookingRepo,
		isAvailableForHireUC,
		getBookingDatesUC,
	)

	// Create controller layer
	tokenCtrl := tokenController.New(infoLog, errorLog)
	vehicleCtrl := vehicleController.New(infoLog, errorLog, getVehicleUC)
	bookingCtrl := bookingController.New(infoLog, errorLog, db, bookVehicleUC, getAllBookingsUC)

	return &Dependencies{
		TokenController:   tokenCtrl,
		VehicleController: vehicleCtrl,
		BookingController: bookingCtrl,
	}
}
