package pStorage

import (
	"context"
	"testing"

	dto "github.com/Eanhain/gofermart/internal/api"
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
		log     Logger
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

func TestGetUserFromDB(t *testing.T) {
	logger := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     Logger
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

			defer psInst.Close()
			rows := pgxmock.NewRows([]string{"username", "hash"}).AddRow("test", "hash2")
			mock.ExpectQuery(selectUser).WithArgs(tt.user.Login).WillReturnRows(rows)
			if _, _, err := psInst.GetUserFromDB(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("GetUserFromDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	logger := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     Logger
		users   dto.UserArray
		wantErr error
	}{
		{
			name: "OK",
			ctx:  context.Background(),
			log:  logger,
			users: dto.UserArray{
				dto.User{
					Login: "test",
					Hash:  "hash1",
				},
				dto.User{
					Login: "test2",
					Hash:  "hash2",
				},
				dto.User{
					Login: "test3",
					Hash:  "hash3",
				}},
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
			for _, user := range tt.users {
				mockBatch.ExpectExec(InsertUser.Name).WithArgs(user.Login, user.Hash).WillReturnResult(connCommand)
			}
			mock.ExpectCommit()
			defer psInst.Close()

			if err := psInst.RegisterUser(tt.ctx, tt.users); err != tt.wantErr {
				t.Fatalf("GetUserFromDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
