package handlers

import (
	"context"

	domain "github.com/Eanhain/gofermart/internal/domain"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type app struct {
	*fiber.App
	service domain.Service
	logger  domain.Logger
	server  string
	jwtConf jwtware.Config
}

func InitialApp(log domain.Logger, service domain.Service, server string, supersecret string) app {
	routeFiber := fiber.New()
	return app{routeFiber, service, log, server, jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(supersecret)}, ErrorHandler: jwtError}}
}

func (r *app) StartServer(ctx context.Context) error {
	err := r.Listen(r.server)
	return err
}

// rate limit error - 429 - accrual
func (r *app) CreateHandlers(ctx context.Context) error {
	r.Post("/api/user/register", r.HandlerRegUser)
	r.Post("/api/user/login", r.LoginJWT)
	r.Use(jwtware.New(r.jwtConf))
	r.Post("/api/user/orders", r.HandlerPushOrder)
	r.Get("/api/user/orders", r.HandlerGetUserOrders)
	r.Get("/api/user/balance", r.HandlerGetUserBalance)
	r.Post("/api/user/balance/withdraw", r.HandlersUserBalanceWithdraw)
	err := r.Listen(r.server)
	return err
}
