package noted

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/rs/zerolog/log"
)

var (
	//go:embed all:templates
	res   embed.FS
)


func GetIndex(writer http.ResponseWriter, request *http.Request) {
	page := "templates/index.gohtml"
	log.Debug().Msg(fmt.Sprintf("about to process template: %s", page))
	tpl, err := template.ParseFS(res, page)
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("page %s not found in pages cache...", request.RequestURI))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"name": "David",
	}
	if err := tpl.Execute(writer, data); err != nil {
		log.Info().Msg(fmt.Sprintf("error encountered in executing template: %s", err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetTriggerDelay(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Param q from r.URL.Query:%s\n", r.URL.Query().Get("bo"))
	q := r.URL.Query()
	_, err := io.WriteString(w, fmt.Sprintf("<h3>%s</h3>", q.Get("q")))
	if err != nil {
		log.Warn().Msg("Unable to write after trigger")
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
}
