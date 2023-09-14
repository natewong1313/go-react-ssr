package hot_reload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)

var watcher *fsnotify.Watcher

// https://gist.github.com/sdomino/74980d69f9fa80cb9d73#file-watch_recursive-go
// Watches for file changes in the src directory
func StartWatching() {
	go StartServer()
	go func() {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk(config.C.FrontendDir, watchFilesInDir); err != nil {
			fmt.Println("ERROR", err)
		}
		for {
			select {
			// Watch for file changes
			case event := <-watcher.Events:
				if event.Op.String() != "CHMOD" && !strings.Contains(event.Name, "-gossr-temporary") {
					fmt.Println(event.Name)
					fmt.Printf("EVENT! %#v\n", event)
					parentFilePath := react_renderer.UpdateCacheOnFileChange(event.Name)
					go BroadcastFileUpdateToClients(parentFilePath)
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
}

func watchFilesInDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}
