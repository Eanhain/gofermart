package service

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	hash "github.com/Eanhain/gofermart/internal/hash"
)

type Service struct {
	c   domain.Cache
	log domain.Logger
}

func InitialService(ctx context.Context, c domain.Cache, log domain.Logger) (*Service, error) {
	if err := c.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("couldn't initialize service layer: %w", err)
	}
	return &Service{c: c, log: log}, nil
}

func (s *Service) AuthUser(ctx context.Context, user dto.UserInput) (bool, error) {
	tUser, err := s.c.CheckUser(ctx, user)
	if err != nil {
		return false, err
	}
	ok := hash.VerifyUserHash(s.log, user, tUser)
	return ok, nil
}

func (s *Service) RegUser(ctx context.Context, user dto.UserInput) error {
	hashedUser := hash.CreateUserHash(s.log, user)
	err := s.c.RegisterUser(ctx, hashedUser)
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
