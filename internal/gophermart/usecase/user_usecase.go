package usecase

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

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
		return nil, errors.Wrapf(gophermart.ErrUserExists, "register user err, username: %s exists", username)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "register user err, username: %s", username)
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, errors.Wrapf(err, "register user err, can't hash password")
	}

	newUser := &entity.User{
		Login:    username,
		Password: passwordHash,
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, errors.Wrapf(err, "register user err, can't create user %s", username)
	}

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "register user err, can't get user %s", username)
	}
	return user, nil
}

func (u *UserUseCase) AuthenticateUser(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := u.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "auth user err, can't get username: %s", username)
	}
	if user == nil {
		return nil, errors.Wrapf(gophermart.ErrUnauthorized, "auth user err, username: %s", username)
	}
	if err := comparePasswords(user.Password, password); err != nil {
		return nil, errors.Wrapf(err, "auth user err, probably wrong password, username: %s", username)
	}
	return user, nil
}

func (u *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	existingUser, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "get user by username err, username %s", username)
	}
	if existingUser == nil {
		return nil, errors.Wrapf(gophermart.ErrUserNotFound, "get user by username err, username %s", username)
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
		return gophermart.ErrUnauthorized
	}
	return nil
}
