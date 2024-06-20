package server

import (
	"net/http"

	"github.com/spf13/viper"
)

func (s *server) serveReactStaticFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(
			"/static/",
			http.FileServer(
				http.Dir(viper.GetString("reactBuildDir")+"/static"),
			),
		).ServeHTTP(w, r)
	}
}

func (s *server) serveReactIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, viper.GetString("reactBuildDir")+"/index.html")
	}
}
