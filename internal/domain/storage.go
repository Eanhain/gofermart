package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Storage interface {
	InitSchema(ctx context.Context, log Logger) error

	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.UserInput) (dto.User, error)
	GetUserID(ctx context.Context, user string) (int, error)

	InsertNewUserOrder(ctx context.Context, order string, userID int, status string, accrual float64) error
	InsertOrderWithdrawn(ctx context.Context, userID int, order dto.Withdrawn) error
	CheckUserOrderNonExist(ctx context.Context, userID int, orders string) error
	GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error)
	CheckOrderNonExist(ctx context.Context, orders string) error
	GetUserOrdersWithdrawn(ctx context.Context, userID int) (dto.Withdrawns, error)

	InsertNewUserBalance(ctx context.Context, userID int, balance float64) error
	GetUserBalance(ctx context.Context, userID int) (dto.Amount, error)
	UpdateUserBalance(ctx context.Context, userID int, balance float64) error
	UpdateUserWithdrawn(ctx context.Context, userID int, withdrawn float64) error
}

type AccrualAPI interface {
	GetOrder(order string) (dto.OrderDesc, error)
}
