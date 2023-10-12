package react

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"runtime"

	"github.com/natewong1313/go-react-ssr/internal/templates"
)

// Renders the HTML template in internal/templates with the given parameters
func renderHTMLString(params HTMLParams) []byte {
	params.IsDev = os.Getenv("APP_ENV") != "production"
	t := template.Must(template.New("").Parse(templates.HTML_TEMPLATE))
	var output bytes.Buffer
	t.Execute(&output, params)
	return output.Bytes()
}

// Renders the error template in internal/templates and passes the error message
func renderErrorHTMLString(err error) []byte {
	t := template.Must(template.New("").Parse(templates.ERROR_TEMPLATE))
	var output bytes.Buffer
	_, filename, line, _ := runtime.Caller(1)
	t.Execute(&output, ErrorParams{
		Error: fmt.Sprintf("%s line %d: %v", filename, line, err),
	})
	return output.Bytes()
}
