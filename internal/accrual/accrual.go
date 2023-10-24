package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"net/http"
	"time"
)

type Response struct {
	Order    string
	Response *http.Response
	Error    error
}

func SendOrder(order string, cfg config.ServerConfig, responseChan chan<- Response) {
	client := &http.Client{Timeout: 10 * time.Second}
	server := cfg.AccrualAddress
	url := fmt.Sprintf("http://%s/api/orders/%s", server, order)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса для заказа %s: %v\n", order, err)
		responseChan <- Response{Order: order, Error: err}
		return
	}
	do, err := client.Do(request)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса для заказа %s: %v\n", order, err)
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

func Sender(ctx context.Context,
	u usecase.OrderUseCase,
	b usecase.BalanceUseCase,
	cfg config.ServerConfig,
	rate time.Duration) {
	responseChan := make(chan Response)

	for {
		orders := GetOrders(ctx, u)
		for _, order := range orders {
			SendOrder(order, cfg, responseChan)
		}

		for range orders {
			response := <-responseChan
			if response.Error != nil {
				fmt.Printf("Ошибка при запросе для заказа %s: %v\n", response.Order, response.Error)
			} else {
				ProcessResponse(ctx, response.Response, u, b)
			}
		}
		time.Sleep(rate)
	}
}

func GetOrders(ctx context.Context, u usecase.OrderUseCase) []string {
	send, err := u.GetToSend(ctx)
	if err != nil {
		return nil
	}
	return send
}

func ProcessResponse(ctx context.Context, response *http.Response, u usecase.OrderUseCase, b usecase.BalanceUseCase) {
	if response.StatusCode == http.StatusOK {
		var result struct {
			Order   string  `json:"order"`
			Status  string  `json:"status"`
			Accrual float64 `json:"accrual"`
		}

		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			fmt.Printf("Ошибка при декодировании JSON ответа: %v\n", err)
			return
		}

		switch result.Status {
		case "PROCESSED":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "PROCESSED", &result.Accrual); err != nil {
				fmt.Printf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
			if result.Accrual > 0 {
				o, err := u.GetOrderByID(ctx, result.Order)
				user := o.UserID
				if err != nil {
					fmt.Printf("Ошибка при получении ID пользователя для заказа %s: %v\n", result.Order, err)
				} else {
					if err := b.Add(ctx, user, result.Accrual); err != nil {
						fmt.Printf("Ошибка при обновлении баланса пользователя %d: %v\n", user, err)
					}
				}
			}
		case "INVALID":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "INVALID", &result.Accrual); err != nil {
				fmt.Printf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		case "PROCESSING":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "PROCESSING", &result.Accrual); err != nil {
				fmt.Printf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		case "REGISTERED":
			if _, err := u.UpdateOrderStatus(ctx, result.Order, "REGISTERED", &result.Accrual); err != nil {
				fmt.Printf("Ошибка при обновлении статуса заказа %s: %v\n", result.Order, err)
			}
		default:
			fmt.Printf("Заказ %s не был обработан: статус %s\n", result.Order, result.Status)
		}
	} else {
		fmt.Printf("Ошибка при выполнении запроса. HTTP статус: %s\n", response.Status)
	}
}
