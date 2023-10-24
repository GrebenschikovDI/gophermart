package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("user already exists")
var ErrUnauthorized = errors.New("authentication failed")

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) RegisterUser(ctx context.Context, username, password string) (*entity.User, error) {
	existingUser, err := u.GetUserByUsername(ctx, username)

	if existingUser != nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Login:    username,
		Password: passwordHash,
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u *UserUseCase) AuthenticateUser(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := u.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}
	if err := comparePasswords(user.Password, password); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	existingUser, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}
	return existingUser, nil
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func comparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}
