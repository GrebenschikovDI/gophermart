package gophermart

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrUnauthorized  = errors.New("authentication failed")
	ErrLowBalance    = errors.New("balance is low")
	ErrAlreadyExists = errors.New("order already exists")
	ErrAlreadyTaken  = errors.New("order is taken by another user")
)
