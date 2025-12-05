package postgres

import (
	"context"
	"testing"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	logger "github.com/Eanhain/gofermart/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
)

// func Test
func TestConnectToPersistStorage(t *testing.T) {
	logger := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     domain.Logger
		wantErr error
	}{
		{
			name:    "OK",
			ctx:     context.Background(),
			log:     logger,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Cannot create mock %v", err)
			}
			psInst := PersistStorage{mock, tt.log}
			defer psInst.Close()
			mock.ExpectBegin()
			ddls := []string{ddlUsers, ddlOrders, ddlBalance}
			for _, ddl := range ddls {
				mock.ExpectExec(ddl).WillReturnResult(pgconn.NewCommandTag("CREATE TABLE"))
			}
			mock.ExpectCommit()
			if err = psInst.InitSchema(tt.ctx, tt.log); err != tt.wantErr {
				t.Fatalf("ConnectToPersistStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckUser(t *testing.T) {
	logger := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     domain.Logger
		user    dto.UserInput
		wantErr error
	}{
		{
			name: "OK",
			ctx:  context.Background(),
			log:  logger,
			user: dto.UserInput{
				Login:    "test",
				Password: "passw",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Cannot create mock %v", err)
			}
			psInst := PersistStorage{mock, tt.log}

			defer psInst.Close()
			rows := pgxmock.NewRows([]string{"username", "hash"}).AddRow("test", "hash2")
			mock.ExpectQuery(selectUser).WithArgs(tt.user.Login).WillReturnRows(rows)
			if _, err := psInst.CheckUser(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("CheckUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	logger := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     domain.Logger
		user    dto.User
		wantErr error
	}{
		{
			name: "OK",
			ctx:  context.Background(),
			log:  logger,
			user: dto.User{
				Login: "test",
				Hash:  "hash1",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Cannot create mock %v", err)
			}
			psInst := PersistStorage{mock, tt.log}
			mock.ExpectBeginTx(pgx.TxOptions{})
			mock.ExpectPrepare(InsertUser.Name, InsertUser.DML)
			connCommand := pgconn.NewCommandTag("INSERT")
			mockBatch := mock.ExpectBatch()
			mockBatch.ExpectExec(InsertUser.Name).WithArgs(tt.user.Login, tt.user.Hash).WillReturnResult(connCommand)
			mock.ExpectCommit()
			defer psInst.Close()

			if err := psInst.RegisterUser(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("CheckUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
