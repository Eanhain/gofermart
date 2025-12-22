package migration

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
	entity "github.com/Eanhain/gofermart/internal/storage/entity"
	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"
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

type Migration struct {
	entity.PgxIface
	log  domain.Logger
	psql sq.StatementBuilderType
}

func InitialMigration(ctx context.Context, log domain.Logger, pgxPool *pgxpool.Pool) (*Migration, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	MigrationInstance := &Migration{pgxPool, log, psql}
	return MigrationInstance, nil
}

func (ps *Migration) InitSchema(ctx context.Context, log domain.Logger) error {
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
