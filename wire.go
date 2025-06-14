//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"net/http"

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
	"github.com/iondodon/go-vbs/server"
	"github.com/iondodon/go-vbs/vehicle/controller/vehicleController"
	"github.com/iondodon/go-vbs/vehicle/repository/vehicleRepository"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness/availabilityService"
	"github.com/iondodon/go-vbs/vehicle/vehicleBusiness/getVehicleService"
)

// Application struct to hold all application dependencies
type application struct {
	server      *http.Server
	database    *sql.DB
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
	getVehicleUC vehicleBusiness.GetVehicleUseCase,
) *vehicleController.Controller {
	return vehicleController.New(getVehicleUC)
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
	vehicleRepo vehicleBusiness.VehicleRepository,
	customerRepo customerBusiness.CustomerRepository,
	bookingRepo bookingBusiness.BookingRepository,
	bookingDateRepo bookingBusiness.BookingDateRepository,
	vehicleAvailabilityService vehicleBusiness.AvailabilityUseCase,
) bookingBusiness.BookVehicleUseCase {
	return bookVehicleService.New(
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
	db *sql.DB,
	bookVehicleUC bookingBusiness.BookVehicleUseCase,
	getAllBookingsUC bookingBusiness.GetAllBookingsUseCase,
) *bookingController.Controller {
	return bookingController.New(db, bookVehicleUC, getAllBookingsUC)
}

// Auth domain providers
func ProvideAuthController() *authController.Controller {
	return authController.New()
}

func ProvideServer(
	authCtrl *authController.Controller,
	vehicleCtrl *vehicleController.Controller,
	bookingCtrl *bookingController.Controller,
) *http.Server {
	return server.NewServer(authCtrl, vehicleCtrl, bookingCtrl)
}

// Application provider
func ProvideApplication(server *http.Server, db *sql.DB) *application {
	return &application{
		server:   server,
		database: db,
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

// Main provider set that combines all others
var ApplicationSet = wire.NewSet(
	DatabaseSet,
	RepositorySet,
	VehicleSet,
	BookingSet,
	AuthSet,
	ProvideServer,
	ProvideApplication,
)

// Wire injector function for application
func InitializeApplication() (*application, error) {
	wire.Build(ApplicationSet)
	return &application{}, nil
}
