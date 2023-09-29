package react_renderer

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

func serverRenderReactFile(reactFilePath, props string, serverBuildResultChan chan<- ServerBuildResult) {
	// Check if the build is cached
	serverRendererBuild, ok := getCachedServerBuild(reactFilePath)
	if !ok {
		var err error
		serverRendererBuild, err = buildReactServerRendererFile(reactFilePath)
		if err != nil {
			serverBuildResultChan <- ServerBuildResult{Error: err}
			return
		}
		setCachedServerBuild(reactFilePath, serverRendererBuild)
	}
	serverHTML, err := renderReactToHTML(serverRendererBuild.JS, props)
	if err != nil {
		serverBuildResultChan <- ServerBuildResult{Error: err}
		return
	}
	serverBuildResultChan <- ServerBuildResult{HTML: serverHTML, CSS: serverRendererBuild.CSS}
}

func buildReactServerRendererFile(reactFilePath string) (ServerRendererBuild, error) {
	defer utils.Timer("buildReactServerRendererFile")()
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents: fmt.Sprintf(`import { renderToString } from "react-dom/server";
			import React from "react";
			function render() {
				const App = require("%s").default;
				return renderToString(<App {...props} />);
			  }
			  globalThis.render = render;
		  `, reactFilePath),
			Loader:     esbuildApi.LoaderTSX,
			ResolveDir: config.C.FrontendDir,
		},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Outdir:            "/",
		AssetNames:        fmt.Sprintf("%s/[name]", strings.TrimPrefix(config.C.AssetRoute, "/")),
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
		return ServerRendererBuild{}, fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, buildResult.Errors[0].Location.File, buildResult.Errors[0].Location.LineText)
	}

	var js string
	var css string
	for _, file := range buildResult.OutputFiles {
		if strings.HasSuffix(file.Path, "stdin.js") {
			js = string(file.Contents)
		} else if strings.HasSuffix(file.Path, "stdin.css") {
			css = string(file.Contents)
		}
	}
	return ServerRendererBuild{JS: js, CSS: css}, nil
}

func renderReactToHTML(rendererJS, props string) (string, error) {
	defer utils.Timer("renderReactToHTML")()
	vm := goja.New()
	err := injectTextEncoderPolyfill(vm)
	if err != nil {
		return "", err
	}
	_, err = vm.RunString(rendererJS + fmt.Sprintf(`
var props = %s;`, props))
	if err != nil {
		return "", err
	}
	render, ok := goja.AssertFunction(vm.Get("render"))
	if !ok {
		return "", fmt.Errorf("render is not a function")
	}
	res, err := render(goja.Undefined())
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func injectTextEncoderPolyfill(vm *goja.Runtime) error {
	_, err := vm.RunString(`function TextEncoder() {
	}
	
	TextEncoder.prototype.encode = function (string) {
	  var octets = [];
	  var length = string.length;
	  var i = 0;
	  while (i < length) {
		var codePoint = string.codePointAt(i);
		var c = 0;
		var bits = 0;
		if (codePoint <= 0x0000007F) {
		  c = 0;
		  bits = 0x00;
		} else if (codePoint <= 0x000007FF) {
		  c = 6;
		  bits = 0xC0;
		} else if (codePoint <= 0x0000FFFF) {
		  c = 12;
		  bits = 0xE0;
		} else if (codePoint <= 0x001FFFFF) {
		  c = 18;
		  bits = 0xF0;
		}
		octets.push(bits | (codePoint >> c));
		c -= 6;
		while (c >= 0) {
		  octets.push(0x80 | ((codePoint >> c) & 0x3F));
		  c -= 6;
		}
		i += codePoint >= 0x10000 ? 2 : 1;
	  }
	  return octets;
	};
	
	function TextDecoder() {
	}
	
	TextDecoder.prototype.decode = function (octets) {
	  var string = "";
	  var i = 0;
	  while (i < octets.length) {
		var octet = octets[i];
		var bytesNeeded = 0;
		var codePoint = 0;
		if (octet <= 0x7F) {
		  bytesNeeded = 0;
		  codePoint = octet & 0xFF;
		} else if (octet <= 0xDF) {
		  bytesNeeded = 1;
		  codePoint = octet & 0x1F;
		} else if (octet <= 0xEF) {
		  bytesNeeded = 2;
		  codePoint = octet & 0x0F;
		} else if (octet <= 0xF4) {
		  bytesNeeded = 3;
		  codePoint = octet & 0x07;
		}
		if (octets.length - i - bytesNeeded > 0) {
		  var k = 0;
		  while (k < bytesNeeded) {
			octet = octets[i + k + 1];
			codePoint = (codePoint << 6) | (octet & 0x3F);
			k += 1;
		  }
		} else {
		  codePoint = 0xFFFD;
		  bytesNeeded = octets.length - i;
		}
		string += String.fromCodePoint(codePoint);
		i += bytesNeeded + 1;
	  }
	  return string
	};`)
	return err
}
