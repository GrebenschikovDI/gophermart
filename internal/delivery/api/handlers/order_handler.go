package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type OrderHandler struct {
	OrderUseCase usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: orderUseCase,
	}
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	currentUserID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	orders, err := o.OrderUseCase.GetByUserID(r.Context(), currentUserID)
	if err != nil {
		http.Error(w, "Cant get orders list", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		http.Error(w, "Order list is empty", http.StatusNoContent)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(orders); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (o *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}

	orderIDStr := strings.TrimSpace(string(body))

	if !isLuhnValid(orderIDStr) {
		http.Error(w, "Bad format of order", http.StatusUnprocessableEntity)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusInternalServerError)
		return
	}

	currentUserID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}
	_, err = o.OrderUseCase.CreateOrder(r.Context(), orderID, currentUserID, "NEW")
	if errors.Is(err, usecase.AlreadyTaken) {
		http.Error(w, "order is taken by another user", http.StatusConflict)
		return
	} else if errors.Is(err, usecase.AlreadyExists) {
		http.Error(w, "order already exists", http.StatusOK)
		return
	} else if err != nil {
		http.Error(w, "Unknown mistake", http.StatusInternalServerError)
	} else {
		http.Error(w, "New order is taken", http.StatusAccepted)
		return
	}
}

func isLuhnValid(purportedCC string) bool {
	sum := 0
	nDigits := len(purportedCC)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(purportedCC[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
