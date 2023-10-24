package api

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/delivery/api/handlers"
	mw "github.com/GrebenschikovDI/gophermart.git/internal/delivery/api/middleware"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(
	userUseCase usecase.UserUseCase,
	orderUseCase usecase.OrderUseCase,
	balanceUseCase usecase.BalanceUseCase,
	withdrawalUseCase usecase.WithdrawalUseCase,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		orderHandler := handlers.NewOrderHandler(orderUseCase)
		r.Get("/api/user/orders", orderHandler.GetOrders)
		r.Post("/api/user/orders", orderHandler.UploadOrder)
		balanceHandler := handlers.NewBalanceHandler(balanceUseCase, withdrawalUseCase)
		r.Get("/api/user/balance", balanceHandler.GetBalance)
		r.Post("/api/user/balance/withdraw", balanceHandler.Withdraw)

		r.Get("/api/user/withdrawals", balanceHandler.GetWithdrawals)
	})

	userHandler := handlers.NewUserHandler(userUseCase)
	r.Post("/api/user/register", userHandler.RegisterUser)
	r.Post("/api/user/login", userHandler.LoginUser)
	return r
}
