package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Logger interface {
	Warnln(args ...any)
	Infoln(args ...any)
}

type Service interface {
}

type app struct {
	*fiber.App
	logger Logger
	server string
}

func InitialApp(log Logger, host string, port string) app {
	routeFiber := fiber.New()
	server := fmt.Sprintf("%v:%v", host, port)
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
