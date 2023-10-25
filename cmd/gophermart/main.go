package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/accrual"
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

	userUseCase := usecase.NewUserUseCase(db.UserRepo)
	orderUseCase := usecase.NewOrderUseCase(db.OrderRepo)
	balanceUseCase := usecase.NewBalanceUseCase(db.BalanceRepo)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(db.WithdrawalRepo)

	fmt.Println(cfg.AccrualAddress)
	fmt.Println(cfg.Dsn)
	fmt.Println(cfg.RunAddress)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: api.Router(*userUseCase, *orderUseCase, *balanceUseCase, *withdrawalUseCase),
	}

	go accrual.Sender(context.Background(), *orderUseCase, *balanceUseCase, *cfg, 0)

	fmt.Println("Running server at", cfg.RunAddress)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error with server", err)
	}
}
