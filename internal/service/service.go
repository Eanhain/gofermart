package service

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Service struct {
	c domain.Cache
}

func InitialService(c domain.Cache) *Service {
	return &Service{c: c}
}

func (s *Service) AddUser(ctx context.Context, user dto.UserArray) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}
