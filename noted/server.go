package noted

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/rs/zerolog/log"
)
var (
	//go:embed static/templates/index.gohtml
	res   embed.FS
	pages = map[string]string{
		"/": "static/templates/index.gohtml",
	}
)

func InitServer() {
	log.Info().Msg("initializing server")
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/api/hello", getGreeting)
	log.Info().Msg("handlers installed")
	log.Info().Msg("file server installing")
	http.FileServer(http.FS(res))
	log.Info().Msg("server initialized")

	err := http.ListenAndServe(":5555", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("server closed")
	} else if err != nil {
		log.Fatal().Msg(fmt.Sprintf("unable to start server: %s", err))
	}
}

func getIndex(writer http.ResponseWriter, request *http.Request) {
	page, ok := pages[request.URL.Path]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	log.Info().Msg(fmt.Sprintf("about to process template: %s", page))
	tpl, err := template.ParseFS(res, page)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("page %s not found in pages cache...", request.RequestURI))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	data := map[string]interface{}{
		"userAgent": request.UserAgent(),
	}
	if err := tpl.Execute(writer, data); err != nil {
		log.Info().Msg(fmt.Sprintf("error encountered in executing template: %s", err))
		return
	}
}

func getHello(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("got /hello request")
	_, err := io.WriteString(w, "hello!")
	if err != nil {
		log.Info().Msg("unable to handle writing in getHello")
		return
	}
}

func getGreeting(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("got /api/hello request")
	_, err := io.WriteString(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<h1>Hello from api request</h1>
		<br/>
		<h2> Explanation </h2>
		<p>Strange to think that this is an api not a web page</p>
	</body>
	</html>
	`)
	if err != nil {
		log.Info().Msg("unable to handle writing in getGreeting")
		return
	}
}
