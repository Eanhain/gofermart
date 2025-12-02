package cache

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Cache struct {
	storage   domain.Storage
	log       domain.Logger
	authUsers map[string]bool
}

func InitCache(ctx context.Context, log domain.Logger, storage domain.Storage) (Cache, error) {
	return Cache{storage: storage,
		log: log, authUsers: make(map[string]bool)}, nil
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
	c.authUsers[user.Login] = true
	return tUser, nil
}

func (c *Cache) GetUserID(ctx context.Context, username string) (int, error) {
	id, err := c.storage.GetUserID(ctx, username)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Cache) InsertNewUserOrder(ctx context.Context, order string, userID int) error {
	err := c.storage.InsertNewUserOrder(ctx, order, userID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) CheckAuthUser(user string) bool {
	return c.authUsers[user]
}

func (c *Cache) GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error) {
	orders, err := c.storage.GetUserOrders(ctx, userID)
	if err != nil {
		return dto.OrdersDesc{}, err
	}
	return orders, nil
}
