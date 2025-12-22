package handlers

import (
	"errors"
	"fmt"
	"time"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func (r *app) CreateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"login": username,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenJWT, err := token.SignedString(r.jwtConf.SigningKey.Key)
	if err != nil {
		return "", err
	}
	return tokenJWT, nil
}

func (r *app) LoginJWT(c *fiber.Ctx) error {
	username, err := r.AuthUser(c)
	if err != nil {
		r.logger.Warnln("Can't auth user: ", err)
		return fiber.ErrUnauthorized
	}
	tokenJWT, err := r.CreateJWT(username)
	if err != nil {
		r.logger.Warnln("Can't create jwt token", err)
		return fiber.ErrInternalServerError
	}
	c.Set("Authorization", "Bearer "+tokenJWT)
	return c.JSON(fiber.Map{"token": tokenJWT})
}

func (r *app) HandlerRegUser(c *fiber.Ctx) error {
	var user dto.UserInput
	if err := c.BodyParser(&user); err != nil {
		r.logger.Warnln("can't parse body for registr", err)
		return fiber.ErrInternalServerError
	}
	if err := r.service.RegUser(c.Context(), user); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return fiber.ErrConflict
		} else {
			return fiber.ErrInternalServerError
		}
	}
	r.service.AuthUser(c.Context(), user)
	if ok, err := r.service.AuthUser(c.Context(), user); err != nil || !ok {
		r.logger.Warnln("Can't auth user: ", err)
		return fiber.ErrBadRequest
	}

	tokenJWT, err := r.CreateJWT(user.Login)
	if err != nil {
		r.logger.Warnln("Can't create jwt token", err)
		return fiber.ErrInternalServerError
	}

	c.Set("Authorization", "Bearer "+tokenJWT)

	return nil
}

func (r *app) AuthUser(c *fiber.Ctx) (string, error) {
	var user dto.UserInput
	if err := c.BodyParser(&user); err != nil {
		r.logger.Warnln("can't parse body for registr", err)
		return "", err
	}
	if ok, err := r.service.AuthUser(c.Context(), user); err != nil || !ok {
		return "", fmt.Errorf("user not auth %v", user.Login)
	}
	return user.Login, nil
}
