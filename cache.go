package go_ssr

import "sync"

type Cache struct {
	ServerBuilds             *ServerBuilds
	ClientBuilds             *ClientBuilds
	RouteIDToParentFile      *RouteIDToParentFile
	ParentFileToDependencies *ParentFileToDependencies
}

func NewCache() *Cache {
	return &Cache{
		ServerBuilds: &ServerBuilds{
			Builds: make(map[string]ServerBuild),
			Lock:   sync.RWMutex{},
		},
		ClientBuilds: &ClientBuilds{
			Builds: make(map[string]ClientBuild),
			Lock:   sync.RWMutex{},
		},
		RouteIDToParentFile: &RouteIDToParentFile{
			ReactFiles: make(map[string]string),
			Lock:       sync.RWMutex{},
		},
		ParentFileToDependencies: &ParentFileToDependencies{
			Dependencies: make(map[string][]string),
			Lock:         sync.RWMutex{},
		},
	}
}

type ServerBuilds struct {
	Builds map[string]ServerBuild
	Lock   sync.RWMutex
}

func (cache *Cache) GetServerBuild(filePath string) (ServerBuild, bool) {
	cache.ServerBuilds.Lock.RLock()
	defer cache.ServerBuilds.Lock.RUnlock()
	build, ok := cache.ServerBuilds.Builds[filePath]
	return build, ok
}

func (cache *Cache) SetServerBuild(filePath string, build ServerBuild) {
	cache.ServerBuilds.Lock.Lock()
	defer cache.ServerBuilds.Lock.Unlock()
	cache.ServerBuilds.Builds[filePath] = build
}

func (cache *Cache) RemoveServerBuild(filePath string) {
	cache.ServerBuilds.Lock.Lock()
	defer cache.ServerBuilds.Lock.Unlock()
	if _, ok := cache.ServerBuilds.Builds[filePath]; !ok {
		return
	}
	delete(cache.ServerBuilds.Builds, filePath)
}

type ClientBuilds struct {
	Builds map[string]ClientBuild
	Lock   sync.RWMutex
}

func (cache *Cache) GetClientBuild(filePath string) (ClientBuild, bool) {
	cache.ClientBuilds.Lock.RLock()
	defer cache.ClientBuilds.Lock.RUnlock()
	build, ok := cache.ClientBuilds.Builds[filePath]
	return build, ok
}

func (cache *Cache) SetClientBuild(filePath string, build ClientBuild) {
	cache.ClientBuilds.Lock.Lock()
	defer cache.ClientBuilds.Lock.Unlock()
	cache.ClientBuilds.Builds[filePath] = build
}

func (cache *Cache) RemoveClientBuild(filePath string) {
	cache.ClientBuilds.Lock.Lock()
	defer cache.ClientBuilds.Lock.Unlock()
	if _, ok := cache.ClientBuilds.Builds[filePath]; !ok {
		return
	}
	delete(cache.ClientBuilds.Builds, filePath)
}

type RouteIDToParentFile struct {
	ReactFiles map[string]string
	Lock       sync.RWMutex
}

func (cache *Cache) SetParentFile(routeID, filePath string) {
	cache.RouteIDToParentFile.Lock.Lock()
	defer cache.RouteIDToParentFile.Lock.Unlock()
	cache.RouteIDToParentFile.ReactFiles[routeID] = filePath
}

func (cache *Cache) GetRouteIDSForParentFile(filePath string) []string {
	cache.RouteIDToParentFile.Lock.RLock()
	defer cache.RouteIDToParentFile.Lock.RUnlock()
	var routes []string
	for route, file := range cache.RouteIDToParentFile.ReactFiles {
		if file == filePath {
			routes = append(routes, route)
		}
	}
	return routes
}

func (cache *Cache) GetAllRouteIDS() []string {
	cache.RouteIDToParentFile.Lock.RLock()
	defer cache.RouteIDToParentFile.Lock.RUnlock()
	routes := make([]string, 0, len(cache.RouteIDToParentFile.ReactFiles))
	for route := range cache.RouteIDToParentFile.ReactFiles {
		routes = append(routes, route)
	}
	return routes
}

type ParentFileToDependencies struct {
	Dependencies map[string][]string
	Lock         sync.RWMutex
}

func (cache *Cache) SetParentFileDependencies(filePath string, dependencies []string) {
	cache.ParentFileToDependencies.Lock.Lock()
	defer cache.ParentFileToDependencies.Lock.Unlock()
	cache.ParentFileToDependencies.Dependencies[filePath] = dependencies
}

func (cache *Cache) GetParentFilesFromDependency(dependencyPath string) []string {
	cache.ParentFileToDependencies.Lock.RLock()
	defer cache.ParentFileToDependencies.Lock.RUnlock()
	var parentFilePaths []string
	for parentFilePath, dependencies := range cache.ParentFileToDependencies.Dependencies {
		for _, dependency := range dependencies {
			if dependency == dependencyPath {
				parentFilePaths = append(parentFilePaths, parentFilePath)
			}
		}
	}
	return parentFilePaths
}
