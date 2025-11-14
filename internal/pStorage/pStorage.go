package pStorage

import (
	"context"
	"fmt"
	"log"

	dto "github.com/Eanhain/gofermart/internal/api"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type Logger interface {
	Warnln(args ...any)
	Infoln(args ...any)
}

const (
	ddlUsers = ` 
		CREATE TABLE IF NOT EXISTS users (
			ID      SERIAL PRIMARY KEY,
			USERNAME 	TEXT NOT NULL UNIQUE,
			HASH  	TEXT NOT NULL UNIQUE,
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
	ddlOrders = `
		CREATE TABLE IF NOT EXISTS orders (
			ID			SERIAL PRIMARY KEY,
			USER_ID 	INTEGER REFERENCES users (ID),
			STATUS 		TEXT NOT NULL,
			ACCURAL 	REAL,
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
	ddlBalance = `
		CREATE TABLE IF NOT EXISTS balance (
			ID			SERIAL PRIMARY KEY,
			USER_ID 	INTEGER REFERENCES users (ID),
			BALANCE 	REAL NOT NULL,
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
)

type PersistStorage struct {
	*pgxpool.Pool
}

func ConnectToPersistStorage(ctx context.Context, log Logger, connString string) (*PersistStorage, error) {
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	if err := pgxPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping doesn't work %w", err)
	}
	PersistStorageInstance := &PersistStorage{pgxPool}
	if err := PersistStorageInstance.InitSchema(ctx, log); err != nil {
		return nil, fmt.Errorf("can't initialize schema, %w", err)
	}
	return PersistStorageInstance, nil
}

func (ps *PersistStorage) InitSchema(ctx context.Context, log Logger) error {
	ddls := []string{ddlUsers, ddlOrders, ddlBalance}
	tx, err := ps.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("can't complete ddl transaction")
	}
	for _, ddl := range ddls {
		commandTag, err := tx.Exec(ctx, ddl)
		log.Infoln(commandTag)
		if err != nil {
			return fmt.Errorf("couldn't parse DDL output %w", err)
		}
	}
	log.Infoln("tables created")
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PersistStorage) RegisterUser(ctx context.Context, users dto.UserArray) error {
	tx, err := ps.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("can't register user (transaction start) %w", err)
	}
	prepareState, err := tx.Prepare(ctx, "Add Users (batch)", `
		INSERT INTO users (username, hash) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("can't prepare statement %w", err)
	}

	batch := pgx.Batch{}

	for _, user := range users {
		batch.Queue(prepareState.Name, user.Login, user.Hash)
	}

	res := tx.SendBatch(ctx, &batch)

	ct, err := res.Exec()
	if err != nil {
		return fmt.Errorf("error with sending batch data %w", err)
	}
	log.Println("batch data was sending, rows:", ct.RowsAffected())

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (ps *PersistStorage) CheckUserPermissions(ctx context.Context, untrustedUser dto.User) (dto.User, error) {
	var orUser dto.User
	sql := `
		SELECT * FROM users
		WHERE USERNAME = "$1"
	`
	row := ps.QueryRow(ctx, sql, untrustedUser.Login)

	if err := row.Scan(&orUser); err != nil {
		return dto.User{}, err
	}

	return orUser, nil
}
