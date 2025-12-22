package service

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Service struct {
	auth    domain.Auth
	balance domain.Balance
	orders  domain.Orders
	migr    domain.Migration
	aAPI    domain.AccrualAPI
	log     domain.Logger
}

func InitialService(ctx context.Context, auth domain.Auth, balance domain.Balance, orders domain.Orders,
	migr domain.Migration, accrual domain.AccrualAPI, log domain.Logger) (*Service, error) {
	if err := migr.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("couldn't initialize service layer: %w", err)
	}
	return &Service{auth: auth, balance: balance,
		orders: orders, migr: migr,
		aAPI: accrual, log: log}, nil
}
