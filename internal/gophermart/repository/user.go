package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, userID int) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
