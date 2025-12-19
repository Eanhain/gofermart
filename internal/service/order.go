package service

import (
	"context"
	"fmt"
	"strconv"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	"github.com/theplant/luhn"
)

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
	if err != nil {
		return err
	}
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
	if accrual != 0 && status == "PROCESSED" {
		if err := s.c.UpdateUserBalance(ctx, id, accrual); err != nil {
			return err
		}
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

func (s *Service) GetUserWithdrawals(ctx context.Context, username string) (dto.Withdrawns, error) {
	id, err := s.c.GetUserID(ctx, username)
	if err != nil {
		return dto.Withdrawns{}, err
	}
	orders, err := s.c.GetUserOrdersWithdrawn(ctx, id)
	if err != nil {
		return dto.Withdrawns{}, err
	}
	return orders, nil
}
