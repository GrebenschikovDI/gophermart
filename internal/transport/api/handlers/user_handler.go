package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"net/http"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req Auth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	username := req.Login
	password := req.Password

	if username == "" || password == "" {
		http.Error(w, "Username and password must not be empty", http.StatusBadRequest)
		return
	}

	user, err := u.UserUseCase.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	token, err := createAuthToken(user)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req Auth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	username := req.Login
	password := req.Password

	if username == "" || password == "" {
		http.Error(w, "Username and password must not be empty", http.StatusBadRequest)
		return
	}

	user, err := u.UserUseCase.RegisterUser(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, gophermart.ErrUserExists) {
			http.Error(w, "Username is already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	token, err := createAuthToken(user)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}
