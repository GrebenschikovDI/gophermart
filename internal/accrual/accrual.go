package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Response struct {
	Order    string
	Response *http.Response
	Error    error
}

type Accrual struct {
	log    *logrus.Logger
	client *http.Client
}

func NewAccrual(log *logrus.Logger) *Accrual {
	return &Accrual{
		log:    log,
		client: &http.Client{},
	}
}

func (a *Accrual) SendOrder(order string, cfg config.ServerConfig, responseChan chan<- Response) {
	server := cfg.AccrualAddress
	url := fmt.Sprintf("%s/api/orders/%s", server, order)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		a.log.Errorf("Ошибка при создании запроса для заказа %s: %v\n", order, err)
		responseChan <- Response{Order: order, Error: err}
		return
	}
	do, err := a.client.Do(request)
	if err != nil {
		a.log.Errorf("Ошибка при выполнении запроса для заказа %s: %v\n", order, err)
		responseChan <- Response{Order: order, Error: err}
		return
	}

	defer do.Body.Close()
	if do.StatusCode == http.StatusOK {
		responseChan <- Response{Order: order, Response: do}
	} else {
		responseChan <- Response{Order: order, Error: fmt.Errorf("ошибка при выполнении запроса для заказа %s. HTTP статус: %s", order, do.Status)}
	}
}

func (a *Accrual) Sender(ctx context.Context,
	u usecase.OrderUseCase,
	b usecase.BalanceUseCase,
	cfg config.ServerConfig,
	rate time.Duration,
	offset, limit int) {
	responseChan := make(chan Response)

	for {
		orders := a.GetOrders(ctx, u, offset, limit)
		for _, order := range orders {
			go a.SendOrder(order, cfg, responseChan)
		}

		for range orders {
			response := <-responseChan
			if response.Error != nil {
				a.log.Errorf("Ошибка при запросе для заказа %s: %v\n", response.Order, response.Error)
			} else {
				a.ProcessResponse(ctx, response.Response, u, b)
			}
		}
		time.Sleep(rate)
	}
}

func (a *Accrual) GetOrders(ctx context.Context, u usecase.OrderUseCase, offset, limit int) []string {
	send, err := u.GetToSend(ctx, offset, limit)
	if err != nil {
		a.log.Error(err)
		return nil
	}
	return send
}

func (a *Accrual) ProcessResponse(ctx context.Context, response *http.Response, u usecase.OrderUseCase, b usecase.BalanceUseCase) {
	if response.StatusCode == http.StatusOK {
		var result struct {
			Order   string  `json:"order"`
			Status  string  `json:"status"`
			Accrual float64 `json:"accrual"`
		}

		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			a.log.Errorf("Ошибка при декодировании JSON ответа: %v\n", err)
			return
		}

		switch result.Status {
		case "PROCESSED":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "PROCESSED", &result.Accrual); err != nil {
				a.log.Errorf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
			if result.Accrual > 0 {
				o, err := u.GetOrderByID(ctx, result.Order)
				user := o.UserID
				if err != nil {
					a.log.Errorf("Ошибка при получении ID пользователя для заказа %s: %v\n", result.Order, err)
				} else {
					if err := b.Add(ctx, user, result.Accrual); err != nil {
						a.log.Errorf("Ошибка при обновлении баланса пользователя %d: %v\n", user, err)
					}
				}
			}
		case "INVALID":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "INVALID", &result.Accrual); err != nil {
				a.log.Errorf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		case "PROCESSING":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "PROCESSING", &result.Accrual); err != nil {
				a.log.Errorf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		case "REGISTERED":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "REGISTERED", &result.Accrual); err != nil {
				a.log.Errorf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		default:
			a.log.Errorf("Заказ %s не был обработан: статус %s\n", result.Order, result.Status)
		}
	} else {
		a.log.Errorf("Ошибка при выполнении запроса. HTTP статус: %s\n", response.Status)
	}
}
