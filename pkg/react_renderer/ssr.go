package reactrenderer

import (
	"encoding/json"
	"errors"
	"fmt"
	"gossr/config"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/gin-gonic/gin"
)

// Uses esbuild to render the React page at the given file path and inserts it into an html page
func Render(c *gin.Context, filePath string, propsStruct ...interface{}) {
    // propsJson := ""
    props := "null"
    if(len(propsStruct) > 0){
        propsJSON, err := json.Marshal(propsStruct[0])
        if err != nil {
            c.JSON(500, gin.H{"error": "Invalid prop types"})
            return
        }
        props = string(propsJSON)
    }
    // Get the full path of the file
    filePath = config.Config.Web.SrcDirectory + "/" + filePath
    compiledJS, err := buildFile(filePath, props)
    if err != nil {
        c.String(500, err.Error())
    }else{
        c.HTML(http.StatusOK, "index.html", gin.H{
            "src": template.JS(compiledJS),
        })
    }
}

// func returnError(c *gin.Context, err error) {
//     c.JSON(500, gin.H{"error": err.Error()})
// }

func buildFile(filePath, props string) (string, error){
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
    newFilePath := strings.Replace(route, fileExtension, "-temporary"+fileExtension, 1)
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
