package flags

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type ServerFlags struct {
	Host string `short:"h" help:"Server hostname." default:"localhost" env:"HOST"`
	Port string `short:"p" help:"Server port." default:"8080" env:"PORT"`
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
