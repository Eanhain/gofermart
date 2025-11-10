package flags

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type ServerFlags struct {
	Host             string `short:"h" help:"Server hostname." default:"localhost" env:"HOST"`
	Port             string `short:"p" help:"Server port." default:"8080" env:"PORT"`
	PostgresHost     string `long:"phost" help:"Postgres connection host." default:"127.0.0.1" env:"POSTGRES_HOST"`
	PostgresPort     string `long:"pport" help:"Postgres connection port." default:"5432" env:"POSTGRES_PORT"`
	PostgresUser     string `long:"puser" help:"Postgres connection user." default:"db_user" env:"POSTGRES_USER"`
	PostgresPassword string `long:"ppass" help:"Postgres connection password." default:"s3cret" env:"POSTGRES_PASSWORD"`
	PostgresSchema   string `long:"pschema" help:"Postgres connection schema." default:"gofemart" env:"POSTGRES_SCHEMA"`
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
	return sf.Host
}

func (sf ServerFlags) GetPort() string {
	return sf.Port
}

func (sf ServerFlags) GetAddr() string {
	return fmt.Sprintf("%v:%v", sf.Host, sf.Port)
}
