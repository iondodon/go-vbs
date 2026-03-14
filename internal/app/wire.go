//go:build wireinject
// +build wireinject

package app

import (
	"database/sql"
	"net/http"

	"github.com/google/wire"
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

// Application struct to hold all application dependencies
func ProvideDatabase() (*sql.DB, error) {
	return repository.NewInMemDBConn()
}

// Repository providers
func ProvideQueries(db *sql.DB) *repository.Queries {
	return repository.New(db)
}

// Vehicle domain providers
func ProvideVehicleRepository(queries *repository.Queries) vehicleServices.VehicleRepository {
	return vehicleRepository.New(queries)
}

func ProvideGetVehicleUseCase(vehicleRepo vehicleServices.VehicleRepository) vehicleServices.GetVehicleUseCase {
	return getVehicleService.New(vehicleRepo)
}

func ProvideAvailabilityUseCase(vehicleRepo vehicleServices.VehicleRepository) vehicleServices.AvailabilityUseCase {
	return availabilityService.New(vehicleRepo)
}

func ProvideVehicleController(
	getVehicleUC vehicleServices.GetVehicleUseCase,
) *vehicleController.Controller {
	return vehicleController.New(getVehicleUC)
}

// Customer domain providers
func ProvideCustomerRepository(queries *repository.Queries) customerServices.CustomerRepository {
	return customerRepository.New(queries)
}

// Booking domain providers
func ProvideBookingRepository(queries *repository.Queries) bookingServices.BookingRepository {
	return bookingRepository.New(queries)
}

func ProvideBookingDateRepository(queries *repository.Queries) bookingServices.BookingDateRepository {
	return bookingDateRepository.New(queries)
}

func ProvideBookVehicleUseCase(
	vehicleRepo vehicleServices.VehicleRepository,
	customerRepo customerServices.CustomerRepository,
	bookingRepo bookingServices.BookingRepository,
	bookingDateRepo bookingServices.BookingDateRepository,
	vehicleAvailabilityService vehicleServices.AvailabilityUseCase,
) bookingServices.BookVehicleUseCase {
	return bookVehicleService.New(
		vehicleRepo,
		customerRepo,
		bookingRepo,
		bookingDateRepo,
		vehicleAvailabilityService,
	)
}

func ProvideGetAllBookingsUseCase(bookingRepo bookingServices.BookingRepository) bookingServices.GetAllBookingsUseCase {
	return getAllBookingsService.New(bookingRepo)
}

func ProvideBookingController(
	db *sql.DB,
	bookVehicleUC bookingServices.BookVehicleUseCase,
	getAllBookingsUC bookingServices.GetAllBookingsUseCase,
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
func ProvideApplication(server *http.Server, db *sql.DB) *Application {
	return &Application{
		Server:   server,
		Database: db,
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
func InitializeApplication() (*Application, error) {
	wire.Build(ApplicationSet)
	return &Application{}, nil
}
