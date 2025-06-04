package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iondodon/go-vbs/handler"
	"github.com/iondodon/go-vbs/middleware"
	"github.com/iondodon/go-vbs/repository"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := repository.NewInMemDBConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Fatal(err)
		}
	}(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize all controllers using Wire
	controllers, err := InitializeControllers(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("GET /login", handler.Handler(controllers.Auth.Login))
	router.Handle("GET /refresh", handler.Handler(controllers.Auth.Refresh))
	router.Handle("GET /vehicles/{uuid}", handler.Handler(controllers.Vehicle.HandleGetVehicleByUUID))
	router.Handle("POST /bookings", handler.Handler(controllers.Booking.HandleBookVehicle))
	router.Handle("GET /bookings", middleware.JWT(handler.Handler(controllers.Booking.HandleGetAllBookings)))

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
