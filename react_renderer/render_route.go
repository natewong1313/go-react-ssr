package react_renderer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"runtime"
	"strings"

	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/logger"
)

type Config struct {
	File     string
	Title    string
	MetaTags map[string]string
	Props    interface{}
}

// Converts the given react file path to a full html page
func RenderRoute(renderConfig Config) []byte {
	// _, file, lineNum, _ := runtime.Caller(1)
	// routeID := file + ":" + fmt.Sprint(lineNum)
	// Get the program counter for the caller of this function and use that for the id
	pc, _, _, _ := runtime.Caller(1)
	// Create a unique route ID to differentiate between routes that return the same file
	routeID := fmt.Sprint(pc)
	// Props are passed to the renderer as a JSON string, or set to null if no props are passed
	props := "null"
	if renderConfig.Props != nil {
		// Convert props to a JSON string
		propsJSON, err := json.Marshal(renderConfig.Props)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to convert props to JSON")
			return renderErrorHTMLString(err)
		}
		props = string(propsJSON)
	}
	// Get the full path of the react component file
	reactFilePath := getFullFilePath(config.C.FrontendDir + "/" + renderConfig.File)
	// updateRouteToFileMap(routeID, filePath)
	// cachedBuild, ok := checkForCachedBuild(filePath)
	// if !ok {
	// 	// If a build hasn't been cached for this file (or the file has been updated), build it
	// 	var err error
	// 	var metafile string
	// 	// metafile contains the dependencies of the file, this is used to invalidate the cache when a dependency changes
	// 	cachedBuild, metafile, err = BuildFile(filePath, props)
	// 	if err != nil {
	// 		logger.L.Error().Err(err).Msg("Error occured building file")
	// 		return renderErrorHTMLString(err)
	// 	}
	// 	updateFileToDependenciesMap(filePath, getDependenciesFromMetafile(metafile))
	// 	cacheBuild(filePath, cachedBuild)
	// }
	// builtReactFile, cachedBuildFound := checkForCachedBuild(routeID, reactFilePath)

	// if !cachedBuildFound {
	// 	newBuild, dependencyPaths, err := buildReactFile(routeID, reactFilePath, props)
	// 	if err != nil {
	// 		logger.L.Err(err).Msg("Error occured building file")
	// 		return renderErrorHTMLString(err)
	// 	}
	// 	updateFileToDependenciesMap(reactFilePath, dependencyPaths)
	// 	cacheBuild(routeID, reactFilePath, newBuild)
	// 	builtReactFile = newBuild
	// }
	builtReactFile, _, err := buildReactFile(routeID, reactFilePath, props)
	if err != nil {
		logger.L.Err(err).Msg("Error occured building file")
		return renderErrorHTMLString(err)
	}
	return renderHTMLString(HTMLParams{
		Title:      renderConfig.Title,
		MetaTags:   getMetaTags(renderConfig.MetaTags),
		OGMetaTags: getOGMetaTags(renderConfig.MetaTags),
		JS:         template.JS(builtReactFile.CompiledJS),
		CSS:        template.CSS(builtReactFile.CompiledCSS),
		RouteID:    routeID,
	})
}

// Differentiate between meta tags and open graph meta tags

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
