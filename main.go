package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	// Initialize all dependencies using Wire
	app, err := InitializeApplication()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := app.database.Close(); err != nil {
			slog.Error("Error closing database", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		slog.Info("Server started")
		if err := app.server.ListenAndServe(); err != nil {
			slog.Error("Error starting server", "error", err)
			os.Exit(1)
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
	if err := app.server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", "error", err)
		os.Exit(1)
	}

	slog.Info("Shutting down...")

	os.Exit(0)
}
