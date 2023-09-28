package react_renderer

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/config"
)

func buildForServer(reactFilePath, props string, c chan<- ServerBuildResult) {
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			// const App = require("%s").default;
			// const props = %s
			// console.log(renderToString(<App {...props} />));
			Contents: fmt.Sprintf(`import { renderToString } from "react-dom/server";
			import React from "react";
	

			function render() {
				const App = require("%s").default;
				const props = %s
				return renderToString(<App {...props} />);
			  }
			  globalThis.render = render;
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

	vm := goja.New()
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
	};
	(function(global) {
		'use strict';
		if (!global.console) {
		  global.console = {};
		}
		var con = global.console;
		var prop, method;
		var dummy = function() {};
		var properties = ['memory'];
		var methods = ('assert,clear,count,debug,dir,dirxml,error,exception,group,' +
		   'groupCollapsed,groupEnd,info,log,markTimeline,profile,profiles,profileEnd,' +
		   'show,table,time,timeEnd,timeline,timelineEnd,timeStamp,trace,warn,timeLog,trace').split(',');
		while (prop = properties.pop()) if (!con[prop]) con[prop] = {};
		while (method = methods.pop()) if (!con[method]) con[method] = dummy;
	  })(typeof window === 'undefined' ? this : window);`)
	if err != nil {
		c <- ServerBuildResult{Error: err}
		return
	}
	_, err = vm.RunString(string(buildResult.OutputFiles[0].Contents))
	if err != nil {
		c <- ServerBuildResult{Error: err}
		return
	}
	render, ok := goja.AssertFunction(vm.Get("render"))
	if !ok {
		c <- ServerBuildResult{Error: err}
		return
	}
	res, err := render(goja.Undefined())
	if err != nil {
		c <- ServerBuildResult{Error: err}
		return
	}
	c <- ServerBuildResult{HTML: res.String(), CSS: css}
}
