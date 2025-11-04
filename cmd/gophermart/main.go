package main

import (
	"context"

	flags "github.com/Eanhain/gofermart/internal/flags"
	route "github.com/Eanhain/gofermart/internal/handlers"
	logger "github.com/Eanhain/gofermart/internal/logger"
)

func flagsInitalize(log *logger.Logger) (flags.ServerFlags, error) {
	flagInstance, err := flags.InitialFlags()
	if err != nil {
		log.Warnln(err)
	}
	flagInstance.Parse()
	log.Infoln(flagInstance)
	return flagInstance, nil
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	log := logger.InitialLogger()

	defer log.Sync()

	flagsIn, err := flagsInitalize(log)

	if err != nil {
		log.Warnln(err)
	}

	r := route.InitialApp(log, flagsIn.GetHost(), flagsIn.GetPort())

	if err := r.StartServer(ctx); err != nil {
		log.Errorln("cannot start server", err)
	}

}
