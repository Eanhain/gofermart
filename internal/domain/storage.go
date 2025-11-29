package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Cache interface {
	InitSchema(ctx context.Context, log Logger) error
	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.UserInput) (dto.User, error)
	InsertNewUserOrder(ctx context.Context, order int, userID int) error
	GetUserID(ctx context.Context, user string) (int, error)
	CheckAuthUser(user string) bool
}

type Storage interface {
	InitSchema(ctx context.Context, log Logger) error
	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.UserInput) (dto.User, error)
	InsertNewUserOrder(ctx context.Context, order int, userID int) error
	GetUserID(ctx context.Context, user string) (int, error)
}
