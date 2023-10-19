package go_ssr

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/dop251/goja"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/internal/html"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/rs/zerolog"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type RenderConfig struct {
	File     string
	Title    string
	MetaTags map[string]string
	Props    interface{}
}

type renderTask struct {
	engine             *Engine
	logger             zerolog.Logger
	routeID            string
	filePath           string
	props              string
	config             RenderConfig
	serverRenderResult chan ServerRenderResult
	clientRenderResult chan ClientRenderResult
}

func (engine *Engine) RenderRoute(renderConfig RenderConfig) []byte {
	task := renderTask{
		engine:             engine,
		logger:             zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		routeID:            getRouteID(),
		filePath:           filepath.ToSlash(utils.GetFullFilePath(engine.Config.FrontendDir + "/" + renderConfig.File)),
		config:             renderConfig,
		serverRenderResult: make(chan ServerRenderResult),
		clientRenderResult: make(chan ClientRenderResult),
	}
	serverRenderResult, clientRenderResult, err := task.Render()
	if err != nil {
		return html.RenderError(err, task.routeID)
	}
	return html.RenderHTMLString(html.Params{
		Title:      renderConfig.Title,
		MetaTags:   getMetaTags(renderConfig.MetaTags),
		OGMetaTags: getOGMetaTags(renderConfig.MetaTags),
		JS:         template.JS(clientRenderResult.js),
		CSS:        template.CSS(serverRenderResult.css),
		RouteID:    task.routeID,
		ServerHTML: template.HTML(serverRenderResult.html),
	})
}

func getRouteID() string {
	pc, _, _, _ := runtime.Caller(2)
	return fmt.Sprint(pc)
}

func getMetaTags(metaTags map[string]string) map[string]string {
	newMetaTags := make(map[string]string)
	for key, value := range metaTags {
		if !strings.HasPrefix(key, "og:") {
			newMetaTags[key] = value
		}
	}
	return newMetaTags
}

func getOGMetaTags(metaTags map[string]string) map[string]string {
	newMetaTags := make(map[string]string)
	for key, value := range metaTags {
		if strings.HasPrefix(key, "og:") {
			newMetaTags[key] = value
		}
	}
	return newMetaTags
}

func (rt *renderTask) Render() (ServerRenderResult, ClientRenderResult, error) {
	props, err := propsToString(rt.config.Props)
	if err != nil {
		return ServerRenderResult{}, ClientRenderResult{}, err
	}
	rt.props = props

	rt.engine.CacheManager.SetParentFile(rt.routeID, rt.filePath)

	go rt.serverRender()
	go rt.clientRender()

	serverRenderResult := <-rt.serverRenderResult
	if serverRenderResult.err != nil {
		return ServerRenderResult{}, ClientRenderResult{}, serverRenderResult.err
	}
	clientBuildResult := <-rt.clientRenderResult
	if clientBuildResult.err != nil {
		return ServerRenderResult{}, ClientRenderResult{}, clientBuildResult.err
	}

	go rt.engine.CacheManager.SetParentFileDependencies(rt.filePath, clientBuildResult.dependencies)
	return serverRenderResult, clientBuildResult, nil
}

// Convert props to JSON string, or set to null if no props are passed
func propsToString(props interface{}) (string, error) {
	if props != nil {
		propsJSON, err := json.Marshal(props)
		return string(propsJSON), err
	}
	return "null", nil
}

type ServerRenderResult struct {
	html string
	css  string
	err  error
}

func (rt *renderTask) serverRender() {
	serverBuild, ok := rt.engine.CacheManager.GetServerBuild(rt.filePath)
	if !ok {
		build, err := rt.buildReactServerFile()
		if err != nil {
			rt.serverRenderResult <- ServerRenderResult{err: err}
			return
		}
		serverBuild = build
		rt.engine.CacheManager.SetServerBuild(rt.filePath, serverBuild)
	}
	renderedHTML, err := rt.renderReactToHTML(serverBuild.js)
	rt.serverRenderResult <- ServerRenderResult{html: renderedHTML, css: serverBuild.css, err: err}
}

type ServerBuild struct {
	js  string
	css string
}

