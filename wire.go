//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/wire"
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

// Controllers struct to hold all controllers
type Controllers struct {
	Auth    *authController.Controller
	Vehicle *vehicleController.Controller
	Booking *bookingController.Controller
}

// Logger wrapper types to distinguish between different loggers
type InfoLogger struct {
	*log.Logger
}

type ErrorLogger struct {
	*log.Logger
}

// Logger providers
func ProvideInfoLogger() InfoLogger {
	return InfoLogger{log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)}
}

func ProvideErrorLogger() ErrorLogger {
	return ErrorLogger{log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)}
}

// Repository providers
func ProvideQueries(db *sql.DB) *repository.Queries {
	return repository.New(db)
}

// Vehicle domain providers
func ProvideVehicleRepository(queries *repository.Queries) vehicleBusiness.VehicleRepository {
	return vehicleRepository.New(queries)
}

func ProvideGetVehicleUseCase(vehicleRepo vehicleBusiness.VehicleRepository) vehicleBusiness.GetVehicleUseCase {
	return getVehicleService.New(vehicleRepo)
}

func ProvideAvailabilityUseCase(vehicleRepo vehicleBusiness.VehicleRepository) vehicleBusiness.AvailabilityUseCase {
	return availabilityService.New(vehicleRepo)
}

func ProvideVehicleController(
	infoLog InfoLogger,
	errorLog ErrorLogger,
	getVehicleUC vehicleBusiness.GetVehicleUseCase,
) *vehicleController.Controller {
	return vehicleController.New(infoLog.Logger, errorLog.Logger, getVehicleUC)
}

// Customer domain providers
func ProvideCustomerRepository(queries *repository.Queries) customerBusiness.CustomerRepository {
	return customerRepository.New(queries)
}

// Booking domain providers
func ProvideBookingRepository(queries *repository.Queries) business.BookingRepository {
	return bookingRepository.New(queries)
}

func ProvideBookingDateRepository(queries *repository.Queries) business.BookingDateRepository {
	return bookingDateRepository.New(queries)
}

func ProvideBookVehicleUseCase(
	infoLog InfoLogger,
	errorLog ErrorLogger,
	vehicleRepo vehicleBusiness.VehicleRepository,
	customerRepo customerBusiness.CustomerRepository,
	bookingRepo business.BookingRepository,
	bookingDateRepo business.BookingDateRepository,
	vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase,
) business.BookVehicleUseCase {
	return bookVehicleService.New(
		infoLog.Logger,
		errorLog.Logger,
		vehicleRepo,
		customerRepo,
		bookingRepo,
		bookingDateRepo,
		vehicleAvailabilityService,
	)
}

func ProvideGetAllBookingsUseCase(bookingRepo business.BookingRepository) business.GetAllBookingsUseCase {
	return getAllBookingsService.New(bookingRepo)
}

func ProvideBookingController(
	infoLog InfoLogger,
	errorLog ErrorLogger,
	db *sql.DB,
	bookVehicleUC business.BookVehicleUseCase,
	getAllBookingsUC business.GetAllBookingsUseCase,
) *bookingController.Controller {
	return bookingController.New(infoLog.Logger, errorLog.Logger, db, bookVehicleUC, getAllBookingsUC)
}

// Auth domain providers
func ProvideAuthController(infoLog InfoLogger, errorLog ErrorLogger) *authController.Controller {
	return authController.New(infoLog.Logger, errorLog.Logger)
}

// Controllers provider
func ProvideControllers(
	authCtrl *authController.Controller,
	vehicleCtrl *vehicleController.Controller,
	bookingCtrl *bookingController.Controller,
) *Controllers {
	return &Controllers{
		Auth:    authCtrl,
		Vehicle: vehicleCtrl,
		Booking: bookingCtrl,
	}
}

// Provider sets
var RepositorySet = wire.NewSet(
	ProvideQueries,
	ProvideVehicleRepository,
	ProvideCustomerRepository,
	ProvideBookingRepository,
	ProvideBookingDateRepository,
)

var VehicleSet = wire.NewSet(
	ProvideGetVehicleUseCase,
	ProvideAvailabilityUseCase,
	ProvideVehicleController,
)

var BookingSet = wire.NewSet(
	ProvideBookVehicleUseCase,
	ProvideGetAllBookingsUseCase,
	ProvideBookingController,
)

var AuthSet = wire.NewSet(
	ProvideAuthController,
)

var LoggerSet = wire.NewSet(
	ProvideInfoLogger,
	ProvideErrorLogger,
)

// Main provider set that combines all others
var ApplicationSet = wire.NewSet(
	LoggerSet,
	RepositorySet,
	VehicleSet,
	BookingSet,
	AuthSet,
	ProvideControllers,
)

// Wire injector function for Controllers
func InitializeControllers(db *sql.DB) (*Controllers, error) {
	wire.Build(ApplicationSet)
	return &Controllers{}, nil
}
