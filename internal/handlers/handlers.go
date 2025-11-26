package handlers

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type app struct {
	*fiber.App
	service domain.Service
	logger  domain.Logger
	server  string
}

func InitialApp(log domain.Logger, service domain.Service, server string) app {
	routeFiber := fiber.New()
	return app{routeFiber, service, log, server}
}

func (r *app) StartServer(ctx context.Context) error {
	err := r.Listen(r.server)
	return err
}

func (r *app) CreateHandlers(ctx context.Context) error {
	r.Post("/api/user/register", r.HandlerRegUser)
	r.Post("/api/user/login", r.HandlerAuthUser)
	err := r.Listen(r.server)
	return err
}

func (r *app) HandlerRegUser(c *fiber.Ctx) error {
	var user dto.UserInput
	if err := c.BodyParser(&user); err != nil {
		r.logger.Warnln("can't parse body for registr", err)
		return nil
	}
	if err := r.service.RegUser(context.TODO(), user); err != nil {
		return err
	}
	return nil
}

func (r *app) HandlerAuthUser(c *fiber.Ctx) error {
	var user dto.UserInput
	if err := c.BodyParser(&user); err != nil {
		r.logger.Warnln("can't parse body for registr", err)
		return nil
	}
	if _, err := r.service.AuthUser(context.TODO(), user); err != nil {
		return err
	}
	return nil
}