func (rt *renderTask) buildReactServerFile() (ServerBuild, error) {
	var layoutCSSImport string
	var layoutImport string
	renderStatement := `renderToString(<App {...props} />)`
	if rt.engine.Config.LayoutCSSFilePath != "" {
		layoutCSSImport = fmt.Sprintf(`import "%s";`, rt.engine.CachedLayoutCSSFilePath)
	}
	if rt.engine.Config.LayoutFilePath != "" {
		layoutImport = fmt.Sprintf(`import Layout from "%s";`, rt.engine.Config.LayoutFilePath)
		renderStatement = `renderToString(<Layout><App {...props} /></Layout>)`
	}
	buildResult := esbuildApi.Build(esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents: fmt.Sprintf(`import { renderToString } from "react-dom/server";
			import React from "react";
			%s
			%s
			function render() {
				const App = require("%s").default;
				return %s;
			  }
			  globalThis.render = render;
		  `, layoutCSSImport, layoutImport, rt.filePath, renderStatement),
			Loader:     esbuildApi.LoaderTSX,
			ResolveDir: rt.engine.Config.FrontendDir,
		},
		Bundle:            true,
		Write:             false,
		Outdir:            "/",
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		AssetNames:        fmt.Sprintf("%s/[name]", strings.TrimPrefix(rt.engine.Config.AssetRoute, "/")),
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
		fileLocation := "unknown"
		lineNum := "unknown"
		if buildResult.Errors[0].Location != nil {
			fileLocation = buildResult.Errors[0].Location.File
			lineNum = buildResult.Errors[0].Location.LineText
		}
		return ServerBuild{}, fmt.Errorf("%s <br>in %s <br>at %s", buildResult.Errors[0].Text, fileLocation, lineNum)
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

	// Write js to file, future use
	cacheDir, err := utils.GetServerBuildCacheDir(rt.config.File)
	if err != nil {
		return ServerBuild{}, err
	}
	jsFilePath := fmt.Sprintf("%s/render.js", cacheDir)
	if err := os.WriteFile(jsFilePath, []byte(js), 0644); err != nil {
		return ServerBuild{}, err
	}

	return ServerBuild{js: js, css: css}, nil
}

func (rt *renderTask) renderReactToHTML(rendererJS string) (string, error) {
	vm := goja.New()
	if err := injectTextEncoderPolyfill(vm); err != nil {
		return "", err
	}
	if err := injectConsolePolyfill(vm); err != nil {
		return "", err
	}
	if _, err := vm.RunString(rendererJS + fmt.Sprintf(`var props = %s;`, rt.props)); err != nil {
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

func injectConsolePolyfill(vm *goja.Runtime) error {
	_, err := vm.RunString(`(function(global) {
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
	return err
}

type ClientRenderResult struct {
	js           string
	dependencies []string
	err          error
}

func (rt *renderTask) clientRender() {
	clientBuild, ok := rt.engine.CacheManager.GetClientBuild(rt.filePath)
	if !ok {
		build, err := rt.buildReactClientFile()
		if err != nil {
			rt.clientRenderResult <- ClientRenderResult{err: err}
			return
		}
		clientBuild = build
		rt.engine.CacheManager.SetClientBuild(rt.filePath, clientBuild)
	}
	js := injectProps(clientBuild.js, rt.props)
	rt.clientRenderResult <- ClientRenderResult{js: js, dependencies: clientBuild.dependencies}
}

type ClientBuild struct {
	js           string
	dependencies []string
}

func (rt *renderTask) buildReactClientFile() (ClientBuild, error) {
	var layoutCSSImport string
	if rt.engine.CachedLayoutCSSFilePath != "" {
		layoutCSSImport = fmt.Sprintf(`import "%s";`, rt.engine.CachedLayoutCSSFilePath)
	}
	var layoutImport string
	renderStatement := `hydrateRoot(document.getElementById("root"), <App {...props} />);`
	if rt.engine.Config.LayoutFilePath != "" {
		layoutImport = fmt.Sprintf(`import Layout from "%s";`, rt.engine.Config.LayoutFilePath)
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
				layoutCSSImport, layoutImport, rt.filePath, renderStatement),
			Loader:     getLoaderType(rt.filePath),
			ResolveDir: rt.engine.Config.FrontendDir,
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
	return ClientBuild{js: js, dependencies: getDependencyPathsFromMetafile(buildResult.Metafile)}, nil
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

// injectProps injects the props into the already compiled JS
func injectProps(compiledJS, props string) string {
	return fmt.Sprintf(`window.props = %s; %s`, props, compiledJS)
}
