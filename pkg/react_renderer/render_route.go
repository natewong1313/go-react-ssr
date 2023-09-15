package react_renderer

import (
	"encoding/json"
	"html/template"
	"strings"

	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/pkg/config"

	"github.com/gin-gonic/gin"
)

type Config struct {
	File     string
	Title    string
	MetaTags map[string]string
	Props    interface{}
}

// Uses esbuild to render the React page at the given file path and inserts it into an html page
func RenderRoute(c *gin.Context, renderConfig Config) {
	props := "null"
	if renderConfig.Props != nil {
		propsJSON, err := json.Marshal(renderConfig.Props)
		if err != nil {
			logger.L.Error().Err(err).Msg("Failed to convert props to JSON")
			c.JSON(500, gin.H{"error": "Invalid prop types"})
			return
		}
		props = string(propsJSON)
	}
	// Get the full path of the file
	filePath := getFullFilePath(config.C.FrontendDir + "/" + renderConfig.File)
	updateRouteToFileMap(c.Request.URL.Path, filePath)
	cachedBuild, ok := checkForCachedBuild(filePath)
	if !ok {
		var err error
		var metafile string
		cachedBuild, metafile, err = BuildFile(filePath, props)
		if err != nil {
			logger.L.Error().Err(err).Msg("Error occured building file")
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		updateFileToDependenciesMap(filePath, getDependenciesFromMetafile(metafile))
		cacheBuild(filePath, cachedBuild)
	}
	c.Writer.Write(renderHTMLString(HTMLParams{
		Title:      renderConfig.Title,
		MetaTags:   getMetaTags(renderConfig.MetaTags),
		OGMetaTags: getOGMetaTags(renderConfig.MetaTags),
		JS:         template.JS(cachedBuild.CompiledJS),
		CSS:        template.CSS(cachedBuild.CompiledCSS),
		Route:      c.Request.URL.Path,
	}))
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
