package handlers

import (
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("login")
	password := r.FormValue("password")

	user, err := u.UserUseCase.RegisterUser(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserExists) {
			http.Error(w, "Username is alredy taken", http.StatusConflict)
			return
		}
		http.Error(w, "Registration failes", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.ID),
	})

	w.WriteHeader(http.StatusOK)
}
