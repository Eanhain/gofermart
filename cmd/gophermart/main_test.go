package main

import (
	"reflect"
	"testing"

	flags "github.com/Eanhain/gofermart/internal/flags"
	logger "github.com/Eanhain/gofermart/internal/logger"
)

func Test_flagsInitalize(t *testing.T) {
	type args struct {
		log *logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    flags.ServerFlags
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flagsInitalize(tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Fatalf("flagsInitalize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flagsInitalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
