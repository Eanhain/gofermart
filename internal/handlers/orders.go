package handlers

import (
	"errors"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

func (r *app) HandlerGetUserOrders(c *fiber.Ctx) error {
	var ordersOut dto.OrdersDescOut
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	orders, err := r.service.GetUserOrders(c.Context(), username)
	if err != nil {
		r.logger.Warnln("Can't get order", err)
		return err
	}
	copier.Copy(&ordersOut, orders)
	if err := c.JSON(ordersOut); err != nil {
		r.logger.Warnln("Can't get json", err)
		return err
	}
	return nil
}

func (r *app) HandlerPushOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	order := c.Body()
	orderStr := string(order)
	if err := r.service.PostUserOrder(c.Context(), username, orderStr); err != nil {
		if errors.Is(err, domain.ErrOrderExistWrongUser) {
			return fiber.ErrConflict
		} else if errors.Is(err, domain.ErrOrderExist) {
			return nil
		} else if errors.Is(err, domain.ErrOrderInvalid) {
			return fiber.ErrUnprocessableEntity
		} else if errors.Is(err, domain.ErrRequestCount) {
			return errors.Join(fiber.ErrTooManyRequests, err)
		}
		return fiber.ErrInternalServerError
	}
	c.SendStatus(fiber.StatusAccepted)
	return nil
}

func (r *app) HandlersWithdrawals(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["login"].(string)
	orders, err := r.service.GetUserWithdrawals(c.Context(), username)
	if errors.Is(err, domain.ErrEmptyOrdersList) {
		c.SendStatus(fiber.StatusNoContent)
		return nil
	} else if err != nil {
		r.logger.Warnln("Can't get order", err)
		return err
	}
	if err := c.JSON(orders); err != nil {
		r.logger.Warnln("Can't get json", err)
		return err
	}
	return nil
}
