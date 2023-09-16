package main

import (
	"context"
	"go-vbs/controller"
	"go-vbs/integration"
	"go-vbs/repository"
	"go-vbs/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := integration.NewInMemDBConn()
	defer func(db integration.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Fatal(err)
		}
	}(db)
	if err != nil {
		errorLog.Fatal(err)
	}
	vrp := repository.NewVehicleRepository(db)
	gvuc := usecase.NewGetVehicleUseCase(vrp)
	vc := controller.NewVehicleController(infoLog, errorLog, gvuc)

	r := mux.NewRouter()
	r.HandleFunc("/vehicles/{uuid}", vc.HandleGetVehicleByUUID)

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     errorLog,
	}

	go func() {
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
