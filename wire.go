//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/wire"
	"github.com/iondodon/go-vbs/auth/controller/authController"
	"github.com/iondodon/go-vbs/booking/bookingBusiness"
	"github.com/iondodon/go-vbs/booking/bookingBusiness/bookVehicleService"
	"github.com/iondodon/go-vbs/booking/bookingBusiness/getAllBookingsService"
	"github.com/iondodon/go-vbs/booking/bookingController/bookingController"
	"github.com/iondodon/go-vbs/booking/repository/bookingDateRepository"
	"github.com/iondodon/go-vbs/booking/repository/bookingRepository"
	"github.com/iondodon/go-vbs/customer/customerBusiness"
	"github.com/iondodon/go-vbs/customer/repository/customerRepository"
	"github.com/iondodon/go-vbs/repository"
	"github.com/iondodon/go-vbs/vehicle/controller/vehicleController"
	"github.com/iondodon/go-vbs/vehicle/repository/vehicleRepository"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness/availabilityService"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness/getVehicleService"
)

// Controllers struct to hold all controllers
type Controllers struct {
	Auth    *authController.Controller
	Vehicle *vehicleController.Controller
	Booking *bookingController.Controller
}

// Application struct to hold all application dependencies
type Application struct {
	Controllers *Controllers
	Database    *sql.DB
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

// Database provider
func ProvideDatabase() (*sql.DB, error) {
	return repository.NewInMemDBConn()
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
func ProvideBookingRepository(queries *repository.Queries) bookingBusiness.BookingRepository {
	return bookingRepository.New(queries)
}

func ProvideBookingDateRepository(queries *repository.Queries) bookingBusiness.BookingDateRepository {
	return bookingDateRepository.New(queries)
}

func ProvideBookVehicleUseCase(
	infoLog InfoLogger,
	errorLog ErrorLogger,
	vehicleRepo vehicleBusiness.VehicleRepository,
	customerRepo customerBusiness.CustomerRepository,
	bookingRepo bookingBusiness.BookingRepository,
	bookingDateRepo bookingBusiness.BookingDateRepository,
	vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase,
) bookingBusiness.BookVehicleUseCase {
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

func ProvideGetAllBookingsUseCase(bookingRepo bookingBusiness.BookingRepository) bookingBusiness.GetAllBookingsUseCase {
	return getAllBookingsService.New(bookingRepo)
}

func ProvideBookingController(
	infoLog InfoLogger,
	errorLog ErrorLogger,
	db *sql.DB,
	bookVehicleUC bookingBusiness.BookVehicleUseCase,
	getAllBookingsUC bookingBusiness.GetAllBookingsUseCase,
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

// Application provider
func ProvideApplication(controllers *Controllers, db *sql.DB) *Application {
	return &Application{
		Controllers: controllers,
		Database:    db,
	}
}

// Provider sets
var DatabaseSet = wire.NewSet(
	ProvideDatabase,
	ProvideQueries,
)

var RepositorySet = wire.NewSet(
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
	DatabaseSet,
	RepositorySet,
	VehicleSet,
	BookingSet,
	AuthSet,
	ProvideControllers,
	ProvideApplication,
)

// Wire injector function for Application
func InitializeApplication() (*Application, error) {
	wire.Build(ApplicationSet)
	return &Application{}, nil
}
