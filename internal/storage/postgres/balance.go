package postgres

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	sq "github.com/Masterminds/squirrel"
)

func (ps *PersistStorage) InsertNewUserBalance(ctx context.Context, userID int, balance float64) error {
	sql, args, err := ps.psql.
		Insert("balance").
		Columns("user_id", "balance", "withdrawn").
		Values(userID, balance, 0).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't insert user balance", err)
		return err
	}
	ps.log.Infoln("Insert user balance", tag)
	return nil
}

func (ps *PersistStorage) GetUserBalance(ctx context.Context, userID int) (dto.Amount, error) {
	var balance dto.Amount

	sql, args, err := ps.psql.
		Select("balance", "withdrawn").
		From("balance").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return dto.Amount{}, err
	}

	row := ps.QueryRow(ctx, sql, args...)
	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
		return dto.Amount{}, err
	}
	return balance, nil
}

func (ps *PersistStorage) UpdateUserBalance(ctx context.Context, userID int, amount float64) error {
	sql, args, err := ps.psql.
		Update("balance").
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't update user balance", err)
		return err
	}
	ps.log.Infoln("Update user balance", tag)
	return nil
}

func (ps *PersistStorage) UpdateUserWithdrawn(ctx context.Context, userID int, amount float64) error {
	sql, args, err := ps.psql.
		Update("balance").
		Set("withdrawn", sq.Expr("withdrawn + ?", amount)).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't update user withdrawn", err)
		return err
	}
	ps.log.Infoln("Update user withdrawn", tag)
	return nil
}
