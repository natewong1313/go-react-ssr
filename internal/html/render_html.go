package html

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"runtime"
	"strings"
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
	params.OGMetaTags = getOGMetaTags(params.MetaTags)
	params.MetaTags = getMetaTags(params.MetaTags)
	t := template.Must(template.New("").Parse(BaseTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, params)
	if err != nil {
		return RenderError(err, params.RouteID)
	}
	return output.Bytes()
}

func getMetaTags(metaTags map[string]string) map[string]string {
	newMetaTags := make(map[string]string)
	for key, value := range metaTags {
		if !strings.HasPrefix(key, "og:") {
			newMetaTags[key] = value
		}
	}
	return newMetaTags
}

func getOGMetaTags(metaTags map[string]string) map[string]string {
	newMetaTags := make(map[string]string)
	for key, value := range metaTags {
		if strings.HasPrefix(key, "og:") {
			newMetaTags[key] = value
		}
	}
	return newMetaTags
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
