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

func buildForClient(reactFilePath, props string, c chan<- ClientBuildResult) {
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
			const props = %s
			ReactDOM.hydrate(<App {...props} />, document.getElementById("root"));`,
				globalCssImport, filepath.Base(reactFilePath), props),
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
		c <- ClientBuildResult{Error: fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)}
		return
	}
	compiledJS := string(buildResult.OutputFiles[0].Contents)
	// Return the compiled build
	c <- ClientBuildResult{JS: compiledJS, Dependencies: getDependencyPathsFromMetafile(buildResult.Metafile)}
}

// Get the esbuild loader type for the react file, dependeing on the file extension
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
