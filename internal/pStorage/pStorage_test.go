package pStorage

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	dto "github.com/Eanhain/gofermart/internal/api"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

func TestConnectToPersistStorage(t *testing.T) {
	type args struct {
		ctx    context.Context
		log    Logger
		user   string
		passw  string
		host   string
		port   string
		schema string
	}
	pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v", "user1", "passw1", "host1", "port1", "schema1"))
	tests := []struct {
		name    string
		args    args
		want    *PersistStorage
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnectToPersistStorage(tt.args.ctx, tt.args.log, tt.args.user, tt.args.passw, tt.args.host, tt.args.port, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ConnectToPersistStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConnectToPersistStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersistStorage_InitSchema(t *testing.T) {
	type fields struct {
		Pool *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		log Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PersistStorage{
				Pool: tt.fields.Pool,
			}
			if err := ps.InitSchema(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("PersistStorage.InitSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersistStorage_RegisterUser(t *testing.T) {
	type fields struct {
		Pool *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		users dto.UserArray
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PersistStorage{
				Pool: tt.fields.Pool,
			}
			if err := ps.RegisterUser(tt.args.ctx, tt.args.users); (err != nil) != tt.wantErr {
				t.Errorf("PersistStorage.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersistStorage_CheckUserPermissions(t *testing.T) {
	type fields struct {
		Pool *pgxpool.Pool
	}
	type args struct {
		ctx           context.Context
		untrustedUser dto.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PersistStorage{
				Pool: tt.fields.Pool,
			}
			got, err := ps.CheckUserPermissions(tt.args.ctx, tt.args.untrustedUser)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersistStorage.CheckUserPermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PersistStorage.CheckUserPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}
