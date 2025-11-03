package main

import (
	"log"

	flags "github.com/Eanhain/gofermart/internal/flags"
)

func flagsInitalize() error {
	flagInstance, err := flags.InitialFlags()
	if err != nil {
		log.Println(err)
	}
	flagInstance.Parse()
	log.Println(flagInstance)
	return nil
}

func main() {
	if err := flagsInitalize(); err != nil {
		log.Println(err)
	}
}
