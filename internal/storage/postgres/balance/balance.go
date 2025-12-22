package balance

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	"github.com/Eanhain/gofermart/internal/domain"
	entity "github.com/Eanhain/gofermart/internal/storage/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Balance struct {
	entity.PgxIface
	log  domain.Logger
	psql sq.StatementBuilderType
}

func InitialBalance(ctx context.Context, log domain.Logger, pgxPool *pgxpool.Pool) (*Balance, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	MigrationInstance := &Balance{pgxPool, log, psql}
	return MigrationInstance, nil
}

func (ps *Balance) InsertNewUserBalance(ctx context.Context, userID int, balance float64) error {
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

func (ps *Balance) GetUserBalance(ctx context.Context, userID int) (dto.Amount, error) {
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

func (ps *Balance) UpdateUserBalance(ctx context.Context, userID int, amount float64) error {
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

func (ps *Balance) UpdateUserWithdrawn(ctx context.Context, userID int, amount float64) error {
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
