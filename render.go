package go_ssr

import (
	"encoding/json"
	"fmt"
	"github.com/natewong1313/go-react-ssr/internal/html"
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

func (engine *Engine) RenderRoute(renderConfig RenderConfig) []byte {
	// routeID is the program counter of the caller
	pc, _, _, _ := runtime.Caller(1)
	routeID := fmt.Sprint(pc)

	props, err := propsToString(renderConfig.Props)
	if err != nil {
		return html.RenderError(err, routeID)
	}
	task := renderTask{
		engine:             engine,
		logger:             zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		routeID:            routeID,
		props:              props,
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

// Convert props to JSON string, or set to null if no props are passed
func propsToString(props interface{}) (string, error) {
	if props != nil {
		propsJSON, err := json.Marshal(props)
		return string(propsJSON), err
	}
	return "null", nil
}
