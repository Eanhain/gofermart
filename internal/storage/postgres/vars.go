package postgres

import (
	"context"

	"github.com/Eanhain/gofermart/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		DML: `INSERT INTO balance (user_id, balance, withdrawn)
		VALUES ($1, $2, $3)`,
	}

	UpdateUserBalance = DMLUserStruct{
		Name: "Update user balance",
		DML: `UPDATE balance
		SET balance = balance + $1
		WHERE user_id = $2;`,
	}

	InsertOrderWithdrawn = DMLUserStruct{
		Name: "Insert user withdrawn order",
		DML: `INSERT INTO withdraw_orders (user_id, order_id, sum)
		VALUES ($1, $2, $3)`,
	}

	UpdateUserWithdrawn = DMLUserStruct{
		Name: "Update user withdrawn",
		DML: `UPDATE balance
		SET withdrawn = withdrawn + $1
		WHERE user_id = $2;`,
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
	selectWithdrawnUserOrders = `
		select order_id as order, sum, UPLOADED_AT as processed_at from withdraw_orders
		where user_id = $1
		order by processed_at desc
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

type PersistStorage struct {
	PgxIface
	log domain.Logger
}
