package main

import (
	flags "github.com/Eanhain/gofermart/internal/flags"
	logger "github.com/Eanhain/gofermart/internal/logger"
)

func flagsInitalize(log *logger.Logger) error {
	flagInstance, err := flags.InitialFlags()
	if err != nil {
		log.Warnln(err)
	}
	flagInstance.Parse()
	log.Infoln(flagInstance)
	return nil
}

func main() {
	log := logger.InitialLogger()

	defer log.Sync()

	if err := flagsInitalize(log); err != nil {
		log.Warnln(err)
	}

}
