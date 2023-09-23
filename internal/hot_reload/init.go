package hot_reload

// Init initializes the hot reload server and file watcher
func Init() {
	go StartServer()
	go WatchForFileChanges()
}
