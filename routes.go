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

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"config",
		"GET",
		"/config",
		ReadConfig,
	},
	Route{
		"ContentIndex",
		"GET",
		"/content",
		ReadContentIndex,
	},
	Route{
		"ContentType",
		"GET",
		"/content/{contentType}",
		ReadContentType,
	},
	Route{
		"ContentTypeFile",
		"GET",
		"/content/{contentType}/{fileName}",
		ReadContentTypeFile,
	},
	Route{
		"SplitContentTypeFile",
		"GET",
		"/content/{contentType}/{fileName}/{filePart}",
		ReadSplitContentTypeFile,
	},
}
