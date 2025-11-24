package handlers

import (
	"context"

	domain "github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type app struct {
	*fiber.App
	logger domain.Logger
	server string
}

func InitialApp(log domain.Logger, service domain.Service, server string) app {
	routeFiber := fiber.New()
	return app{routeFiber, log, server}
}

func (r *app) StartServer(ctx context.Context) error {
	err := r.Listen(r.server)
	return err
}

func (r *app) CreateHandlers(ctx context.Context) error {
	r.Post("/api/user/register")
	err := r.Listen(r.server)
	return err
}
