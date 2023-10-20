package api

import (
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/persistence"
)

type Dependency struct {
	UserUseCase usecase.UserUseCase
}

func InitDependencies(storage *persistence.PgStorage) *Dependency {
	return &Dependency{
		UserUseCase: *usecase.NewUserUseCase(storage.UserRepo),
	}
}
