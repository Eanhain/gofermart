package postgres

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Close()
}

type DMLUserStruct struct {
	DML  string
	Name string
}

var (
	InsertUser = DMLUserStruct{
		Name: "Add Users (batch)",
		DML: ` INSERT INTO users (username, hash) 
		VALUES ($1, $2)`,
	}
	InsertOrder = DMLUserStruct{
		Name: "Insert user order",
		DML: `INSERT INTO orders (ID, USER_ID, STATUS, ACCURAL)
		VALUES ($1, $2, 'NEW', 0)`,
	}
	selectUserID = `
		SELECT id FROM users 
		WHERE USERNAME = $1
	`
	selectUser = `
		SELECT username,hash FROM users
		WHERE USERNAME = $1
	`
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
			ID			BIGINT PRIMARY KEY,
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
			UPLOADED_AT TIMESTAMPTZ DEFAULT now()
		)`
)

type PersistStorage struct {
	PgxIface
	log domain.Logger
}

func InitialPersistStorage(ctx context.Context, log domain.Logger, connString string) (*PersistStorage, error) {
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	PersistStorageInstance := &PersistStorage{pgxPool, log}
	return PersistStorageInstance, nil
}

func (ps *PersistStorage) InitSchema(ctx context.Context, log domain.Logger) error {
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
	ps.log.Infoln("tables created")
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

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

func (ps *PersistStorage) InsertNewUserOrder(ctx context.Context, order int, userID int) error {
	tag, err := ps.Exec(ctx, InsertOrder.DML, order, userID)
	if err != nil {
		ps.log.Warnln("Can't insert user order", err)
		return err
	}
	ps.log.Infoln("Insert user order", tag)
	return nil

}
