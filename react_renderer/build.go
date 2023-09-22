package react_renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buger/jsonparser"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
)

type Build struct {
	CompiledJS  string
	CompiledCSS string
}

func buildReactFile(routeID, reactFilePath, props string) (Build, []string, error) {
	var build Build
	// Build with esbuild
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		// Build contents from a string rather than file
		Stdin: &esbuildApi.StdinOptions{
			Contents:   getBuildContents(reactFilePath, props),
			Loader:     esbuildApi.LoaderTSX,
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
		return build, nil, fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)
	}
	// First output file is the react build
	build.CompiledJS = string(buildResult.OutputFiles[0].Contents)
	// Check for css file
	for _, file := range buildResult.OutputFiles {
		if strings.HasSuffix(string(file.Path), ".css") {
			build.CompiledCSS = string(file.Contents)
			break
		}
	}
	// Return the compiled build
	return build, getDependencyPathsFromMetafile(buildResult.Metafile), nil
}

// Parse dependencies from esbuild metafile
func getDependencyPathsFromMetafile(metafile string) []string {
	var dependencyPaths []string
	jsonparser.ObjectEach([]byte(metafile), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if !strings.Contains(string(key), "/node_modules/") {
			dependencyPaths = append(dependencyPaths, getFullFilePath(string(key)))
		}
		return nil
	}, "inputs")
	return dependencyPaths
}

// This will imports the desired react file and sets props
func getBuildContents(reactFilePath, props string) string {
	return fmt.Sprintf(`import * as React from "react";
	import * as ReactDOM from "react-dom";
	import App from "./%s";
	const props = %s
	ReactDOM.render(<App {...props} />, document.getElementById("root"));`,
		filepath.Base(reactFilePath), props)
}
