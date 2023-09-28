package react_renderer

import "html/template"

// Build stores the compiled JS and CSS for a React file
type Build struct {
	CompiledJS  string
	CompiledCSS string
}

// Config stores the configuration for a React file to be rendered
type Config struct {
	File     string
	Title    string
	MetaTags map[string]string
	Links    []LinkElement
	Props    interface{}
}

// LinkElement stores the attributes for a link element
type LinkElement struct {
	Href     string
	Rel      string
	Media    string
	Hreflang string
	Type     string
	Title    string
}

// HTMLParams stores the parameters for the html template
type HTMLParams struct {
	Title      string
	MetaTags   map[string]string
	OGMetaTags map[string]string
	Links      []LinkElement
	JS         template.JS
	CSS        template.CSS
	RouteID    string
	IsDev      bool
	ServerHTML template.HTML
}

// ErrorParams stores the parameters for the error template
type ErrorParams struct {
	Error string
}

type ClientBuildResult struct {
	JS           string
	Dependencies []string
	Error        error
}

type ServerBuildResult struct {
	HTML  string
	CSS   string
	Error error
}
