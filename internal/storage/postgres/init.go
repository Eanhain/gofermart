package postgres

import (
	"context"
	"fmt"

	domain "github.com/Eanhain/gofermart/internal/domain"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

// TODO Migration
// sql-c?

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
