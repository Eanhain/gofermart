package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Storage interface {
	InitSchema(ctx context.Context, log Logger) error
	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.UserInput) (dto.User, error)
	InsertNewUserOrder(ctx context.Context, order string, userID int, status string, accrual float64) error
	GetUserID(ctx context.Context, user string) (int, error)
	GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error)
	GetUserBalance(ctx context.Context, userID int) (dto.Amount, error)
	CheckOrderNonExist(ctx context.Context, orders string) error
	CheckUserOrderNonExist(ctx context.Context, userID int, orders string) error
	InsertNewUserBalance(ctx context.Context, userID int, balance float64) error
}

type AccrualAPI interface {
	GetOrder(order string) (dto.OrderDesc, error)
}
