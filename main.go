package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iondodon/go-vbs/controller"
	"github.com/iondodon/go-vbs/integration"
	"github.com/iondodon/go-vbs/middleware"
	"github.com/iondodon/go-vbs/repository"
	bookingRepoPkg "github.com/iondodon/go-vbs/repository/booking"
	bdRepoPkg "github.com/iondodon/go-vbs/repository/bookingdate"
	custRepo "github.com/iondodon/go-vbs/repository/customer"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"
	bookVehUCPkg "github.com/iondodon/go-vbs/usecase/booking"
	bookingDateUCPkg "github.com/iondodon/go-vbs/usecase/bookingdate"
	vehicleUCPKG "github.com/iondodon/go-vbs/usecase/vehicle"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := integration.NewInMemDBConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Fatal(err)
		}
	}(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	queries := repository.New(db)

	vehicleRepository := vehRepo.NewVehicleRepositoryBuilder().
		Queries(queries).
		Build()
	customerRepository := custRepo.NewCustomerRepositoryBuilder().
		Queries(queries).
		Build()
	bookingRepository := bookingRepoPkg.NewBookingRepositoryBuilder().
		Queries(queries).
		Build()
	bookingDateRepository := bdRepoPkg.NewBookingDateRepositoryBuilder().
		Queries(queries).
		Build()

	getVehicle := vehicleUCPKG.NewGetVehicleBuilder().
		VehicleRepository(*vehicleRepository).
		Build()
	isAvaiForHire := vehicleUCPKG.NewIsAvailableForHireBuilder().
		VehRepo(*vehicleRepository).
		Build()
	getBookingDates := bookingDateUCPkg.NewGetBookingDatesBuilder().
		BdRepo(*bookingDateRepository).
		Build()
	bookVehicle := bookVehUCPkg.NewBookVehicleBuilder().
		InfoLog(infoLog).
		ErrorLog(errorLog).
		VehRepo(*vehicleRepository).
		CustRepo(*customerRepository).
		BookingRepo(*bookingRepository).
		IsAvailableForHire(*isAvaiForHire).
		GetBookingDates(*getBookingDates).
		Build()
	getAllBookins := bookVehUCPkg.NewGetAllBookingsBuilder().
		BookingRepo(*bookingRepository).
		Build()

	tokenController := controller.NewTokenControllerBuilder().
		InfoLog(infoLog).
		ErrorLog(errorLog).
		Build()
	vehicleController := controller.NewVehicleControllerBuilder().
		InfoLog(infoLog).
		ErrorLog(errorLog).
		GetVehicleUseCase(*getVehicle).
		Build()
	bookingController := controller.NewBookingControllerBuilder().
		InfoLog(infoLog).
		ErrorLog(errorLog).
		Db(db).
		BookVehicleUseCase(*bookVehicle).
		GetAllBookings(*getAllBookins).
		Build()

	router := http.NewServeMux()
	router.Handle("GET /login", controller.Handler(tokenController.Login))
	router.Handle("GET /refresh", controller.Handler(tokenController.Refresh))
	router.Handle("GET /vehicles/{uuid}", controller.Handler(vehicleController.HandleGetVehicleByUUID))
	router.Handle("POST /bookings", controller.Handler(bookingController.HandleBookVehicle))
	router.Handle("GET /bookings", middleware.JWT(controller.Handler(bookingController.HandleGetAllBookings)))

	// Mount Swagger UI only in development mode
	if os.Getenv("GO_ENV") == "development" {
		infoLog.Println("Running in development mode - Swagger UI enabled")
		router.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./swagger-ui"))))
		router.Handle("/docs/openapi.yaml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/docs/openapi.yaml" {
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, "openapi.yaml")
		}))
		infoLog.Println("Swagger UI available at http://localhost:8000/docs")
	}

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     errorLog,
	}

	go func() {
		infoLog.Println("Server started")
		if err := srv.ListenAndServe(); err != nil {
			errorLog.Print(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	if err := srv.Shutdown(ctx); err != nil {
		errorLog.Println(err)
	}

	infoLog.Println("Shutting down...")

	os.Exit(0)
}
