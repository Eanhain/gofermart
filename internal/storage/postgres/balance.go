package postgres

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

func (ps *PersistStorage) InsertNewUserBalance(ctx context.Context, userID int, balance float64) error {
	tag, err := ps.Exec(ctx, InsertUserBalance.DML, userID, balance, 0)
	if err != nil {
		ps.log.Warnln("Can't insert user balance", err)
		return err
	}
	ps.log.Infoln("Insert user balance", tag)
	return nil

}

func (ps *PersistStorage) GetUserBalance(ctx context.Context, userID int) (dto.Amount, error) {
	var balance dto.Amount
	row := ps.QueryRow(ctx, selectUserBalance, userID)
	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
		return dto.Amount{}, err
	}
	return balance, nil

}
