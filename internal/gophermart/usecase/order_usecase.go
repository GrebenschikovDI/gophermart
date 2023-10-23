package usecase

import (
	"context"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

var AlreadyExists = errors.New("order already exists")
var AlreadyTaken = errors.New("order is taken by another user")

type OrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, id, userID int, status string) (*entity.Order, error) {
	newOrder := &entity.Order{
		ID:     id,
		UserID: userID,
		Status: status,
	}

	existingOrder, _ := u.orderRepo.GetByID(ctx, id)
	if existingOrder != nil {
		if existingOrder.UserID != newOrder.UserID {
			return existingOrder, AlreadyTaken
		} else {
			return existingOrder, AlreadyExists
		}
	} else {
		if err := u.orderRepo.Create(ctx, newOrder); err != nil {
			return nil, err
		}
		order, err := u.orderRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
}

func (u *OrderUseCase) GetOrderByID(ctx context.Context, id int) (*entity.Order, error) {
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

func (u *OrderUseCase) UpdateOrderStatus(ctx context.Context, id int, status string) (*entity.Order, error) {
	order, err := u.orderRepo.Update(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return order, nil
}
