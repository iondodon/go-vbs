package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/iondodon/go-vbs/auth/controller/authController"
	"github.com/iondodon/go-vbs/booking/bookingController/bookingController"
	"github.com/iondodon/go-vbs/handler"
	"github.com/iondodon/go-vbs/middleware"
	"github.com/iondodon/go-vbs/vehicle/controller/vehicleController"
)

func NewServer(
	authCtrl *authController.Controller,
	vehicleCtrl *vehicleController.Controller,
	bookingCtrl *bookingController.Controller,
) *http.Server {
	router := http.NewServeMux()
	router.Handle("GET /login", handler.Handler(authCtrl.Login))
	router.Handle("GET /refresh", handler.Handler(authCtrl.Refresh))
	router.Handle("GET /vehicles/{uuid}", handler.Handler(vehicleCtrl.HandleGetVehicleByUUID))
	router.Handle("POST /bookings", handler.Handler(bookingCtrl.HandleBookVehicle))
	router.Handle("GET /bookings", middleware.JWT(handler.Handler(bookingCtrl.HandleGetAllBookings)))

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

	return &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}