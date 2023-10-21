package handlers

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/dgrijalva/jwt-go"
)

// TODO: hide
var jwtSecret = []byte("super-secret-key")

func createAuthToken(user *entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Login

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
