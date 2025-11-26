package cache

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Cache struct {
	storage domain.Storage
}

func InitCache(ctx context.Context, storage domain.Storage) (Cache, error) {
	return Cache{storage}, nil
}

func (c *Cache) InitSchema(ctx context.Context, log domain.Logger) error {
	if err := c.storage.InitSchema(ctx, log); err != nil {
		return err
	}
	return nil
}

func (c *Cache) RegisterUser(ctx context.Context, user dto.User) error {
	if err := c.storage.RegisterUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (c *Cache) CheckUser(ctx context.Context, user dto.User) (dto.User, error) {
	user, err := c.storage.CheckUser(ctx, user)
	if err == nil {
		return dto.User{}, err
	}
	return user, nil
}
