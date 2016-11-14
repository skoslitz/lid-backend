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
			"CreateDir",
			"POST",
			"/dir/{path:.*}",
			h.CreateDir,
		},
		Route{
			"UpdateDir",
			"OPTIONS",
			"/dir/{path:.*}",
			h.UpdateDir,
		},
		Route{
			"DeleteDir",
			"DELETE",
			"/dir/{path:.*}",
			h.DeleteDir,
		},
		Route{
			"ReadRegionRelationships",
			"GET",
			"/regionen/{id}/{type}",
			h.ReadRegionRelationships,
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
			"OPTIONS",
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
		Route{
			"UpdateAsset",
			"OPTIONS",
			"/asset/{path:.*}",
			h.UpdateAsset,
		},
		Route{
			"DeleteAsset",
			"DELETE",
			"/asset/{path:.*}",
			h.DeleteAsset,
		},

		// misc
		Route{
			"PublishSite",
			"POST",
			"/site/publish",
			h.PublishSite,
		},
		Route{
			"PreviewSite",
			"POST",
			"/site/preview",
			h.PreviewSite,
		},
	}
}
