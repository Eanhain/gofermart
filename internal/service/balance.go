package service

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	"github.com/Eanhain/gofermart/internal/domain"
)

func (s *Service) GetUserBalance(ctx context.Context, username string) (dto.Amount, error) {
	id, err := s.storage.GetUserID(ctx, username)
	if err != nil {
		return dto.Amount{}, err
	}
	balance, err := s.storage.GetUserBalance(ctx, id)
	if err != nil {
		return dto.Amount{}, err
	}
	return balance, nil
}

func (s *Service) PostUserWithdrawOrder(ctx context.Context, username string, order dto.Withdrawn) error {

	userID, err := s.storage.GetUserID(ctx, username)
	if err != nil {
		return err
	}
	if err := s.CheckUserOrderDubl(ctx, userID, order.Order); err != nil {
		return err
	}
	if _, err = s.CheckOrderByLuna(ctx, order.Order); err != nil {
		return err
	}

	if err := s.storage.InsertNewUserOrder(ctx, order.Order, userID, "NEW", 0); err != nil {
		return err
	}
	amount, err := s.storage.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if amount.Current < order.Sum {
		return domain.ErrBalanceWithdrawn
	}
	if err := s.storage.UpdateUserBalance(ctx, userID, order.Sum*-1); err != nil {
		s.log.Warnln("can't update user balance", username, amount, order)
		return err
	}
	if err := s.storage.UpdateUserWithdrawn(ctx, userID, order.Sum); err != nil {
		s.log.Warnln("can't update user withdrawn", username, amount, order)
		return err
	}
	if err := s.storage.InsertOrderWithdrawn(ctx, userID, order); err != nil {
		s.log.Warnln("can't insert user withdrawn", username, amount, order)
		return err
	}

	return err
}
