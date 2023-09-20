package hot_reload

func Init() {
	go StartServer()
	go WatchForFileChanges()
}
