package react

import (
	"sync"
)

var cachedServerBuilds = make(map[string]ServerRendererBuild)
var cachedServerBuildsLock = sync.RWMutex{}
var cachedClientBuilds = make(map[string]ClientBuild)
var cachedClientBuildsLock = sync.RWMutex{}

// Get the cached server build for the given routeID
func getCachedServerBuild(reactFilePath string) (ServerRendererBuild, bool) {
	cachedServerBuildsLock.RLock()
	defer cachedServerBuildsLock.RUnlock()
	build, ok := cachedServerBuilds[reactFilePath]
	return build, ok
}

// Get the cached client build for the given routeID
func getCachedClientBuild(reactFilePath string) (ClientBuild, bool) {
	cachedClientBuildsLock.RLock()
	defer cachedClientBuildsLock.RUnlock()
	build, ok := cachedClientBuilds[reactFilePath]
	return build, ok
}

// Set the cached server build for the given routeID
func setCachedServerBuild(routeID string, build ServerRendererBuild) {
	cachedServerBuildsLock.Lock()
	defer cachedServerBuildsLock.Unlock()
	cachedServerBuilds[routeID] = build
}

// Set the cached client build for the given routeID
func setCachedClientBuild(reactFilePath string, build ClientBuild) {
	cachedClientBuildsLock.Lock()
	defer cachedClientBuildsLock.Unlock()
	cachedClientBuilds[reactFilePath] = build
}

func RemoveCachedServerBuild(reactFilePath string) {
	cachedServerBuildsLock.Lock()
	defer cachedServerBuildsLock.Unlock()
	if _, ok := cachedServerBuilds[reactFilePath]; !ok {
		return
	}
	delete(cachedServerBuilds, reactFilePath)
}

func RemoveCachedClientBuild(reactFilePath string) {
	cachedClientBuildsLock.Lock()
	defer cachedClientBuildsLock.Unlock()
	if _, ok := cachedClientBuilds[reactFilePath]; !ok {
		return
	}
	delete(cachedClientBuilds, reactFilePath)
}
