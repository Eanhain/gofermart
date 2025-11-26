package cache

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Cache struct {
	storage domain.Storage
	log     domain.Logger
}

func InitCache(ctx context.Context, log domain.Logger, storage domain.Storage) (Cache, error) {
	return Cache{storage: storage,
		log: log}, nil
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

func (c *Cache) CheckUser(ctx context.Context, user dto.UserInput) (dto.User, error) {
	tUser, err := c.storage.CheckUser(ctx, user)
	if err != nil {
		c.log.Warnln("Check user error", tUser, err)
		return dto.User{}, err
	}
	return tUser, nil
}
