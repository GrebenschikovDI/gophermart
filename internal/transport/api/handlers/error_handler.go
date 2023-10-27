package handlers

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ErrorHandler struct {
	Log *logrus.Logger
}

func NewErrorHandler(log *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		Log: log,
	}
}

func (eh *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	eh.Log.WithFields(logrus.Fields{
		"error":  err,
		"method": r.Method,
		"path":   r.URL.Path,
	}).Error("Error in handler")

	switch {
	case errors.Is(err, gophermart.ErrUserExists):
		http.Error(w, "Username is already taken", http.StatusConflict)
	case errors.Is(err, gophermart.ErrUserNotFound):
		http.Error(w, "User not found", http.StatusNotFound)
	case errors.Is(err, gophermart.ErrUnauthorized):
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	case errors.Is(err, gophermart.ErrLowBalance):
		http.Error(w, "Balance is low", http.StatusPaymentRequired)
	case errors.Is(err, gophermart.ErrAlreadyExists):
		http.Error(w, "Order already exists", http.StatusConflict)
	case errors.Is(err, gophermart.ErrAlreadyTaken):
		http.Error(w, "Order is taken by another user", http.StatusConflict)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
