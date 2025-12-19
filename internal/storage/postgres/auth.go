package postgres

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	pgx "github.com/jackc/pgx/v5"
)

var (
	selectUser = `
		SELECT username,hash FROM users
		WHERE USERNAME = $1
	`
	selectUserID = `
		SELECT id FROM users 
		WHERE USERNAME = $1
	`
	InsertUser = DMLUserStruct{
		Name: "Add Users (batch)",
		DML: ` INSERT INTO users (username, hash) 
		VALUES ($1, $2)`,
	}
)

func (ps *PersistStorage) RegisterUser(ctx context.Context, user dto.User) error {
	tx, err := ps.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("can't register user (transaction start) %w", err)
	}
	prepareState, err := tx.Prepare(ctx, InsertUser.Name, InsertUser.DML)
	if err != nil {
		return fmt.Errorf("can't prepare statement %w", err)
	}

	batch := &pgx.Batch{}

	batch.Queue(prepareState.Name, user.Login, user.Hash)

	res := tx.SendBatch(ctx, batch)
	ct, err := res.Exec()
	if err != nil {
		return fmt.Errorf("error with sending batch data %w, with user %v", err, user)
	}

	res.Close()

	if err := tx.Commit(ctx); err != nil {
		ps.log.Warnln("cannot commit register tx", err)
		return err
	}

	ps.log.Infoln("batch data was sending, rows:", ct.RowsAffected(), "for user:", user.Login)

	return nil
}

func (ps *PersistStorage) CheckUser(ctx context.Context, untrustedUser dto.UserInput) (dto.User, error) {
	var orUser dto.User

	row := ps.QueryRow(ctx, selectUser, untrustedUser.Login)

	if err := row.Scan(&orUser.Login, &orUser.Hash); err != nil {
		return dto.User{}, err
	}

	ps.log.Infoln("Get trust user from db:", orUser.Login)
	return orUser, nil
}

func (ps *PersistStorage) GetUserID(ctx context.Context, username string) (int, error) {
	var id int
	row := ps.QueryRow(ctx, selectUserID, username)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
