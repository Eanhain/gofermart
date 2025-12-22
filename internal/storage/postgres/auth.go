package postgres

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	sq "github.com/Masterminds/squirrel"
)

func (ps *PersistStorage) RegisterUser(ctx context.Context, user dto.User) error {
	sql, args, err := ps.psql.
		Insert("users").
		Columns("username", "hash").
		Values(user.Login, user.Hash).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql: %w", err)
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("can't register user %w, with user %v", err, user.Login)
	}

	ps.log.Infoln("User registered, rows affected:", tag.RowsAffected(), "user:", user.Login)
	return nil
}

func (ps *PersistStorage) CheckUser(ctx context.Context, untrustedUser dto.UserInput) (dto.User, error) {
	var orUser dto.User

	sql, args, err := ps.psql.
		Select("username", "hash").
		From("users").
		Where(sq.Eq{"username": untrustedUser.Login}).
		ToSql()

	if err != nil {
		return dto.User{}, fmt.Errorf("failed to build sql: %w", err)
	}

	row := ps.QueryRow(ctx, sql, args...)
	if err := row.Scan(&orUser.Login, &orUser.Hash); err != nil {
		return dto.User{}, err
	}

	ps.log.Infoln("Get trust user from db:", orUser.Login)
	return orUser, nil
}

func (ps *PersistStorage) GetUserID(ctx context.Context, username string) (int, error) {
	var id int

	sql, args, err := ps.psql.
		Select("id").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return -1, fmt.Errorf("failed to build sql: %w", err)
	}

	row := ps.QueryRow(ctx, sql, args...)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
