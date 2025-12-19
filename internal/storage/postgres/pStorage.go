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
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
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
		VALUES ($1, $2, $3, $4)`,
	}
	InsertUserBalance = DMLUserStruct{
		Name: "Insert user balance",
		DML: `INSERT INTO balance (user_id, balance)
		VALUES ($1, $2)`,
	}

	selectUserID = `
		SELECT id FROM users 
		WHERE USERNAME = $1
	`
	selectUser = `
		SELECT username,hash FROM users
		WHERE USERNAME = $1
	`
	selectUserOrders = `
		select ID as number, status, accural accrual, uploaded_at 
		from orders 
		where user_id = $1;
	`
	selectUserBalance = `
		select balance as current, withdrawn 
		from balance 
		where user_id = $1;
	`
	selectUserOrder = `
		select ID as number 
		from orders 
		where user_id = $1 and id = $2;
	`
	selectOrder = `
		select ID as number 
		from orders 
		where id = $1;
	`
)

// TODO Migration
// sql-c?
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
		SUM			REAL
	)`
)

type PersistStorage struct {
	PgxIface
	log domain.Logger
}

// select for update -> coins
func InitialPersistStorage(ctx context.Context, log domain.Logger, connString string) (*PersistStorage, error) {
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	PersistStorageInstance := &PersistStorage{pgxPool, log}
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

func (ps *PersistStorage) InsertNewUserOrder(ctx context.Context, order string, userID int, status string, accrual float64) error {
	tag, err := ps.Exec(ctx, InsertOrder.DML, order, userID, status, accrual)
	if err != nil {
		ps.log.Warnln("Can't insert user order", err)
		return err
	}
	ps.log.Infoln("Insert user order", tag)
	return nil

}

func (ps *PersistStorage) InsertNewUserBalance(ctx context.Context, userID int, balance float64) error {
	tag, err := ps.Exec(ctx, InsertUserBalance.DML, userID, balance)
	if err != nil {
		ps.log.Warnln("Can't insert user balance", err)
		return err
	}
	ps.log.Infoln("Insert user balance", tag)
	return nil

}

func (ps *PersistStorage) GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error) {
	var orders dto.OrdersDesc
	var order dto.OrderDesc
	rows, err := ps.Query(ctx, selectUserOrders, userID)
	if err != nil {
		ps.log.Warnln("Can't get user order", err)
		return dto.OrdersDesc{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			ps.log.Warnln("can't parse order", err)
			return dto.OrdersDesc{}, err
		}
		orders = append(orders, order)
	}

	ps.log.Infoln("Get user orders", orders)
	return orders, nil

}

func (ps *PersistStorage) GetUserBalance(ctx context.Context, userID int) (dto.Amount, error) {
	var balance dto.Amount
	row := ps.QueryRow(ctx, selectUserBalance, userID)
	if err := row.Scan(&balance.Current, &balance.Withdrawn); err != nil {
		return dto.Amount{}, err
	}
	return balance, nil

}

func (ps *PersistStorage) CheckUserOrderNonExist(ctx context.Context, userID int, orders string) error {
	var orderOut string
	out := ps.QueryRow(ctx, selectUserOrder, userID, orders)
	out.Scan(&orderOut)
	if orderOut != "" {
		return domain.ErrOrderExist
	}
	return nil
}

func (ps *PersistStorage) CheckOrderNonExist(ctx context.Context, orders string) error {
	var orderOut string
	out := ps.QueryRow(ctx, selectOrder, orders)
	out.Scan(&orderOut)
	if orderOut != "" {
		return domain.ErrOrderExistWrongUser
	}
	return nil
}
