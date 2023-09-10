package react_renderer

import (
	"path/filepath"
	"sync"
)

var cachedBuilds = make(map[string]string)
var cachedBuildsLock = sync.RWMutex{}

func checkForCachedBuild(filePath string) (string, bool) {
	cachedBuildsLock.RLock()
    defer cachedBuildsLock.RUnlock()
	cachedBuild, ok := cachedBuilds[getFullFilePath(filePath)]
	return cachedBuild, ok
}

func cacheBuild(filePath, build string) {
	cachedBuildsLock.Lock()
    defer cachedBuildsLock.Unlock()
	cachedBuilds[getFullFilePath(filePath)] = build
}

func UpdateCacheOnFileChange(filePath string) {
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	delete(cachedBuilds, getFullFilePath(filePath))
}

func getFullFilePath(filePath string) string {
	fp, _ := filepath.Abs(filePath)
	return fp
}