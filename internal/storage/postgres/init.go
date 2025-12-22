package postgres

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

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
			ID			TEXT PRIMARY KEY,
			USER_ID 	INTEGER REFERENCES users (ID),
			STATUS 		TEXT NOT NULL,
			ACCURAL 	REAL NOT NULL,
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
	ddlBalance = `
		CREATE TABLE IF NOT EXISTS balance (
			ID			SERIAL PRIMARY KEY,
			USER_ID 	INTEGER REFERENCES users (ID),
			BALANCE 	REAL NOT NULL,
			WITHDRAWN	REAL,
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
	ddlWithdrawOrders = `
	CREATE TABLE IF NOT EXISTS withdraw_orders (
		ID			SERIAL PRIMARY KEY,
		USER_ID 	INTEGER REFERENCES users (ID),
		ORDER_ID	TEXT REFERENCES orders (ID),
		SUM			REAL,
		UPLOADED_AT TIMESTAMPTZ DEFAULT now()
	)`
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

type DMLUserStruct struct {
	DML  string
	Name string
}

type PersistStorage struct {
	PgxIface
	log  domain.Logger
	psql sq.StatementBuilderType
}

// TODO Migration
// sql-c?

// select for update -> coins
func InitialPersistStorage(ctx context.Context, log domain.Logger, connString string) (*PersistStorage, error) {
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	PersistStorageInstance := &PersistStorage{pgxPool, log, psql}
	return PersistStorageInstance, nil
}

func (ps *PersistStorage) InitSchema(ctx context.Context, log domain.Logger) error {
	ddls := []string{ddlUsers, ddlOrders, ddlBalance, ddlWithdrawOrders}
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
	ps.log.Infoln("tables created")
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
