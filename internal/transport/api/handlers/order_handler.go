package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	mw "github.com/GrebenschikovDI/gophermart.git/internal/transport/api/middleware"
	"io"
	"net/http"
	"strings"
	"time"
)

type OrderHandler struct {
	OrderUseCase usecase.OrderUseCase
}

type OrderRequest struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrderHandler(orderUseCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: orderUseCase,
	}
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}

	orders, err := o.OrderUseCase.GetByUserID(r.Context(), currentUserID)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get orders list", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		mw.LogError(w, r, gophermart.ErrNoOrders)
		http.Error(w, "Order list is empty", http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)

	ordersRequest := make([]*OrderRequest, 0)
	for _, order := range orders {
		ord := &OrderRequest{
			Number:     order.ID,
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		}
		if order.Accrual != nil {
			ord.Accrual = *order.Accrual
		}
		ordersRequest = append(ordersRequest, ord)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(ordersRequest); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (o *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}

	orderID := strings.TrimSpace(string(body))

	if !isLuhnValid(orderID) {
		mw.LogError(w, r, gophermart.ErrOrderBadFormat)
		http.Error(w, "Bad format of order", http.StatusUnprocessableEntity)
		return
	}

	currentUserID, err := getCurrentUser(r)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Cant get user id", http.StatusUnauthorized)
		return
	}
	_, err = o.OrderUseCase.CreateOrder(r.Context(), orderID, currentUserID, "NEW")
	if errors.Is(err, gophermart.ErrAlreadyTaken) {
		mw.LogError(w, r, err)
		http.Error(w, "order is taken by another user", http.StatusConflict)
		return
	} else if errors.Is(err, gophermart.ErrAlreadyExists) {
		mw.LogError(w, r, err)
		http.Error(w, "order already exists", http.StatusOK)
		return
	} else if err != nil {
		mw.LogError(w, r, err)
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
