package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type PgStorage struct {
	db             *sql.DB
	userRepo       repository.UserRepository
	orderRepo      repository.OrderRepository
	balanceRepo    repository.BalanceRepository
	withdrawalRepo repository.WithdrawalRepository
}

func InitDB(_ context.Context, dsn, migrations string) (*PgStorage, error) {
	db, err := sql.Open("pgx", dsn)
	userRepo := NewUserRepo(db)
	orderRepo := NewOrderRepo(db)
	balanceRepo := NewBalanceRepo(db)
	withdrawalRepo := NewWithdrawalRepo(db)
	storage := &PgStorage{
		db:             db,
		userRepo:       userRepo,
		orderRepo:      orderRepo,
		balanceRepo:    balanceRepo,
		withdrawalRepo: withdrawalRepo,
	}

	err = storage.runMigrations(dsn, migrations)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *PgStorage) runMigrations(dsn, migrations string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrations), dsn)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
