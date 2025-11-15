package pStorage

import (
	"context"
	"testing"

	dto "github.com/Eanhain/gofermart/internal/api"
	logger "github.com/Eanhain/gofermart/internal/logger"
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
			psInst := PersistStorage{mock}
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
			psInst := PersistStorage{mock}

			defer psInst.Close()
			rows := pgxmock.NewRows([]string{"username", "hash"}).AddRow("test", "hash2")
			mock.ExpectQuery(selectUser).WithArgs(tt.user.Login).WillReturnRows(rows)
			if unt, trust, err := psInst.GetUserFromDB(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("GetUserFromDB() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log("Input user:", unt, "Output user:", trust)
			}
		})
	}
}
