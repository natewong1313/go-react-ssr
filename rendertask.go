package go_ssr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/natewong1313/go-react-ssr/internal/reactbuilder"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/rs/zerolog"
)

type renderTask struct {
	engine             *Engine
	logger             zerolog.Logger
	routeID            string
	filePath           string
	props              string
	config             RenderConfig
	serverRenderResult chan serverRenderResult
	clientRenderResult chan clientRenderResult
}

type serverRenderResult struct {
	html string
	css  string
	err  error
}

type clientRenderResult struct {
	js           string
	dependencies []string
	err          error
}

// Start starts the render task, returns the rendered html, css, and js for hydration
func (rt *renderTask) Start() (string, string, string, error) {
	rt.serverRenderResult = make(chan serverRenderResult)
	rt.clientRenderResult = make(chan clientRenderResult)
	// Assigns the parent file to the routeID so that the cache can be invalidated when the parent file changes
	rt.engine.CacheManager.SetParentFile(rt.routeID, rt.filePath)

	// Render for server and client concurrently
	go rt.doRender("server")
	go rt.doRender("client")

	// Wait for both to finish
	srResult := <-rt.serverRenderResult
	if srResult.err != nil {
		rt.logger.Error().Err(srResult.err).Msg("Failed to build for server")
		return "", "", "", srResult.err
	}
	crResult := <-rt.clientRenderResult
	if crResult.err != nil {
		rt.logger.Error().Err(crResult.err).Msg("Failed to build for client")
		return "", "", "", crResult.err
	}

	// Set the parent file dependencies so that the cache can be invalidated a dependency changes
	go rt.engine.CacheManager.SetParentFileDependencies(rt.filePath, crResult.dependencies)
	return srResult.html, srResult.css, crResult.js, nil
}

func (rt *renderTask) doRender(buildType string) {
	// Check if the build is in the cache
	build, buildFound := rt.getBuildFromCache(buildType)
	if !buildFound {
		// Build the file if it's not in the cache
		newBuild, err := rt.buildFile(buildType)
		if err != nil {
			rt.handleBuildError(err, buildType)
			return
		}
		rt.updateBuildCache(newBuild, buildType)
		build = newBuild
	}
	// JS is built without props so that the props can be injected into cached JS builds
	js := injectProps(build.JS, rt.props)
	if buildType == "server" {
		// Save the server js to a file to be executed by node
		jsFilePath, err := rt.saveServerRenderFile(js)
		if err != nil {
			rt.handleBuildError(err, buildType)
			return
		}
		// Then call that file with node to get the rendered HTML
		renderedHTML, err := renderReactToHTML(jsFilePath, rt.engine.Config.JSRuntime)
		rt.serverRenderResult <- serverRenderResult{html: renderedHTML, css: build.CSS, err: err}
	} else {
		rt.clientRenderResult <- clientRenderResult{js: js, dependencies: build.Dependencies}
	}
}

// getBuild returns the build from the cache if it exists
func (rt *renderTask) getBuildFromCache(buildType string) (reactbuilder.BuildResult, bool) {
	if buildType == "server" {
		return rt.engine.CacheManager.GetServerBuild(rt.filePath)
	} else {
		return rt.engine.CacheManager.GetClientBuild(rt.filePath)
	}
}

// buildFile gets the contents of the file to be built and builds it with reactbuilder
func (rt *renderTask) buildFile(buildType string) (reactbuilder.BuildResult, error) {
	buildContents, err := rt.getBuildContents(buildType)
	if err != nil {
		return reactbuilder.BuildResult{}, err
	}
	if buildType == "server" {
		return reactbuilder.BuildServer(buildContents, rt.engine.Config.FrontendDir, rt.engine.Config.AssetRoute)
	} else {
		return reactbuilder.BuildClient(buildContents, rt.engine.Config.FrontendDir, rt.engine.Config.AssetRoute)
	}
}

// getBuildContents gets the required imports based on the config and returns the contents to be built with reactbuilder
func (rt *renderTask) getBuildContents(buildType string) (string, error) {
	var imports []string
	if rt.engine.CachedLayoutCSSFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import "%s";`, rt.engine.CachedLayoutCSSFilePath))
	}
	if rt.engine.Config.LayoutFilePath != "" {
		imports = append(imports, fmt.Sprintf(`import Layout from "%s";`, rt.engine.Config.LayoutFilePath))
	}
	if buildType == "server" {
		return reactbuilder.GenerateServerBuildContents(imports, rt.filePath, rt.engine.Config.LayoutFilePath != "")
	} else {
		return reactbuilder.GenerateClientBuildContents(imports, rt.filePath, rt.engine.Config.LayoutFilePath != "")
	}
}

// handleBuildError handles the error from building the file and sends it to the appropriate channel
func (rt *renderTask) handleBuildError(err error, buildType string) {
	if buildType == "server" {
		rt.serverRenderResult <- serverRenderResult{err: err}
	} else {
		rt.clientRenderResult <- clientRenderResult{err: err}
	}
}

// updateBuildCache updates the cache with the new build
func (rt *renderTask) updateBuildCache(build reactbuilder.BuildResult, buildType string) {
	if buildType == "server" {
		rt.engine.CacheManager.SetServerBuild(rt.filePath, build)
	} else {
		rt.engine.CacheManager.SetClientBuild(rt.filePath, build)
	}
}

// injectProps injects the props into the already compiled JS
func injectProps(compiledJS, props string) string {
	return fmt.Sprintf(`var props = %s; %s`, props, compiledJS)
}

// saveServerRenderFile saves the generated server js to a file to be executed by node
func (rt *renderTask) saveServerRenderFile(js string) (string, error) {
	cacheDir, err := utils.GetServerBuildCacheDir(rt.routeID)
	if err != nil {
		return "", err
	}
	jsFilePath := fmt.Sprintf("%s/render.js", cacheDir)
	return jsFilePath, os.WriteFile(jsFilePath, []byte(js), 0644)
}

// renderReactToHTML uses node to execute the server js file which outputs the rendered HTML
func renderReactToHTML(jsFilePath string, jsRuntime string) (string, error) {
	command := []string{"node"}
	if jsRuntime == "bun" {
		command = []string{"bun", "run"}
	}
	if jsRuntime == "yarn" {
		command = []string{"yarn", "run"}
	}
	if jsRuntime == "pnpm" {
		command = []string{"pnpm", "run"}
	}
	command = append(command, jsFilePath)
	cmd := exec.Command(command[0], command[1:]...)
	stdOut := new(strings.Builder)
	stdErr := new(strings.Builder)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err := cmd.Run()
	if err != nil {
		return "", errors.New(stdErr.String())
	}
	return stdOut.String(), nil
}
