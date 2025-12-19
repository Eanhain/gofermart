package postgres

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
)

var (
	InsertUserBalance = DMLUserStruct{
		Name: "Insert user balance",
		DML: `INSERT INTO balance (user_id, balance, withdrawn)
		VALUES ($1, $2, $3)`,
	}

	UpdateUserBalance = DMLUserStruct{
		Name: "Update user balance",
		DML: `UPDATE balance
		SET balance = balance + $1
		WHERE user_id = $2;`,
	}

	UpdateUserWithdrawn = DMLUserStruct{
		Name: "Update user withdrawn",
		DML: `UPDATE balance
		SET withdrawn = withdrawn + $1
		WHERE user_id = $2;`,
	}

	selectUserBalance = `
		select balance as current, withdrawn 
		from balance 
		where user_id = $1;
	`
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

func (ps *PersistStorage) UpdateUserBalance(ctx context.Context, userID int, balance float64) error {
	tag, err := ps.Exec(ctx, UpdateUserBalance.DML, userID)
	if err != nil {
		ps.log.Warnln("Can't update user balance", err)
		return err
	}
	ps.log.Infoln("Update user balance", tag)
	return nil

}

func (ps *PersistStorage) UpdateUserWithdrawn(ctx context.Context, userID int, withdrawn float64) error {
	tag, err := ps.Exec(ctx, UpdateUserWithdrawn.DML, userID)
	if err != nil {
		ps.log.Warnln("Can't update user withdrawn", err)
		return err
	}
	ps.log.Infoln("Update user withdrawn", tag)
	return nil

}
