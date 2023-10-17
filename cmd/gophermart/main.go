package main

import (
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/delivery/api"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"net/http"
)

const migrations = "migrations"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error with config", err)
	}

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: api.Router(),
	}

	fmt.Println("Running server at", cfg.RunAddress)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error with server", err)
	}
}
