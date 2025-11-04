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

type route struct {
	*fiber.App
	logger Logger
	server string
}

func InitialHandler(log Logger, host string, port string) route {
	routeFiber := fiber.New()
	server := fmt.Sprintf("%v:%v", host, port)
	return route{routeFiber, log, server}
}

func (r *route) StartServer(ctx context.Context) error {
	err := r.Listen(r.server)
	return err
}
