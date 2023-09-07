package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/gin-gonic/gin"
)

// result := api.Build(api.BuildOptions{
//     EntryPoints:       []string{"./frontend-2/src/App.tsx"},
//     Bundle:            true,
//     MinifyWhitespace:  true,
//     MinifyIdentifiers: true,
//     MinifySyntax:      true,
//     // Engines: []api.Engine{
//     //   {api.EngineChrome, "58"},
//     //   {api.EngineFirefox, "57"},
//     //   {api.EngineSafari, "11"},
//     //   {api.EngineEdge, "16"},
//     // },
//     Write: true,

// })
// if len(result.Errors) > 0 {
//    fmt.Println("Err")
// }
// fmt.Println(result.OutputFiles[0].Contents)
// result.OutputFiles

// func createTmpDir() error {
//     tmpPath := ".tmp"
//     if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
//         err := os.MkdirAll(tmpPath, 0755)
//         if err != nil {
//             return err
//         }
//     }
//     return nil
// }
// var TEMP_PATH = "./.tmp/"

// func createRendererFile(sourceFileName string) error {
//     file, err := os.OpenFile(TEMP_PATH+"go-ssr-main.jsx", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
//     if err != nil {
//         return err
//     }
//     defer file.Close()
//     contents := []byte(fmt.Sprintf(`import * as React from "react";
//     import * as ReactDOM from "react-dom";
//     import App from "./%s";

//     ReactDOM.render(<App />, document.getElementById("root"));`,sourceFileName ))
//     file.Write(contents)
//     file.Sync()
//     return nil
// }

// func addRendererToFile(fileName string) error {
//     err := createTmpDir()
//     if err != nil {
//         panic(err)
//     }
//     r, err := os.Open("./frontend-2/src/"+fileName)
//     if err != nil {
//         panic(err)
//     }
//     defer r.Close()
//     w, err := os.Create(TEMP_PATH+fileName)
//     if err != nil {
//         panic(err)
//     }
//     defer w.Close()
//     w.ReadFrom(r)
//     createRendererFile(fileName)
//     return nil
// }

func makeRendererFile(route string)(string, error) {
    fileExtension := filepath.Ext(route)
    fileName := filepath.Base(route)
    newFilePath := strings.Replace(route, fileExtension, "-temporary"+fileExtension, 1)

    file, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        return "", err
    }
    defer file.Close()
    contents := []byte(fmt.Sprintf(`import * as React from "react";
    import * as ReactDOM from "react-dom";
    import App from "./%s";
    const props = {initialCount: 10}
    ReactDOM.render(<App {...props} />, document.getElementById("root"));`,fileName ))
    file.Write(contents)
    file.Sync()
    return newFilePath, nil
}

func buildFile(filePath string) (string, error){
    newFilePath, err := makeRendererFile(filePath)
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
    return string(result.OutputFiles[0].Contents), nil
}

func renderRoute(c *gin.Context, filePath string) {
    compiledJS, err := buildFile(filePath)
    if err != nil {
        c.String(500, err.Error())
    }else{
        c.HTML(http.StatusOK, "index.html", gin.H{
            "src": template.JS(compiledJS),
        })
    }
}

func main() {
    router := gin.Default()
    router.Static("/public", "./public")
    router.LoadHTMLGlob("public/*")
    // createTmpDir()
    router.GET("/", func(c *gin.Context) {
        // addRendererToFile("App.tsx")
        result := esbuildApi.Build(esbuildApi.BuildOptions{
            EntryPoints:       []string{"./frontend/src/App.tsx"},
            // EntryPoints:       []string{TEMP_PATH+"go-ssr-main.jsx"},
            Bundle:            true,
            MinifyWhitespace:  true,
            MinifyIdentifiers: true,
            MinifySyntax:      true,
            Write: true,
            // Outfile:  ".tmp/out.js",
    
        })
        if len(result.Errors) > 0 {
            for _, err := range result.Errors {
                fmt.Println(err.Text)
            }
           c.String(500, result.Errors[0].Text)
        }else{
            // c.String(200, string(result.OutputFiles[0].Contents))
            c.HTML(http.StatusOK, "index.html", gin.H{
                "src": template.JS(string(result.OutputFiles[0].Contents)),
            })
        }
        
        // c.HTML(http.StatusOK, "index.html", gin.H{
        //     "src": "out.js",
        // })
    })
    router.GET("/home", func(c *gin.Context) {
        renderRoute(c, "./frontend/src/Home.tsx")
    })
    router.GET("/test", func(c *gin.Context) {
        c.HTML(http.StatusOK, "test.html", gin.H{})
    })
    router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}