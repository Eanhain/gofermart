package service

import (
	"context"
	"fmt"
	"strconv"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	hash "github.com/Eanhain/gofermart/internal/hash"
	"github.com/theplant/luhn"
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

func (s *Service) PostUserOrder(ctx context.Context, username string, order string) error {
	id, err := s.c.GetUserID(ctx, username)
	if err != nil {
		return nil
	}
	orderInt, err := strconv.Atoi(order)
	if err != nil {
		return fmt.Errorf("can't convert order to int %w", err)
	}
	s.CheckOrderByLuna(ctx, order)
	if err := s.c.InsertNewUserOrder(ctx, orderInt, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) CheckOrderByLuna(ctx context.Context, order string) (bool, error) {
	orderInt, err := strconv.Atoi(order)
	if err != nil {
		return false, fmt.Errorf("can't convert order to int %w", err)
	}
	if ok := luhn.Valid(orderInt); ok {
		s.log.Infoln("Order id valid: ", order)
		return true, nil
	}
	s.log.Infoln("Order id not valid: ", order)
	return false, nil
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
