package handlers

import (
	"context"
	"fmt"
	"time"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	return app{routeFiber, service, log, server, jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte(supersecret)}}}
}

func (r *app) StartServer(ctx context.Context) error {
	err := r.Listen(r.server)
	return err
}

func (r *app) CreateHandlers(ctx context.Context) error {
	r.Post("/api/user/register", r.HandlerRegUser)
	r.Post("/api/user/login", r.LoginJWT)
	r.Use(jwtware.New(r.jwtConf))
	r.Post("/api/user/orders", r.HandlerPushOrder)
	r.Get("/api/user/orders", r.HandlerGetUserOrders)
	r.Get("/api/user/balance", r.HandlerGetUserBalance)
	err := r.Listen(r.server)
	return err
}

func (r *app) LoginJWT(c *fiber.Ctx) error {
	username, err := r.AuthUser(c)
	if err != nil {
		r.logger.Warnln("Can't get jwt token: ", err)
		c.SendStatus(fiber.StatusUnauthorized)
		return fmt.Errorf("can't get jwt token: %w", err)
	}
	claims := jwt.MapClaims{
		"login": username,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenJWT, err := token.SignedString(r.jwtConf.SigningKey.Key)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": tokenJWT})
}

func (r *app) HandlerGetUserBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	orders, err := r.service.GetUserBalance(context.TODO(), username)
	if err != nil {
		r.logger.Warnln("Can't get balance", err)
		return err
	}

	if err := c.JSON(orders); err != nil {
		r.logger.Warnln("Can't get json", err)
		return err
	}
	return nil
}

func (r *app) HandlerGetUserOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	orders, err := r.service.GetUserOrders(context.TODO(), username)
	if err != nil {
		r.logger.Warnln("Can't get order", err)
		return err
	}

	if err := c.JSON(orders); err != nil {
		r.logger.Warnln("Can't get json", err)
		return err
	}
	return nil
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

func (r *app) AuthUser(c *fiber.Ctx) (string, error) {
	var user dto.UserInput
	if err := c.BodyParser(&user); err != nil {
		r.logger.Warnln("can't parse body for registr", err)
		return "", err
	}
	if ok, err := r.service.AuthUser(context.TODO(), user); err != nil || !ok {
		return "", fmt.Errorf("user not auth %v", user.Login)
	}
	return user.Login, nil
}

func (r *app) HandlerPushOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	order := c.Body()
	orderStr := string(order)
	if err := r.service.PostUserOrder(context.TODO(), username, orderStr); err != nil {
		return err
	}
	return nil
}
