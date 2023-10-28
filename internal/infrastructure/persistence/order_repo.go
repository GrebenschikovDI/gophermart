package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) repository.OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (o *orderRepo) Create(ctx context.Context, order *entity.Order) error {
	_, err := o.db.ExecContext(ctx, "INSERT INTO orders (id, user_id, status) values ($1, $2, $3)", order.ID,
		order.UserID, order.Status)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderRepo) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	row := o.db.QueryRowContext(
		ctx,
		"SELECT  id, user_id, status, accrual, uploaded_at FROM orders WHERE id = $1",
		id,
	)
	order := &entity.Order{}
	var accrual sql.NullFloat64
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &accrual, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	if accrual.Valid {
		acc := accrual.Float64
		order.Accrual = &acc
	} else {
		order.Accrual = nil
	}
	return order, nil
}

func (o *orderRepo) Update(ctx context.Context, id, status string, accrual *float64) (*entity.Order, error) {
	if accrual != nil {
		_, err := o.db.ExecContext(
			ctx,
			"UPDATE orders SET status = $1, accrual = $2 WHERE id = $3",
			status, *accrual, id)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := o.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, id)
		if err != nil {
			return nil, err
		}
	}
	return o.GetByID(ctx, id)
}

func (o *orderRepo) GetByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {
	rows, err := o.db.QueryContext(
		ctx,
		"SELECT id, user_id, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC ",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		var accrual sql.NullFloat64
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			acc := accrual.Float64
			order.Accrual = &acc
		} else {
			order.Accrual = nil
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepo) GetToSend(ctx context.Context, offset, limit int) ([]*entity.Order, error) {
	query := `
        SELECT id, user_id, status, accrual, uploaded_at
        FROM orders
        WHERE status IN ('NEW', 'PROCESSING', 'REGISTERED')
        LIMIT $1 OFFSET $2
    `

	rows, err := o.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		var accrual sql.NullFloat64
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			acc := accrual.Float64
			order.Accrual = &acc
		} else {
			order.Accrual = nil
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
