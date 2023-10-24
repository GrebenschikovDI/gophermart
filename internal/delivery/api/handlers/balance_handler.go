package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"net/http"
	"time"
)

type BalanceHandler struct {
	BalanceUseCase    usecase.BalanceUseCase
	WithdrawalUseCase usecase.WithdrawalUseCase
}

type BalanceRequest struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalRequest struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_At,omitempty"`
}

func NewBalanceHandler(
	balanceUseCase usecase.BalanceUseCase,
	withdrawalUseCase usecase.WithdrawalUseCase) *BalanceHandler {
	return &BalanceHandler{
		BalanceUseCase:    balanceUseCase,
		WithdrawalUseCase: withdrawalUseCase,
	}
}

func (b *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}
	balance, err := b.BalanceUseCase.Get(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "Cant get balance", http.StatusInternalServerError)
		return
	}
	balanceToSend := BalanceRequest{
		Current:   balance.Amount,
		Withdrawn: balance.Withdrawn,
	}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(balanceToSend); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (b *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req WithdrawalRequest

	currentUserID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	order := req.Order
	sum := req.Sum

	if !isLuhnValid(order) {
		http.Error(w, "Order not valid", http.StatusUnprocessableEntity)
		return
	}

	err = b.BalanceUseCase.Withdraw(r.Context(), currentUserID, sum)
	if errors.Is(err, usecase.LowBalance) {
		http.Error(w, "Low balance", http.StatusPaymentRequired)
		return
	} else if err != nil {
		http.Error(w, "Error processing withdrawal", http.StatusInternalServerError)
	}

	newWithdrawal := &entity.Withdrawal{
		UserID:  currentUserID,
		OrderID: order,
		Amount:  sum,
	}

	err = b.WithdrawalUseCase.Add(r.Context(), newWithdrawal)

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	withdrawals, err := b.WithdrawalUseCase.GetList(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "can't get withdrawals", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		http.Error(w, "Order list is empty", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	answer := make([]*WithdrawalRequest, 0)
	for _, withdrawal := range withdrawals {
		wd := &WithdrawalRequest{
			Order:       withdrawal.OrderID,
			Sum:         withdrawal.Amount,
			ProcessedAt: withdrawal.ProcessedAt,
		}
		answer = append(answer, wd)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(answer); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
