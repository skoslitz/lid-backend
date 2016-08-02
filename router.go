package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	Handlers *Handlers
	AdminDir string
}

func NewRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// load routes from routes.go
	for _, route := range GetRoutes(config.Handlers) {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		apiRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// serve assets from static folder
	assetsFs := http.FileServer(http.Dir(config.Handlers.AssetsDir))
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", assetsFs))

	// serve preview from preview folder
	previewFs := http.FileServer(http.Dir(config.Handlers.PreviewDir))
	router.PathPrefix("/preview").Handler(http.StripPrefix("/preview/", previewFs))

	// serve admin client files (html, css, etc)
	adminFs := http.FileServer(http.Dir(config.AdminDir))
	router.PathPrefix("/").Handler(adminFs)

	return router
}
