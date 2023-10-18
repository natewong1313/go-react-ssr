package go_ssr

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HotReload struct {
	Engine           *Engine
	Logger           zerolog.Logger
	ConnectedClients map[string][]*websocket.Conn
	Watcher          *fsnotify.Watcher
}

func newHotReload(engine *Engine) *HotReload {
	return &HotReload{
		Engine:           engine,
		Logger:           engine.Logger,
		ConnectedClients: make(map[string][]*websocket.Conn),
	}
}

func (hr *HotReload) Start() {
	go hr.StartServer()
	go hr.StartWatcher()
}

// StartServer starts the hot reload websocket server at the port specified in the config
func (hr *HotReload) StartServer() {
	hr.Logger.Info().Msgf("Hot reload websocket running on port %d", hr.Engine.Config.HotReloadServerPort)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			hr.Logger.Err(err).Msg("Failed to upgrade websocket")
			return
		}
		// Client should send routeID as first message
		_, routeID, err := ws.ReadMessage()
		if err != nil {
			hr.Logger.Err(err).Msg("Failed to read message from websocket")
			return
		}
		err = ws.WriteMessage(1, []byte("Connected"))
		if err != nil {
			hr.Logger.Err(err).Msg("Failed to write message to websocket")
			return
		}
		// Add client to connectedClients
		hr.ConnectedClients[string(routeID)] = append(hr.ConnectedClients[string(routeID)], ws)
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", hr.Engine.Config.HotReloadServerPort), nil)
	if err != nil {
		hr.Logger.Err(err).Msg("Hot reload server quit unexpectedly")
	}
}

func (hr *HotReload) StartWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		hr.Logger.Err(err).Msg("Failed to start watcher")
		return
	}
	defer watcher.Close()

	if err = filepath.Walk(hr.Engine.Config.FrontendDir, func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}
		return nil
	}); err != nil {
		hr.Logger.Err(err).Msg("Failed to add files in directory to watcher")
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			// Watch for file created, deleted, updated, or renamed events
			if event.Op.String() != "CHMOD" && !strings.Contains(event.Name, "gossr-temporary") {
				filePath := utils.GetFullFilePath(event.Name)
				hr.Logger.Info().Msgf("File changed: %s, reloading", filePath)
				// Store the routes that need to be reloaded
				var routeIDS []string
				switch {
				case filePath == utils.GetFullFilePath(hr.Engine.Config.LayoutFile): // If the layout file has been updated, reload all routes
					routeIDS = hr.Engine.Cache.GetAllRouteIDS()
				case hr.layoutCSSFileUpdated(filePath): // If the global css file has been updated, rebuild it and reload all routes
					err := hr.Engine.BuildLayoutCSSFile()
					if err != nil {
						hr.Logger.Err(err).Msg("Failed to build global css file")
						continue
					}
					routeIDS = hr.Engine.Cache.GetAllRouteIDS()
				case hr.needsTailwindRecompile(filePath): // If tailwind is enabled and a React file has been updated, rebuild the global css file and reload all routes
					err := hr.Engine.BuildLayoutCSSFile()
					if err != nil {
						hr.Logger.Err(err).Msg("Failed to build global css file")
						continue
					}
					fallthrough
				default:
					// Get all route ids that use that file or have it as a dependency
					routeIDS = hr.Engine.Cache.GetRouteIDSForParentFile(filePath)
				}
				parentFiles := hr.Engine.Cache.GetParentFilesFromDependency(filePath)
				for _, parentFile := range parentFiles {
					hr.Engine.Cache.RemoveServerBuild(parentFile)
					hr.Engine.Cache.RemoveClientBuild(parentFile)
				}
				// Tell all browser clients listening for those route ids to reload
				go hr.broadcastFileUpdateToClients(routeIDS)

			}
		case err := <-watcher.Errors:
			hr.Logger.Err(err).Msg("Error watching files")
		}
	}
}

func (hr *HotReload) layoutCSSFileUpdated(filePath string) bool {
	return utils.GetFullFilePath(filePath) == utils.GetFullFilePath(hr.Engine.Config.LayoutCSSFilePath)
}

// Check if tailwind is enabled and if the file is a React file.
func (hr *HotReload) needsTailwindRecompile(filePath string) bool {
	if hr.Engine.Config.TailwindConfigPath == "" {
		return false
	}
	fileTypes := []string{".tsx", ".ts", ".jsx", ".js"}
	for _, fileType := range fileTypes {
		if strings.HasSuffix(filePath, fileType) {
			return true
		}
	}
	return false
}

// Tell all clients listening for a specific routeID to reload
func (hr *HotReload) broadcastFileUpdateToClients(routeIDS []string) {
	// Iterate over each route ID
	for _, routeID := range routeIDS {
		// Find all clients listening for that route ID
		for i, ws := range hr.ConnectedClients[routeID] {
			// Send reload message to client
			err := ws.WriteMessage(1, []byte("reload"))
			if err != nil {
				// remove client if browser is closed or page changed
				hr.ConnectedClients[routeID] = append(hr.ConnectedClients[routeID][:i], hr.ConnectedClients[routeID][i+1:]...)
			}
		}
	}
}
