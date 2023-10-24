package go_ssr

import (
	"encoding/json"
	"fmt"
	"github.com/natewong1313/go-react-ssr/internal/html"
	"github.com/natewong1313/go-react-ssr/internal/reactbuilder"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/rs/zerolog"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
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
	pc, _, _, _ := runtime.Caller(1)
	task := renderTask{
		engine:             engine,
		logger:             zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		routeID:            fmt.Sprint(pc),
		filePath:           filepath.ToSlash(utils.GetFullFilePath(engine.Config.FrontendDir + "/" + renderConfig.File)),
		config:             renderConfig,
		serverRenderResult: make(chan ServerRenderResult),
		clientRenderResult: make(chan ClientRenderResult),
	}
	serverRenderResult, clientRenderResult, err := task.Start()
	if err != nil {
		return html.RenderError(err, task.routeID)
	}
	return html.RenderHTMLString(html.Params{
		Title:      renderConfig.Title,
		MetaTags:   renderConfig.MetaTags,
		JS:         template.JS(clientRenderResult.js),
		CSS:        template.CSS(serverRenderResult.css),
		RouteID:    task.routeID,
		ServerHTML: template.HTML(serverRenderResult.html),
	})
}

func (rt *renderTask) Start() (ServerRenderResult, ClientRenderResult, error) {
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
		rt.logger.Error().Err(serverRenderResult.err).Msg("Failed to build for server")
		return ServerRenderResult{}, ClientRenderResult{}, serverRenderResult.err
	}
	clientBuildResult := <-rt.clientRenderResult
	if clientBuildResult.err != nil {
		rt.logger.Error().Err(clientBuildResult.err).Msg("Failed to build for client")
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
		rt.engine.CacheManager.SetServerBuild(rt.filePath, build)
	}
	js := injectProps(serverBuild.JS, rt.props)
	serverRenderJSFilePath, err := rt.saveServerRenderFile(js)
	if err != nil {
		rt.serverRenderResult <- ServerRenderResult{err: err}
		return
	}
	renderedHTML, err := reactbuilder.RenderReactToHTML(serverRenderJSFilePath)
	rt.serverRenderResult <- ServerRenderResult{html: renderedHTML, css: serverBuild.CSS, err: err}
}

func (rt *renderTask) buildReactServerFile() (reactbuilder.BuildResult, error) {
	var imports []string
	if rt.engine.CachedLayoutCSSFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import "%s";`, rt.engine.CachedLayoutCSSFilePath))
	}
	if rt.engine.Config.LayoutFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import Layout from "%s";`, rt.engine.Config.LayoutFilePath))
	}
	contents, err := reactbuilder.GenerateServerBuildContents(imports, rt.filePath, rt.engine.Config.LayoutFilePath != "")
	if err != nil {
		return reactbuilder.BuildResult{}, err
	}
	buildResult, err := reactbuilder.BuildServer(contents, rt.engine.Config.FrontendDir, rt.engine.Config.AssetRoute)
	if err != nil {
		return reactbuilder.BuildResult{}, err
	}

	return buildResult, nil
}

func (rt *renderTask) saveServerRenderFile(js string) (string, error) {
	cacheDir, err := utils.GetServerBuildCacheDir(rt.routeID)
	if err != nil {
		return "", err
	}
	jsFilePath := fmt.Sprintf("%s/render.js", cacheDir)
	// Write file if not exists
	if err = os.WriteFile(jsFilePath, []byte(js), 0644); err != nil {
		return "", err
	}
	return jsFilePath, nil
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
	js := injectProps(clientBuild.JS, rt.props)
	rt.clientRenderResult <- ClientRenderResult{js: js, dependencies: clientBuild.Dependencies}
}

func (rt *renderTask) buildReactClientFile() (reactbuilder.BuildResult, error) {
	var imports []string
	if rt.engine.CachedLayoutCSSFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import "%s";`, rt.engine.CachedLayoutCSSFilePath))
	}
	if rt.engine.Config.LayoutFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import Layout from "%s";`, rt.engine.Config.LayoutFilePath))
	}
	contents, err := reactbuilder.GenerateClientBuildContents(imports, rt.filePath, rt.engine.Config.LayoutFilePath != "")
	if err != nil {
		return reactbuilder.BuildResult{}, err
	}
	buildResult, err := reactbuilder.BuildClient(contents, rt.engine.Config.FrontendDir, rt.engine.Config.AssetRoute)
	if err != nil {
		return reactbuilder.BuildResult{}, err
	}
	return buildResult, nil
}

// injectProps injects the props into the already compiled JS
func injectProps(compiledJS, props string) string {
	return fmt.Sprintf(`var props = %s; %s`, props, compiledJS)
}
