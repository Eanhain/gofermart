package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	hash "github.com/Eanhain/gofermart/internal/hash"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/theplant/luhn"
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

func (s *Service) AuthUser(ctx context.Context, user dto.UserInput) (bool, error) {
	tUser, err := s.c.CheckUser(ctx, user)
	if err != nil {
		return false, err
	}
	ok := hash.VerifyUserHash(s.log, user, tUser)
	return ok, nil
}

func (s *Service) RegUser(ctx context.Context, user dto.UserInput) error {
	var pgErr *pgconn.PgError
	hashedUser := hash.CreateUserHash(s.log, user)
	err := s.c.RegisterUser(ctx, hashedUser)

	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = domain.ErrConflict
	}

	return err
}

func (s *Service) CheckUserOrderDubl(ctx context.Context, userID int, order string) error {
	if err := s.c.CheckUserOrderNonExist(ctx, userID, order); err != nil {
		return err
	}
	if err := s.c.CheckOrderNonExist(ctx, order); err != nil {
		return err
	}
	return nil
}

func (s *Service) PostUserOrder(ctx context.Context, username string, order string) error {
	var (
		status  string
		accrual float64
	)
	orderDesc, err := s.aAPI.GetOrder(order)
	if err != nil {
		s.log.Warnln(err)
	}
	if orderDesc.Number != "" {
		order = orderDesc.Number
		accrual = orderDesc.Accrual
	}
	if orderDesc.Status != "" {
		status = orderDesc.Status
	} else {
		status = "NEW"
	}
	id, err := s.c.GetUserID(ctx, username)
	if err := s.CheckUserOrderDubl(ctx, id, order); err != nil {
		return err
	}
	if err != nil {
		return nil
	}
	if _, err = s.CheckOrderByLuna(ctx, order); err != nil {
		return err
	}

	if err := s.c.InsertNewUserOrder(ctx, order, id, status, accrual); err != nil {
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
	return false, fmt.Errorf("%w %v", domain.ErrOrderInvalid, order)
}

func (s *Service) GetUserOrders(ctx context.Context, username string) (dto.OrdersDesc, error) {
	id, err := s.c.GetUserID(ctx, username)
	if err != nil {
		return dto.OrdersDesc{}, err
	}
	orders, err := s.c.GetUserOrders(ctx, id)
	if err != nil {
		return dto.OrdersDesc{}, err
	}
	return orders, nil
}

func (s *Service) GetUserBalance(ctx context.Context, username string) (dto.Amount, error) {
	id, err := s.c.GetUserID(ctx, username)
	if err != nil {
		return dto.Amount{}, err
	}
	balance, err := s.c.GetUserBalance(ctx, id)
	if err != nil {
		return dto.Amount{}, err
	}
	return balance, nil
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
