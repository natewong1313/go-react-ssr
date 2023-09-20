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

// Uses esbuild to render the React page at the given file path and inserts it into an html page
func RenderRoute(renderConfig Config) []byte {
	_, file, lineNum, _ := runtime.Caller(1)
	routeID := file + ":" + fmt.Sprint(lineNum)

	props := "null"
	if renderConfig.Props != nil {
		propsJSON, err := json.Marshal(renderConfig.Props)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to convert props to JSON")
			return renderErrorHTMLString(err)
		}
		props = string(propsJSON)
	}
	// Get the full path of the file
	filePath := getFullFilePath(config.C.FrontendDir + "/" + renderConfig.File)
	updateRouteToFileMap(routeID, filePath)
	cachedBuild, ok := checkForCachedBuild(filePath)
	if !ok {
		var err error
		var metafile string
		cachedBuild, metafile, err = BuildFile(filePath, props)
		if err != nil {
			logger.L.Error().Err(err).Msg("Error occured building file")
			return renderErrorHTMLString(err)
		}
		updateFileToDependenciesMap(filePath, getDependenciesFromMetafile(metafile))
		cacheBuild(filePath, cachedBuild)
	}
	return renderHTMLString(HTMLParams{
		Title:      renderConfig.Title,
		MetaTags:   getMetaTags(renderConfig.MetaTags),
		OGMetaTags: getOGMetaTags(renderConfig.MetaTags),
		JS:         template.JS(cachedBuild.CompiledJS),
		CSS:        template.CSS(cachedBuild.CompiledCSS),
		RouteID:    routeID,
	})
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
