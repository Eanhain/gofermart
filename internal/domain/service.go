package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Service interface {
	AuthUser(ctx context.Context, user dto.UserInput) (bool, error)
	RegUser(ctx context.Context, user dto.UserInput) error
	GetUserWithdrawals(ctx context.Context, username string) (*dto.Withdrawns, error)
	PostUserOrder(ctx context.Context, username string, order string) error
	GetUserOrders(ctx context.Context, username string) (dto.OrdersDesc, error)
	GetUserBalance(ctx context.Context, username string) (dto.Amount, error)

	PostUserWithdrawOrder(ctx context.Context, username string, order dto.Withdrawn) error
}
