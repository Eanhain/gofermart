package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Service interface {
	AuthUser(ctx context.Context, user dto.UserInput) (bool, error)
	RegUser(ctx context.Context, user dto.UserInput) error
	PostUserOrders(ctx context.Context, user dto.User) error
	GetUserOrders(ctx context.Context, user dto.User) error
	GetUserBalance(ctx context.Context, user dto.User) error
	GetUserWithdrawals(ctx context.Context, user dto.User) error
	PostUserWithdraw(ctx context.Context, user dto.User) error
}
