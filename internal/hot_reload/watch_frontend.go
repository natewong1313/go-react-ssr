package hot_reload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

var watcher *fsnotify.Watcher

// https://gist.github.com/sdomino/74980d69f9fa80cb9d73#file-watch_recursive-go
// Watches for file changes in the specified src directory
func WatchForFileChanges() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(config.C.FrontendDir, watchFilesInDir); err != nil {
		logger.L.Error().Err(err).Msg("Failed to add files in directory to watcher")
	}
	for {
		select {
		// Watch for file changes
		case event := <-watcher.Events:
			// Watch for file created, deleted, updated, or renamed events
			if event.Op.String() != "CHMOD" && !strings.Contains(event.Name, "-gossr-temporary") {
				filePath := event.Name
				logger.L.Info().Msgf("File changed: %s, reloading", filePath)
				// Get all route ids that use that file or have it as a dependency
				routeIDS := react_renderer.GetRouteIDSWithFile(filePath)
				// Tell all browser clients listening for those route ids to reload
				go BroadcastFileUpdateToClients(routeIDS)
			}
		case err := <-watcher.Errors:
			logger.L.Error().Err(err).Msg("Error watching file")
		}
	}
}

func watchFilesInDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}
