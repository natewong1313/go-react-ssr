package react_renderer

import "html/template"

type Config struct {
	File     string
	Title    string
	MetaTags map[string]string
	Links    []LinkElement
	Props    interface{}
}

type LinkElement struct {
	Href     string
	Rel      string
	Media    string
	Hreflang string
	Type     string
	Title    string
}

type HTMLParams struct {
	Title      string
	MetaTags   map[string]string
	OGMetaTags map[string]string
	Links      []LinkElement
	JS         template.JS
	CSS        template.CSS
	RouteID    string
	IsDev      bool
}

type ErrorParams struct {
	Error string
}
