package main

import (
	"context"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/accrual"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/logger"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/persistence"
	"github.com/GrebenschikovDI/gophermart.git/internal/transport/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.Initialize("info")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithField("error", err).Error("loading config failed")
	}
	log.Info("config loaded successfully")

	db, err := persistence.InitDB(context.Background(), cfg.Dsn, cfg.Migrations)
	if err != nil {
		log.WithField("error", err).Error("init DB failed")
	}
	log.Info("connected to DB")

	userUseCase := usecase.NewUserUseCase(db.UserRepo)
	orderUseCase := usecase.NewOrderUseCase(db.OrderRepo)
	balanceUseCase := usecase.NewBalanceUseCase(db.BalanceRepo)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(db.WithdrawalRepo)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: api.Router(*userUseCase, *orderUseCase, *balanceUseCase, *withdrawalUseCase, log),
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.WithField("error", err).Error("HTTP server shutdown")
		}
		close(stopped)
	}()

	sender := accrual.NewAccrual(log)
	go sender.Sender(context.Background(), *orderUseCase, *balanceUseCase, *cfg, 0, 0, 100)
	
	log.Info("sender activated")

	log.Infof("server running at %s", cfg.RunAddress)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.WithField("error", err).Fatal("Could not start server")
	}

	<-stopped
}
