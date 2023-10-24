package accrual

import (
	"context"
	"fmt"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/usecase"
	"github.com/GrebenschikovDI/gophermart.git/internal/infrastructure/config"
	"net/http"
	"time"
)

func SendOrder(order string, cfg config.ServerConfig) {
	client := &http.Client{Timeout: 10 * time.Second}
	server := cfg.AccrualAddress
	url := fmt.Sprintf("http://%s/api/orders/%s", server, order)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса", err)
		return
	}
	do, err := client.Do(request)
	if err != nil {
		return
	}
	fmt.Println(do.Status)
	do.Body.Close()
}

func Sender(ctx context.Context, u usecase.OrderUseCase, cfg config.ServerConfig, rate time.Duration) {
	for {
		orders := GetOrders(ctx, u)
		for _, order := range orders {
			SendOrder(order, cfg)
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

func UpdateOrders(ctx context.Context, u usecase.OrderUseCase, b usecase.BalanceUseCase) {

}
