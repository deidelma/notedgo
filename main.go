
package main

import (
	"net/http"
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/deidelma/notedgo/noted"

)

//go:embed all:noted/templates
// var res embed.FS

//go:embed all:noted/static
var static embed.FS

func main(){
	noted.InitializeLogger()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	fs := http.FileServer(http.FS(static))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", noted.GetIndex)
	r.Get("/trigger_delay", noted.GetTriggerDelay)
	r.Get("/hello", func( w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello1"))
	})
	http.ListenAndServe(":5823", r)
}