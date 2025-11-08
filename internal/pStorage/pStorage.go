package pStorage

import (
	"context"
	"fmt"

	dto "github.com/Eanhain/gofermart/internal/api"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

const initSchema = `

CREATE TABLE IF NOT EXISTS users (
    ID      BIGINT PRIMARY KEY,
	USER 	TEXT NOT NULL UNIQUE,
    HASH  	TEXT NOT NULL UNIQUE,
	UPLOADED_AT TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
    ORDER_ID	BIGINT PRIMARY KEY,
	USER_ID 	TEXT FOREIGH KEY,
	STATUS 		TEXT NOT NULL,
	ACCURAL 	BIGINT,
	UPLOADED_AT TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS balance (
	ID			BIGINT PRIMARY KEY,
	USER_ID 	TEXT FOREIGH KEY,
	BALANCE 	BIGINT NOT NULL,
	UPLOADED_AT TIMESTAMPTZ DEFAULT now()
);


`

type PersistStorage struct {
	*pgxpool.Pool
	connString string
}

func ConnectToPersistStorage(ctx context.Context, user string, passw string, host string, port string, schema string) (*PersistStorage, error) {
	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, passw, host, port, schema)
	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect db %w", err)
	}
	if err := pgxPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping doesn't work %w", err)
	}
	return &PersistStorage{pgxPool, connString}, nil
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
