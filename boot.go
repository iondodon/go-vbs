package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/iondodon/go-vbs/controller/bookingController"
	"github.com/iondodon/go-vbs/controller/tokenController"
	"github.com/iondodon/go-vbs/controller/vehicleController"
	"github.com/iondodon/go-vbs/repository"
	"github.com/iondodon/go-vbs/repository/bookingDateRepository"
	"github.com/iondodon/go-vbs/repository/bookingRepository"
	"github.com/iondodon/go-vbs/repository/customerRepository"
	"github.com/iondodon/go-vbs/repository/vehicleRepository"
	"github.com/iondodon/go-vbs/usecase"
	"github.com/iondodon/go-vbs/usecase/booking/bookVehicle"
	"github.com/iondodon/go-vbs/usecase/booking/getAllBookings"
	"github.com/iondodon/go-vbs/usecase/bookingdate/getBookingDate"
	"github.com/iondodon/go-vbs/usecase/vehicle/getVehicle"
	"github.com/iondodon/go-vbs/usecase/vehicle/isVehicleAvailable"
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

	var vehicleRepo usecase.VehicleRepository = vehicleRepository.New(queries)
	var customerRepo usecase.CustomerRepository = customerRepository.New(queries)
	var bookingRepo usecase.BookingRepository = bookingRepository.New(queries)
	var bookingDateRepo usecase.BookingDateRepository = bookingDateRepository.New(queries)

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
