package repository

import "github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"

type UserRepository interface {
	FindByLogin(login string) (*entity.User, error)
	Save(user *entity.User) error
}
