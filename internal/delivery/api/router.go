package api

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/delivery/api/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/user/orders", handlers.GetOrders)
	r.Get("/api/user/balance", handlers.GetBalance)
	r.Get("/api/user/withdrawals", handlers.GetWithdrawals)
	r.Post("/api/user/register", handlers.RegisterUser)
	r.Post("/api/user/login", handlers.LoginUser)
	r.Post("/api/user/orders", handlers.CreateOrders)
	r.Post("/api/user/balance/withdraw", handlers.WithdrawBalance)
	return r
}
