package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, user *entity.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE id = $1", id)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1",
		username)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
