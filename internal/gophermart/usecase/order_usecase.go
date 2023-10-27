package usecase

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"github.com/pkg/errors"
)

type OrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, id string, userID int, status string) (*entity.Order, error) {
	newOrder := &entity.Order{
		ID:     id,
		UserID: userID,
		Status: status,
	}

	existingOrder, err := u.orderRepo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		if err := u.orderRepo.Create(ctx, newOrder); err != nil {
			return nil, errors.Wrapf(err, "creating order err, orderID: %s", id)
		}
		order, err := u.orderRepo.GetByID(ctx, id)
		if err != nil {
			return nil, errors.Wrapf(err, "creating order, can't get order: %s", id)
		}
		return order, nil
	}
	if existingOrder != nil {
		if existingOrder.UserID != newOrder.UserID {
			return existingOrder, errors.Wrapf(gophermart.ErrAlreadyTaken, "creating order, order %s is taken", id)
		} else {
			return existingOrder, errors.Wrapf(gophermart.ErrAlreadyExists, "creating order, order %s already exist", id)
		}
	} else {
		return nil, errors.Wrapf(err, "creating order, DB err, orderID: %s", id)
	}
}

func (u *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "get order by id err, id: %s", id)
	}
	return order, nil
}

func (u *OrderUseCase) GetByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {
	orderList, err := u.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "get order by userID err, id: %d", userID)
	}
	return orderList, nil
}

func (u *OrderUseCase) UpdateOrderStatus(ctx context.Context, id, status string, accrual *float64) (*entity.Order, error) {
	order, err := u.orderRepo.Update(ctx, id, status, accrual)
	if err != nil {
		return nil, errors.Wrapf(err, "update order err, id: %s", id)
	}
	return order, nil
}

func (u *OrderUseCase) GetToSend(ctx context.Context, offset, limit int) ([]string, error) {
	orders, err := u.orderRepo.GetToSend(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting orders to sent err")
	}
	var numbers []string
	for _, value := range orders {
		numbers = append(numbers, value.ID)
	}
	return numbers, nil
}
