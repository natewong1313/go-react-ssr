package html

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"runtime"
)

type Params struct {
	Title      string
	MetaTags   map[string]string
	OGMetaTags map[string]string
	Links      []struct {
		Href     string
		Rel      string
		Media    string
		Hreflang string
		Type     string
		Title    string
	}
	JS         template.JS
	CSS        template.CSS
	RouteID    string
	IsDev      bool
	ServerHTML template.HTML
}

// RenderHTMLString Renders the HTML template in internal/html with the given parameters
func RenderHTMLString(params Params) []byte {
	params.IsDev = os.Getenv("APP_ENV") != "production"
	t := template.Must(template.New("").Parse(BASE_TEMPLATE))
	var output bytes.Buffer
	t.Execute(&output, params)
	return output.Bytes()
}

type ErrorParams struct {
	Error string
}

func RenderError(e error) []byte {
	t := template.Must(template.New("").Parse(ERROR_TEMPLATE))
	var output bytes.Buffer
	_, filename, line, _ := runtime.Caller(1)
	t.Execute(&output, ErrorParams{
		Error: fmt.Sprintf("%s line %d: %v", filename, line, e),
	})
	return output.Bytes()
}
