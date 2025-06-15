package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iondodon/go-vbs/handler"
	"github.com/iondodon/go-vbs/middleware"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	// Initialize all dependencies using Wire
	app, err := InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := app.Database.Close(); err != nil {
			// Use panic since we don't have access to errorLog here
			panic(err)
		}
	}()

	router := http.NewServeMux()
	router.Handle("GET /login", handler.Handler(app.Controllers.Auth.Login))
	router.Handle("GET /refresh", handler.Handler(app.Controllers.Auth.Refresh))
	router.Handle("GET /vehicles/{uuid}", handler.Handler(app.Controllers.Vehicle.HandleGetVehicleByUUID))
	router.Handle("POST /bookings", handler.Handler(app.Controllers.Booking.HandleBookVehicle))
	router.Handle("GET /bookings", middleware.JWT(handler.Handler(app.Controllers.Booking.HandleGetAllBookings)))

	// Mount Swagger UI only in development mode
	if os.Getenv("GO_ENV") == "development" {
		slog.Info("Running in development mode - Swagger UI enabled")
		router.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./swagger-ui"))))
		router.Handle("/docs/openapi.yaml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/docs/openapi.yaml" {
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, "openapi.yaml")
		}))
		slog.Info("Swagger UI available at http://localhost:8000/docs")
	}

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		slog.Info("Server started")
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
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
		panic(err)
	}

	slog.Info("Shutting down...")

	os.Exit(0)
}
