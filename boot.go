package main

import (
	"database/sql"
	"log"
	"os"

	authController "github.com/iondodon/go-vbs/auth/in/authController"
	"github.com/iondodon/go-vbs/booking/business"
	"github.com/iondodon/go-vbs/booking/business/bookVehicleService"
	"github.com/iondodon/go-vbs/booking/business/getAllBookingsService"
	bookingController "github.com/iondodon/go-vbs/booking/in/bookingController"
	"github.com/iondodon/go-vbs/booking/out/bookingDateRepository"
	"github.com/iondodon/go-vbs/booking/out/bookingRepository"
	customerBusiness "github.com/iondodon/go-vbs/customer/business"
	"github.com/iondodon/go-vbs/customer/out/customerRepository"
	"github.com/iondodon/go-vbs/repository"
	vehicleBusiness "github.com/iondodon/go-vbs/vehicle/business"
	"github.com/iondodon/go-vbs/vehicle/business/availabilityService"
	"github.com/iondodon/go-vbs/vehicle/business/getVehicleService"
	vehicleController "github.com/iondodon/go-vbs/vehicle/in/vehicleController"
	"github.com/iondodon/go-vbs/vehicle/out/vehicleRepository"
)

type Dependencies struct {
	AuthController    *authController.Controller
	VehicleController *vehicleController.Controller
	BookingController *bookingController.Controller
}

func BootstrapApplication(db *sql.DB) *Dependencies {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create repository layer (out adapters)
	queries := repository.New(db)

	// Vehicle domain
	var vehicleRepo vehicleBusiness.Repository = vehicleRepository.New(queries)
	var getVehicleUC vehicleBusiness.GetVehicleUseCase = getVehicleService.New(vehicleRepo)
	var vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase = availabilityService.New(vehicleRepo)

	// Customer domain
	var customerRepo customerBusiness.Repository = customerRepository.New(queries)

	// Booking domain
	var bookingRepo business.BookingRepository = bookingRepository.New(queries)
	var bookingDateRepo business.BookingDateRepository = bookingDateRepository.New(queries)

	var bookVehicleUC business.BookVehicleUseCase = bookVehicleService.New(
		infoLog,
		errorLog,
		vehicleRepo,  // Cross-domain dependency (vehicle out implements booking business interface)
		customerRepo, // Cross-domain dependency (customer out implements booking business interface)
		bookingRepo,
		bookingDateRepo,
		vehicleAvailabilityService, // Cross-domain dependency
	)

	var getAllBookingsUC business.GetAllBookingsUseCase = getAllBookingsService.New(bookingRepo)

	// Create controller layer (in adapters)
	authCtrl := authController.New(infoLog, errorLog)
	vehicleCtrl := vehicleController.New(infoLog, errorLog, getVehicleUC)
	bookingCtrl := bookingController.New(infoLog, errorLog, db, bookVehicleUC, getAllBookingsUC)

	return &Dependencies{
		AuthController:    authCtrl,
		VehicleController: vehicleCtrl,
		BookingController: bookingCtrl,
	}
}
