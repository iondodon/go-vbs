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

	vehicleRepository := vehRepo.NewVehicleRepository(queries)
	customerRepository := custRepo.NewCustomerRepository(queries)
	bookingRepository := bookingRepoPkg.NewBookingRepository(queries)
	bookingDateRepository := bdRepoPkg.NewBookingDateRepository(queries)

	getVehicle := vehicleUCPKG.NewGetVehicle(vehicleRepository)
	isAvaiForHire := vehicleUCPKG.NewIsAvailableForHire(vehicleRepository)
	getBookingDates := bookingDateUCPkg.NewGetBookingDates(bookingDateRepository)
	bookVehicle := bookVehUCPkg.NewBookVehicle(infoLog, errorLog, vehicleRepository, customerRepository, bookingRepository, isAvaiForHire, getBookingDates)
	getAllBookins := bookVehUCPkg.NewGetAllBookings(bookingRepository)

	tokenController := controller.NewTokenController(infoLog, errorLog)
	vehicleController := controller.NewVehicleController(infoLog, errorLog, getVehicle)
	bookingController := controller.NewBookingController(infoLog, errorLog, db, bookVehicle, getAllBookins)

	router := http.NewServeMux()
	router.Handle("GET /login", controller.Handler(tokenController.Login))
	router.Handle("GET /refresh", controller.Handler(tokenController.Refresh))
	router.Handle("GET /vehicles/{uuid}", controller.Handler(vehicleController.HandleGetVehicleByUUID))
	router.Handle("POST /bookings", controller.Handler(bookingController.HandleBookVehicle))
	router.Handle("GET /bookings", middleware.JWT(controller.Handler(bookingController.HandleGetAllBookings)))

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
