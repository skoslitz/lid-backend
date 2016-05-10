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
}
