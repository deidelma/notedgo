package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/deidelma/notedgo/noted"

	"github.com/rs/zerolog/log"
)


// put finalization logic here -- executed just before exit.
func cleanup() {
	log.Info().Msg("exiting program")
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
	noted.InitializeLogger()
	noted.LoadConfiguration(os.Args[0], os.Args[1:])
	noted.InitServer()
}
