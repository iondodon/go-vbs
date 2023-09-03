package main

import (
	"context"
	"github.com/gorilla/mux"
	"go-vbs/controller"
	"go-vbs/repository"
	"go-vbs/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	vrp := repository.NewVehicleRepository()
	gvuc := usecase.NewGetVehicleUseCase(vrp)
	vc := controller.NewVehicleController(gvuc)

	r := mux.NewRouter()

	r.HandleFunc("/vehicles/{uuid}", vc.HandleGetVehicleByUUID)

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Print(err)
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
		log.Println(err)
	}

	log.Println("Shutting down...")

	os.Exit(0)
}
