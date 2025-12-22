package main

import (
	"context"

	"github.com/Eanhain/gofermart/internal/accrual"
	flags "github.com/Eanhain/gofermart/internal/flags"
	route "github.com/Eanhain/gofermart/internal/handlers"
	logger "github.com/Eanhain/gofermart/internal/logger"
	"github.com/Eanhain/gofermart/internal/service"
	auth "github.com/Eanhain/gofermart/internal/storage/postgres/auth"
	balance "github.com/Eanhain/gofermart/internal/storage/postgres/balance"
	migr "github.com/Eanhain/gofermart/internal/storage/postgres/migration"
	orders "github.com/Eanhain/gofermart/internal/storage/postgres/orders"
	"github.com/jackc/pgx/v5/pgxpool"
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

	pgxPool, err := pgxpool.New(ctx, flagsIn.GetDBConnStr())
	if err != nil {
		log.Errorln("can't create db instance", err)
	}

	migrStore, err := migr.InitialMigration(ctx, log, pgxPool)
	if err != nil {
		log.Errorln("can't create db migration api", err)
	}
	defer migrStore.Close()

	authStore, err := auth.InitialAuth(ctx, log, pgxPool)
	if err != nil {
		log.Errorln("can't create db auth api", err)
	}
	defer authStore.Close()

	balanceStore, err := balance.InitialBalance(ctx, log, pgxPool)
	if err != nil {
		log.Errorln("can't create db balance api", err)
	}
	defer balanceStore.Close()

	ordersStore, err := orders.InitialOrders(ctx, log, pgxPool)
	if err != nil {
		log.Errorln("can't create db orders api", err)
	}
	defer ordersStore.Close()

	accrualAPI, err := accrual.InitialAccrualAPI(ctx, flagsIn.AcAddr, log)
	if err != nil {
		log.Errorln("can't create accrual API instance", err)
	}

	serv, err := service.InitialService(ctx, authStore, balanceStore, ordersStore, migrStore, accrualAPI, log)
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
