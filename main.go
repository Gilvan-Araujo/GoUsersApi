package main

import (
	"goUsersApi/api"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start", "error", err)
		os.Exit(1)
	}
	slog.Info("All systems offline")
}

func run() error {
	handler := api.NewHandler()

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
		Addr:         ":8080",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
