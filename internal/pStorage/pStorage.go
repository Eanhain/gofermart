package pStorage

import (
	"context"
	"fmt"

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

func ConnectToPersistStorage(ctx context.Context, log Logger, user string, passw string, host string, port string, schema string) (*PersistStorage, error) {
	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, passw, host, port, schema)
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	if err := pgxPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping doesn't work %w", err)
	}
	PersistStorageInstance := &PersistStorage{pgxPool}
	if err := PersistStorageInstance.InitSchema(log, ctx); err != nil {
		return nil, fmt.Errorf("can't initialize schema, %w", err)
	}
	return PersistStorageInstance, nil
}

func (ps *PersistStorage) InitSchema(log Logger, ctx context.Context) error {
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

func (ps *PersistStorage) Add(user dto.User) error {
	return nil
}

func (ps *PersistStorage) Del(user dto.User) error {
	return nil
}

func (ps *PersistStorage) MultipleAdd(user dto.UserArray) error {
	return nil
}

func (ps *PersistStorage) MultipleDel(user dto.UserArray) error {
	return nil
}

func (ps *PersistStorage) List() ([]dto.UserArray, error) {
	return []dto.UserArray{}, nil
}

func (ps *PersistStorage) Connect() {

}
