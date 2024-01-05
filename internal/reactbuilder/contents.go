package reactbuilder

import (
	"strings"
	"text/template"
)

var baseTemplate = `
import React from "react";
{{range $import := .Imports}}{{$import}} {{end}}
import App from "{{ .FilePath }}";
{{ if .SuppressConsoleLog }}console.log = () => {};{{ end }}
{{ .RenderFunction }}`
var serverRenderFunction = `globalThis.out = renderToString(<App {...props} />);`
var serverRenderFunctionWithLayout = `globalThis.out = renderToString(<Layout><App {...props} /></Layout>);`
var clientRenderFunction = `hydrateRoot(document.getElementById("root"), <App {...props} />);`
var clientRenderFunctionWithLayout = `hydrateRoot(document.getElementById("root"), <Layout><App {...props} /></Layout>);`

func buildWithTemplate(buildTemplate string, params map[string]interface{}) (string, error) {
	templ, err := template.New("buildTemplate").Parse(buildTemplate)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	err = templ.Execute(&out, params)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func GenerateServerBuildContents(imports []string, filePath string, useLayout bool) (string, error) {
	imports = append(imports, `import { renderToString } from "react-dom/server.browser";`)
	params := map[string]interface{}{
		"Imports":            imports,
		"FilePath":           filePath,
		"RenderFunction":     serverRenderFunction,
		"SuppressConsoleLog": true,
	}
	if useLayout {
		params["RenderFunction"] = serverRenderFunctionWithLayout
	}
	return buildWithTemplate(baseTemplate, params)
}

func GenerateClientBuildContents(imports []string, filePath string, useLayout bool) (string, error) {
	imports = append(imports, `import { hydrateRoot } from "react-dom/client";`)
	params := map[string]interface{}{
		"Imports":        imports,
		"FilePath":       filePath,
		"RenderFunction": clientRenderFunction,
	}
	if useLayout {
		params["RenderFunction"] = clientRenderFunctionWithLayout
	}
	return buildWithTemplate(baseTemplate, params)
}
