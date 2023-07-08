package noted

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/url"
	"net/http"
	"text/template"

	"github.com/rs/zerolog/log"
)

var (
	//go:embed all:templates
	res   embed.FS
	pages = map[string]string{
		"/": "templates/index.gohtml",
	}
	//go:embed all:static
	static embed.FS
)

func InitServer() {
	log.Info().Msg("initializing server")
	RegisterFileServers()
	RegisterHandlers()
	err := http.ListenAndServe(":5823", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("server closed")
	} else if err != nil {
		log.Warn().Msg(fmt.Sprintf("unable to start server: %s", err))
		panic(1)
	}
}

func RegisterFileServers(){
	fs := http.FileServer(http.FS(static))
	http.Handle("/static/", fs)
	log.Info().Msg("static server initialized")
	log.Debug().Msg("file server installing")
	http.FileServer(http.FS(res))
	log.Info().Msg("file server initialized")
}

func RegisterHandlers() {
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/api/hello", getGreeting)
	http.HandleFunc("/trigger_delay", getTriggerDelay)
	log.Info().Msg("handlers installed")
}

func getIndex(writer http.ResponseWriter, request *http.Request) {
	page, ok := pages[request.URL.Path]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	log.Debug().Msg(fmt.Sprintf("about to process template: %s", page))
	tpl, err := template.ParseFS(res, page)
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("page %s not found in pages cache...", request.RequestURI))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	data := map[string]interface{}{
		"name": "David",
	}
	if err := tpl.Execute(writer, data); err != nil {
		log.Info().Msg(fmt.Sprintf("error encountered in executing template: %s", err))
		return
	}
}

func getTriggerDelay(w http.ResponseWriter, r *http.Request){
	log.Info().Msg("got /trigger_delay")
	parms, err := url.ParseQuery(r.RequestURI)
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("Unable to parse URI: %s ", r.RequestURI))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = io.WriteString(w, fmt.Sprintf("<h3>%s</h3>", parms.Get("/trigger_delay?q")))
	if err != nil {
		log.Warn().Msg("Unable to write")
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
}
