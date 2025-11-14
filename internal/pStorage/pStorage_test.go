package pStorage

import (
	"context"
	"testing"

	logger "github.com/Eanhain/gofermart/internal/logger"
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
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("Cannot create mock %v", err)
			}
			psInst := PersistStorage{mock}
			defer psInst.Close()
			mock.ExpectBegin()
			ddls := []string{ddlUsers, ddlOrders, ddlBalance}
			for _, ddl := range ddls {
				mock.ExpectExec(ddl)
			}
			mock.ExpectCommit()
			if err = psInst.InitSchema(tt.ctx, tt.log); err != tt.wantErr {
				t.Fatalf("ConnectToPersistStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestPersistStorage_InitSchema(t *testing.T) {
// 	type fields struct {
// 		Pool *pgxpool.Pool
// 	}
// 	type args struct {
// 		ctx context.Context
// 		log Logger
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ps := &PersistStorage{
// 				Pool: tt.fields.Pool,
// 			}
// 			if err := ps.InitSchema(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
// 				t.Errorf("PersistStorage.InitSchema() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPersistStorage_RegisterUser(t *testing.T) {
// 	type fields struct {
// 		Pool *pgxpool.Pool
// 	}
// 	type args struct {
// 		ctx   context.Context
// 		users dto.UserArray
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ps := &PersistStorage{
// 				Pool: tt.fields.Pool,
// 			}
// 			if err := ps.RegisterUser(tt.args.ctx, tt.args.users); (err != nil) != tt.wantErr {
// 				t.Errorf("PersistStorage.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPersistStorage_CheckUserPermissions(t *testing.T) {
// 	type fields struct {
// 		Pool *pgxpool.Pool
// 	}
// 	type args struct {
// 		ctx           context.Context
// 		untrustedUser dto.User
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    dto.User
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ps := &PersistStorage{
// 				Pool: tt.fields.Pool,
// 			}
// 			got, err := ps.CheckUserPermissions(tt.args.ctx, tt.args.untrustedUser)
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("PersistStorage.CheckUserPermissions() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr {
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("PersistStorage.CheckUserPermissions() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
