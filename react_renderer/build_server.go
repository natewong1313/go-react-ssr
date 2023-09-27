package react_renderer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
)

func buildForServer(reactFilePath, props string, c chan<- ServerBuildResult) {
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents: fmt.Sprintf(`import { renderToString } from "react-dom/server";
			import React from "react";
	
			const App = require("%s").default;
			const props = %s
			console.log(renderToString(<App {...props} />));
		  `, reactFilePath, props),
			Loader:     esbuildApi.LoaderTSX,
			ResolveDir: config.C.FrontendDir,
		},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Outdir:            "/",
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
		c <- ServerBuildResult{Error: fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)}
		return
	}

	var css string
	for _, file := range buildResult.OutputFiles {
		if strings.HasSuffix(string(file.Path), ".css") {
			css = string(file.Contents)
			break
		}
	}

	// exec code with node
	node := exec.Command("node", "-e", string(buildResult.OutputFiles[0].Contents))
	var outb bytes.Buffer

	node.Stdout = &outb
	if err := node.Run(); err != nil {
		c <- ServerBuildResult{Error: err}
		return
	}
	c <- ServerBuildResult{HTML: outb.String(), CSS: css}
}
