package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/deidelma/notedgo/noted"
)

////go:embed all:templates
// var res embed.FS

//go:embed all:static
var static embed.FS

func main(){
	noted.InitializeLogger()
	log.Info().Msg("Notedgo initialized")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	sub, err := fs.Sub(static, ".");
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(sub))
	r.Handle("/static/*", fs)
	// fs := http.FileServer(http.Dir("static"))
	// r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.Get("/", noted.GetIndex)
	r.Get("/trigger_delay", noted.GetTriggerDelay)
	r.Get("/hello", func( w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello1"))
	})
	http.ListenAndServe(":5823", r)
}