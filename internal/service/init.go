package service

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Service struct {
	c    domain.Storage
	aAPI domain.AccrualAPI
	log  domain.Logger
}

// TODO new service (whats service?)
func InitialService(ctx context.Context, c domain.Storage, accrual domain.AccrualAPI, log domain.Logger) (*Service, error) {
	if err := c.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("couldn't initialize service layer: %w", err)
	}
	return &Service{c: c, aAPI: accrual, log: log}, nil
}
