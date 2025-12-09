package domain

import (
	"errors"
)

var (
	ErrUserExist           = errors.New("user exists")
	ErrOrderExist          = errors.New("user already create order")
	ErrOrderExistWrongUser = errors.New("order create for another user")
	ErrJWTToken            = errors.New("can't create jwt token")
	ErrAuthUser            = errors.New("can't auth user")
)
