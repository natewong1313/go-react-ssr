package react_renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	esbuildApi "github.com/evanw/esbuild/pkg/api"
)

func BuildFile(filePath, props string) (CachedBuild, string, error) {
	var cachedBuild CachedBuild
	// Get the path of the renderer file
	newFilePath, err := makeRendererFile(filePath, props)
	if err != nil {
		return cachedBuild, "", err
	}
	// Get temporary directory
	osCacheDir, _ := os.UserCacheDir()
	outDir := filepath.Join(osCacheDir, "gossr")
	// Build the file with esbuild
	result := esbuildApi.Build(esbuildApi.BuildOptions{
		EntryPoints:       []string{newFilePath},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Metafile:          true,
		Outdir:            outDir,
		Loader: map[string]esbuildApi.Loader{
			".png":  esbuildApi.LoaderDataURL,
			".svg":  esbuildApi.LoaderDataURL,
			".jpg":  esbuildApi.LoaderDataURL,
			".jpeg": esbuildApi.LoaderDataURL,
			".gif":  esbuildApi.LoaderDataURL,
			".bmp":  esbuildApi.LoaderDataURL,
		},
	})
	// Remove renderer file
	err = os.Remove(newFilePath)
	if err != nil {
		return cachedBuild, "", err
	}
	if len(result.Errors) > 0 {
		return cachedBuild, "", fmt.Errorf("%s <br>in %s <br>at %s", result.Errors[0].Text, result.Errors[0].Location.File, result.Errors[0].Location.LineText)
	}
	cachedBuild.CompiledJS = string(result.OutputFiles[0].Contents)
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(string(file.Path), ".css") {
			cachedBuild.CompiledCSS = string(file.Contents)
			break
		}
	}
	// Return the compiled build
	return cachedBuild, result.Metafile, nil
}

// Creates a temporary file that imports the file to be rendered
func makeRendererFile(route, props string) (string, error) {
	fileExtension := filepath.Ext(route)
	fileName := filepath.Base(route)
	newFilePath := strings.Replace(route, fileExtension, "-gossr-temporary"+fileExtension, 1)
	// Create the file if it doesn't exist
	file, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer file.Close()
	contents := []byte(fmt.Sprintf(`import * as React from "react";
    import * as ReactDOM from "react-dom";
    import App from "./%s";
    const props = %s
    ReactDOM.render(<App {...props} />, document.getElementById("root"));`, fileName, props))
	file.Write(contents)
	file.Sync()
	return newFilePath, nil
}
