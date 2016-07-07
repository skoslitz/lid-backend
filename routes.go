package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func GetRoutes(h *Handlers) Routes {
	return Routes{
		// directories
		Route{
			"ReadDir",
			"GET",
			"/dir/{path}",
			h.ReadDir,
		},
		Route{
			"ReadDirEdition",
			"GET",
			"/dir/{path:.*}/{edition:[0-9]+}",
			h.ReadDirEdition,
		},
		Route{
			"CreateDir",
			"POST",
			"/dir/{path:.*}",
			h.CreateDir,
		},
		Route{
			"UpdateDir",
			"PUT",
			"/dir/{path:.*}",
			h.UpdateDir,
		},
		Route{
			"DeleteDir",
			"DELETE",
			"/dir/{path:.*}",
			h.DeleteDir,
		},

		// pages
		Route{
			"ReadPage",
			"GET",
			"/page/{path:.*}",
			h.ReadPage,
		},
		Route{
			"CreatePage",
			"POST",
			"/page/{path:.*}",
			h.CreatePage,
		},
		Route{
			"UpdatePage",
			"OPTIONS",
			"/page/{path:.*}",
			h.UpdatePage,
		},
		Route{
			"DeletePage",
			"DELETE",
			"/page/{path:.*}",
			h.DeletePage,
		},

		// config
		Route{
			"ReadConfig",
			"GET",
			"/config",
			h.ReadConfig,
		},
		Route{
			"UpdateConfig",
			"PUT",
			"/config",
			h.UpdateConfig,
		},

		// assets
		Route{
			"CreateAsset",
			"POST",
			"/asset/{path:.*}",
			h.CreateAsset,
		},

		// misc
		Route{
			"PublishSite",
			"POST",
			"/site/publish",
			h.PublishSite,
		},
	}
}
