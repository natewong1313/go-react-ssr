package react_renderer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	esbuildApi "github.com/evanw/esbuild/pkg/api"
)

func BuildFile(filePath, props string) (string, error){
	// Get the path of the renderer file
    newFilePath, err := makeRendererFile(filePath, props)
    if err != nil {
        return "", err
    }
    result := esbuildApi.Build(esbuildApi.BuildOptions{
        EntryPoints:       []string{newFilePath},
        Bundle:            true,
        MinifyWhitespace:  true,
        MinifyIdentifiers: true,
        MinifySyntax:      true,
        // Outfile:  ".tmp/out.js",
    })
    err = os.Remove(newFilePath)
    if err != nil {
        return "", err
    }
    if len(result.Errors) > 0 {
        return "", errors.New(result.Errors[0].Text)
    }
	// Return the compiled Javascript
    return string(result.OutputFiles[0].Contents), nil
}

// Creates a temporary file that imports the file to be rendered
func makeRendererFile(route, props string)(string, error) {
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
    ReactDOM.render(<App {...props} />, document.getElementById("root"));`,fileName, props ))
    file.Write(contents)
    file.Sync()
    return newFilePath, nil
}