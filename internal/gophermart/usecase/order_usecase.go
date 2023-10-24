package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

var ErrAlreadyExists = errors.New("order already exists")
var ErrAlreadyTaken = errors.New("order is taken by another user")

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
			return nil, err
		}
		order, err := u.orderRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
	if existingOrder != nil {
		if existingOrder.UserID != newOrder.UserID {
			return existingOrder, ErrAlreadyTaken
		} else {
			return existingOrder, ErrAlreadyExists
		}
	} else {
		return nil, err
	}
}

func (u *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *OrderUseCase) GetByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {
	orderList, err := u.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orderList, nil
}

func (u *OrderUseCase) UpdateOrderStatus(ctx context.Context, id, status string, accrual *float64) (*entity.Order, error) {
	order, err := u.orderRepo.Update(ctx, id, status, accrual)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *OrderUseCase) GetToSend(ctx context.Context) ([]string, error) {
	orders, err := u.orderRepo.GetToSend(ctx)
	if err != nil {
		return nil, err
	}
	var numbers []string
	for _, value := range orders {
		numbers = append(numbers, value.ID)
	}
	return numbers, nil
}
