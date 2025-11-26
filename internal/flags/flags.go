package flags

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type ServerFlags struct {
	Addr   string `short:"a" help:"Server addr" default:"localhost:8081" env:"RUN_ADDRESS"`
	DdAddr string `short:"d" help:"Postgres connection addr" default:"postgres://db_user:s3cret@127.0.0.1/gofemart" env:"DATABASE_URI"`
	AcAddr string `short:"r" help:"Accrual connection addr" default:"localhost:8080" env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (sf *ServerFlags) Parse() {
	kong.Parse(sf)
}

func InitialFlags() (ServerFlags, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return ServerFlags{}, fmt.Errorf("couldn't import flags %w", err)
	}
	return ServerFlags{}, nil
}

func (sf ServerFlags) GetHost() string {
	splitStr := strings.Split(sf.Addr, ":")
	return splitStr[0]
}

func (sf ServerFlags) GetPort() string {
	splitStr := strings.Split(sf.Addr, ":")
	return splitStr[1]
}

func (sf ServerFlags) GetAddr() string {
	return sf.Addr
}

func (sf ServerFlags) GetDBConnStr() string {
	return sf.DdAddr
}
