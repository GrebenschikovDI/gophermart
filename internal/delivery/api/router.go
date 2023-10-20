package api

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/delivery/api/handlers"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(userUseCase usecase.UserUseCase) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	userHandler := handlers.NewUserHandler(userUseCase)

	r.Get("/api/user/orders", handlers.GetOrders)
	r.Get("/api/user/balance", handlers.GetBalance)
	r.Get("/api/user/withdrawals", handlers.GetWithdrawals)
	r.Post("/api/user/register", userHandler.RegisterUser)
	r.Post("/api/user/login", handlers.LoginUser)
	r.Post("/api/user/orders", handlers.CreateOrders)
	r.Post("/api/user/balance/withdraw", handlers.WithdrawBalance)
	return r
}
