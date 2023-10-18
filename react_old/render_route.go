package react_old

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/natewong1313/go-react-ssr/internal/html"
//	"github.com/natewong1313/go-react-ssr/internal/logger"
//	"html/template"
//	"path/filepath"
//	"runtime"
//	"strings"
//
//	"github.com/natewong1313/go-react-ssr/config"
//	"github.com/natewong1313/go-react-ssr/internal/utils"
//)
//
//// RenderRoute Converts the given react_old file path to a full html page
//func RenderRoute(renderConfig Config) []byte {
//	// Get the program counter for the caller of this function and use that for the id
//	pc, _, _, _ := runtime.Caller(1)
//	routeID := fmt.Sprint(pc)
//
//	task := RenderTask{
//		RouteID:           routeID,
//		FilePath:          filepath.ToSlash(utils.GetFullFilePath(config.C.FrontendDir + "/" + renderConfig.File)),
//		RenderConfig:      renderConfig,
//		ServerBuildResult: make(chan ServerBuildResult),
//		ClientBuildResult: make(chan ClientBuildResult),
//	}
//	// Props are passed to the renderer as a JSON string, or set to null if no props are passed
//	props, err := propsToString(renderConfig.Props)
//	if err != nil {
//		logger.L.Err(err).Msg("Failed to convert props")
//		return html.RenderError(err)
//	}
//	task.Props = props
//	// Update the routeID to file map
//	go updateRouteIDToReactFileMap(routeID, task.FilePath)
//
//	go task.ServerRender()
//	go task.ClientRender()
//
//	serverBuildResult := <-task.ServerBuildResult
//	if serverBuildResult.Error != nil {
//		logger.L.Err(serverBuildResult.Error).Msg("Error occurred building server rendered file")
//		return html.RenderError(serverBuildResult.Error)
//	}
//
//	clientBuildResult := <-task.ClientBuildResult
//	if clientBuildResult.Error != nil {
//		logger.L.Err(clientBuildResult.Error).Msg("Error occurred building client js file")
//		return html.RenderError(clientBuildResult.Error)
//	}
//
//	go updateParentFileDependencies(task.FilePath, clientBuildResult.Dependencies)
//	// Return the rendered html
//	return html.RenderHTMLString(html.Params{
//		Title:      renderConfig.Title,
//		MetaTags:   getMetaTags(renderConfig.MetaTags),
//		OGMetaTags: getOGMetaTags(renderConfig.MetaTags),
//		Links:      renderConfig.Links,
//		JS:         template.JS(clientBuildResult.JS),
//		CSS:        template.CSS(serverBuildResult.CSS),
//		RouteID:    routeID,
//		ServerHTML: template.HTML(serverBuildResult.HTML),
//	})
//}
//
//// Convert props to JSON string, or set to null if no props are passed
//func propsToString(props interface{}) (string, error) {
//	if props != nil {
//		propsJSON, err := json.Marshal(props)
//		if err != nil {
//			return "", err
//		}
//		return string(propsJSON), nil
//	}
//	return "null", nil
//}
//
//// Differentiate between meta tags and open graph meta tags
//
//func getMetaTags(metaTags map[string]string) map[string]string {
//	newMetaTags := make(map[string]string)
//	for key, value := range metaTags {
//		if !strings.HasPrefix(key, "og:") {
//			newMetaTags[key] = value
//		}
//	}
//	return newMetaTags
//}
//
//func getOGMetaTags(metaTags map[string]string) map[string]string {
//	newMetaTags := make(map[string]string)
//	for key, value := range metaTags {
//		if strings.HasPrefix(key, "og:") {
//			newMetaTags[key] = value
//		}
//	}
//	return newMetaTags
//}
