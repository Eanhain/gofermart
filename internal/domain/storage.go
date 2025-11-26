package domain

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

type Cache interface {
	InitSchema(ctx context.Context, log Logger) error
	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.User) (dto.User, error)
}

type Storage interface {
	InitSchema(ctx context.Context, log Logger) error
	RegisterUser(ctx context.Context, users dto.User) error
	CheckUser(ctx context.Context, users dto.User) (dto.User, error)
}
