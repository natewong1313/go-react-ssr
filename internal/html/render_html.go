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
	t := template.Must(template.New("").Parse(BaseTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, params)
	if err != nil {
		return RenderError(err, params.RouteID)
	}
	return output.Bytes()
}

type ErrorParams struct {
	Error   string
	RouteID string
	IsDev   bool
}

// RenderError Renders the error template with the given error
func RenderError(e error, routeID string) []byte {
	t := template.Must(template.New("").Parse(ErrorTemplate))
	var output bytes.Buffer
	_, filename, line, _ := runtime.Caller(1)
	t.Execute(&output, ErrorParams{
		Error:   fmt.Sprintf("%s line %d: %v", filename, line, e),
		RouteID: routeID,
		IsDev:   os.Getenv("APP_ENV") != "production",
	})
	return output.Bytes()
}
