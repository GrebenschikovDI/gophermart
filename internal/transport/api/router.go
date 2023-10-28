package api

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/transport/api/handlers"
	mw "github.com/GrebenschikovDI/gophermart.git/internal/transport/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func Router(
	userUseCase usecase.UserUseCase,
	orderUseCase usecase.OrderUseCase,
	balanceUseCase usecase.BalanceUseCase,
	withdrawalUseCase usecase.WithdrawalUseCase,
	log *logrus.Logger,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.LoggerMiddleware(log))
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
