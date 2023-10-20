package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/delivery/api"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/persistence"
	"net/http"
)

const migrations = "migrations"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error with config", err)
	}

	db, err := persistence.InitDB(context.Background(), cfg.Dsn, cfg.Migrations)
	if err != nil {
		fmt.Println("Error with db", err)
	}

	_ = api.InitDependencies(db)

	userUseCase := usecase.NewUserUseCase(db.UserRepo)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: api.Router(*userUseCase),
	}

	fmt.Println("Running server at", cfg.RunAddress)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error with server", err)
	}
}
