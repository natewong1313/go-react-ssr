package react

import (
	"fmt"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

func (task *RenderTask) ClientRender() {
	// Check if the client build is cached
	clientBuild, ok := getCachedClientBuild(task.FilePath)
	if !ok {
		var err error
		clientBuild, err = task.buildClientJS()
		if err != nil {
			task.ClientBuildResult <- ClientBuildResult{Error: err}
			return
		}
		setCachedClientBuild(task.FilePath, clientBuild)
	}
	js := injectProps(clientBuild.JS, task.Props)
	task.ClientBuildResult <- ClientBuildResult{JS: js, Dependencies: clientBuild.Dependencies}
}

// Build the client JS for the given React file, without props
func (task *RenderTask) buildClientJS() (ClientBuild, error) {
	globalCssImport := ""
	if tempCssFilePath != "" {
		globalCssImport = fmt.Sprintf(`import "%s";`, tempCssFilePath)
	}
	var layoutImport string
	renderStatement := `hydrateRoot(document.getElementById("root"), <App {...props} />);`
	if config.C.LayoutFile != "" {
		layoutImport = fmt.Sprintf(`import Layout from "%s";`, config.C.LayoutFile)
		renderStatement = `hydrateRoot(document.getElementById("root"), <Layout><App {...props} /></Layout>);`
	}
	// Build with esbuild
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents: fmt.Sprintf(`import * as React from "react";
			import { hydrateRoot } from "react-dom/client";
			%s
			%s
			import App from "%s";
			%s`,
				globalCssImport, layoutImport, task.FilePath, renderStatement),
			Loader:     getLoaderType(task.FilePath),
			ResolveDir: config.C.FrontendDir,
		},
		Bundle:            true,
		MinifyWhitespace:  os.Getenv("APP_ENV") == "production", // Minify in production
		MinifyIdentifiers: os.Getenv("APP_ENV") == "production",
		MinifySyntax:      os.Getenv("APP_ENV") == "production",
		Metafile:          true,
		Outdir:            "/", // This is ignored because we are using the metafile
		AssetNames:        "assets/[name]",
		Loader: map[string]esbuildApi.Loader{ // for loading images properly
			".png":  esbuildApi.LoaderFile,
			".svg":  esbuildApi.LoaderFile,
			".jpg":  esbuildApi.LoaderFile,
			".jpeg": esbuildApi.LoaderFile,
			".gif":  esbuildApi.LoaderFile,
			".bmp":  esbuildApi.LoaderFile,
		},
	})
	if len(buildResult.Errors) > 0 {
		// Return formatted error
		return ClientBuild{}, fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)
	}

	var js string
	for _, file := range buildResult.OutputFiles {
		if strings.HasSuffix(file.Path, "stdin.js") {
			js = string(file.Contents)
			break
		}
	}
	// Return the compiled build
	return ClientBuild{JS: js, Dependencies: getDependencyPathsFromMetafile(buildResult.Metafile)}, nil
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
