package service

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Service struct {
	storage domain.Storage
	aAPI    domain.AccrualAPI
	log     domain.Logger
}

func InitialService(ctx context.Context, storage domain.Storage, accrual domain.AccrualAPI, log domain.Logger) (*Service, error) {
	if err := storage.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("couldn't initialize service layer: %w", err)
	}
	return &Service{storage: storage, aAPI: accrual, log: log}, nil
}
