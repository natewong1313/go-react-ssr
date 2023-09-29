package react_renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buger/jsonparser"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

func makeClientBuild(reactFilePath, props string, clientBuildResult chan<- ClientBuildResult) {
	// Check if the client build is cached
	clientBuild, ok := getCachedClientBuild(reactFilePath)
	if !ok {
		var err error
		clientBuild, err = buildClientJS(reactFilePath)
		if err != nil {
			clientBuildResult <- ClientBuildResult{Error: err}
			return
		}
		setCachedClientBuild(reactFilePath, clientBuild)
	}
	js := injectProps(clientBuild.JS, props)
	clientBuildResult <- ClientBuildResult{JS: js, Dependencies: clientBuild.Dependencies}
}

// Build the client JS for the given React file, without props
func buildClientJS(reactFilePath string) (ClientBuild, error) {
	defer utils.Timer("buildForClient")()
	globalCssImport := ""
	if tempCssFilePath != "" {
		globalCssImport = fmt.Sprintf(`import "%s";`, tempCssFilePath)
	}
	// Build with esbuild
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents: fmt.Sprintf(`import * as React from "react";
			import * as ReactDOM from "react-dom";
			%s
			import App from "./%s";
			ReactDOM.hydrate(<App {...window.props} />, document.getElementById("root"));`,
				globalCssImport, filepath.ToSlash(filepath.Base(reactFilePath))),
			Loader:     getLoaderType(reactFilePath),
			ResolveDir: config.C.FrontendDir,
		},
		Bundle:            true,
		MinifyWhitespace:  os.Getenv("APP_ENV") == "production", // Minify in production
		MinifyIdentifiers: os.Getenv("APP_ENV") == "production",
		MinifySyntax:      os.Getenv("APP_ENV") == "production",
		Metafile:          true,
		Outdir:            "/", // This is ignored because we are using the metafile
		Loader: map[string]esbuildApi.Loader{ // for loading images properly
			".png":  esbuildApi.LoaderDataURL,
			".svg":  esbuildApi.LoaderDataURL,
			".jpg":  esbuildApi.LoaderDataURL,
			".jpeg": esbuildApi.LoaderDataURL,
			".gif":  esbuildApi.LoaderDataURL,
			".bmp":  esbuildApi.LoaderDataURL,
		},
	})
	if len(buildResult.Errors) > 0 {
		// Return formatted error
		return ClientBuild{}, fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)
	}
	// Return the compiled build
	return ClientBuild{JS: string(buildResult.OutputFiles[0].Contents), Dependencies: getDependencyPathsFromMetafile(buildResult.Metafile)}, nil
}

// Inject the props into the compiled JS
func injectProps(compiledJS, props string) string {
	return fmt.Sprintf(`window.props = %s; %s`, props, compiledJS)
}

// Get the esbuild loader type for the React file, depending on the file extension
func getLoaderType(reactFilePath string) esbuildApi.Loader {
	loader := esbuildApi.LoaderTSX
	if strings.HasSuffix(reactFilePath, ".ts") {
		loader = esbuildApi.LoaderTS
	}
	if strings.HasSuffix(reactFilePath, ".js") {
		loader = esbuildApi.LoaderJS
	}
	if strings.HasSuffix(reactFilePath, ".jsx") {
		loader = esbuildApi.LoaderJSX
	}
	return loader
}

// Parse dependencies from esbuild metafile
func getDependencyPathsFromMetafile(metafile string) []string {
	var dependencyPaths []string
	// Parse the metafile and get the paths of the dependencies
	// Ignore dependencies in node_modules
	jsonparser.ObjectEach([]byte(metafile), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if !strings.Contains(string(key), "/node_modules/") {
			dependencyPaths = append(dependencyPaths, utils.GetFullFilePath(string(key)))
		}
		return nil
	}, "inputs")
	return dependencyPaths
}
