package reactbuilder

import (
	"strings"
	"text/template"
)

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

var serverTemplate = `
import React from "react";
import { renderToString } from "react-dom/server";
{{range $import := .Imports}} {{$import}} {{end}}
function render() {
	const App = require("{{ .FilePath }}").default;
	return {{ .RenderFunction }};
}
globalThis.render = render;`
var serverRenderFunction = `renderToString(<App {...props} />)`
var serverRenderFunctionWithLayout = `renderToString(<Layout><App {...props} /></Layout>)`

func GenerateServerBuildContents(imports []string, filePath string, useLayout bool) (string, error) {
	params := map[string]interface{}{
		"Imports":        imports,
		"FilePath":       filePath,
		"RenderFunction": serverRenderFunction,
	}
	if useLayout {
		params["RenderFunction"] = serverRenderFunctionWithLayout
	}
	return buildWithTemplate(serverTemplate, params)
}

var clientTemplate = `
import React from "react";
import { hydrateRoot } from "react-dom/client";
{{range $import := .Imports}} {{$import}} {{end}}
import App from "{{ .FilePath }}";
{{ .RenderFunction }}`
var clientRenderFunction = `hydrateRoot(document.getElementById("root"), <App {...props} />);`
var clientRenderFunctionWithLayout = `hydrateRoot(document.getElementById("root"), <Layout><App {...props} /></Layout>);`

func GenerateClientBuildContents(imports []string, filePath string, useLayout bool) (string, error) {
	params := map[string]interface{}{
		"Imports":        imports,
		"FilePath":       filePath,
		"RenderFunction": clientRenderFunction,
	}
	if useLayout {
		params["RenderFunction"] = clientRenderFunctionWithLayout
	}
	return buildWithTemplate(clientTemplate, params)
}
