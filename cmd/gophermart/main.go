package main

import (
	"context"

	"github.com/Eanhain/gofermart/internal/accrual"
	flags "github.com/Eanhain/gofermart/internal/flags"
	route "github.com/Eanhain/gofermart/internal/handlers"
	logger "github.com/Eanhain/gofermart/internal/logger"
	"github.com/Eanhain/gofermart/internal/service"
	store "github.com/Eanhain/gofermart/internal/storage/postgres"
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

	pStore, err := store.InitialPersistStorage(ctx,
		log,
		flagsIn.GetDBConnStr())

	if err != nil {
		log.Errorln("can't create db instance", err)
	}

	defer pStore.Close()

	accrualAPI, err := accrual.InitialAccrualAPI(ctx, flagsIn.AcAddr, log)
	if err != nil {
		log.Errorln("can't create accrual API instance", err)
	}

	serv, err := service.InitialService(ctx, pStore, accrualAPI, log)
	if err != nil {
		log.Errorln("can't create Service instance", err)
	}

	r := route.InitialApp(log, serv, flagsIn.GetAddr(), "supersecret")

	if err := r.CreateHandlers(ctx); err != nil {
		log.Errorln("can't create Handlers instance", err)
	}

	if err := r.StartServer(ctx); err != nil {
		log.Errorln("cannot start server", err)
	}

}
