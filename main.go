package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/deidelma/notedgo/noted"
)

//go:embed all:static
var static embed.FS

func main() {
	// create the server
	server := &http.Server{Addr: "0.0.0.0:5823", Handler: service()}

	// set up the server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// listen for syscall signals to interrupt or quit the program
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// shutdown signal with 30 second grace period
		shutdownCtx, x := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
			x()
		}()

		// trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("Error shutting down %v", err))
		}
		//
		// put cleanup function calls here
		//
		log.Info().Msg("Mock file cleanup initiated")
		time.Sleep(1 * time.Second)
		log.Info().Msg("Notedgo shutting down gracefully")
		serverStopCtx()
	}()

	// run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().AnErr("error running server", err)
	}

	// wait for the server to be stopped
	<-serverCtx.Done()
}

func service() http.Handler {
	noted.InitializeLogger()
	log.Info().Msg("Notedgo initialized")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// set up static files
	// use approach found in https://clavinjune.dev/en/blogs/serving-embedded-static-file-inside-subdirectory-using-go/
	// to ensure that static files are loaded correctly
	sub, err := fs.Sub(static, ".")
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(sub))
	r.Handle("/static/*", fs)

	// register handlers
	r.Get("/", noted.GetIndex)
	r.Get("/trigger_delay", noted.GetTriggerDelay)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello1"))
	})

	return r
}
