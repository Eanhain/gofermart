package postgres

import (
	"context"
	"testing"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	logger "github.com/Eanhain/gofermart/internal/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
)

func TestConnectToPersistStorage(t *testing.T) {
	log := logger.InitialLogger()
	tests := []struct {
		name    string
		ctx     context.Context
		log     domain.Logger
		wantErr error
	}{
		{
			name:    "OK",
			ctx:     context.Background(),
			log:     log,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Cannot create mock %v", err)
			}

			psInst := PersistStorage{
				PgxIface: mock,
				log:      tt.log,
				psql:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
			}

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
	log := logger.InitialLogger()
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
			log:  log,
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

			psInst := PersistStorage{
				PgxIface: mock,
				log:      tt.log,
				psql:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
			}
			defer psInst.Close()

			rows := pgxmock.NewRows([]string{"username", "hash"}).AddRow("test", "hash2")

			expectedSQL := "SELECT username, hash FROM users WHERE username = $1"

			mock.ExpectQuery(expectedSQL).WithArgs(tt.user.Login).WillReturnRows(rows)

			if _, err := psInst.CheckUser(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("CheckUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	log := logger.InitialLogger()
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
			log:  log,
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

			psInst := PersistStorage{
				PgxIface: mock,
				log:      tt.log,
				psql:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
			}
			defer psInst.Close()

			connCommand := pgconn.NewCommandTag("INSERT 0 1")
			expectedSQL := "INSERT INTO users (username, hash) VALUES ($1, $2)"

			mock.ExpectExec(expectedSQL).
				WithArgs(tt.user.Login, tt.user.Hash).
				WillReturnResult(connCommand)

			if err := psInst.RegisterUser(tt.ctx, tt.user); err != tt.wantErr {
				t.Fatalf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
