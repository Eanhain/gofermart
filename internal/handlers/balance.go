package handlers

import (
	"context"
	"errors"

	dto "github.com/Eanhain/gofermart/internal/api"
	"github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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

func (r *app) HandlersUserBalanceWithdraw(c *fiber.Ctx) error {
	var balance dto.Withdrawn
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	if err := c.BodyParser(&balance); err != nil {
		r.logger.Warnln("can't parse json withdrawn order")
		return err
	}
	err := r.service.PostUserWithdrawOrder(context.TODO(), username, balance)
	if errors.Is(err, domain.ErrBalanceWithdrawn) {
		return fiber.ErrPaymentRequired
	} else if errors.Is(err, domain.ErrOrderInvalid) {
		return fiber.ErrUnprocessableEntity
	} else if err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}
