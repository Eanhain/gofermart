package service

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

type Service struct {
	c domain.Cache
}

func InitialService(ctx context.Context, c domain.Cache, log domain.Logger) (*Service, error) {
	if err := c.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("couldn't initialize service layer: %w", err)
	}
	return &Service{c: c}, nil
}

func (s *Service) AuthUser(ctx context.Context, user dto.User) error {
	_, err := s.c.CheckUser(ctx, user)
	return err
}

func (s *Service) RegUser(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}

// TODO
func (s *Service) PostUserOrders(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}

// TODO
func (s *Service) GetUserOrders(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}

// TODO
func (s *Service) GetUserBalance(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}

// TODO
func (s *Service) GetUserWithdrawals(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}

// TODO
func (s *Service) PostUserWithdraw(ctx context.Context, user dto.User) error {
	err := s.c.RegisterUser(ctx, user)
	return err
}
