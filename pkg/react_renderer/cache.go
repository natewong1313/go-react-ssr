package react_renderer

import (
	"path/filepath"
	"sync"
)

type CachedBuild struct {
	CompiledJS string
	CompiledCSS string
}

var cachedBuilds = make(map[string]CachedBuild)
var cachedBuildsLock = sync.RWMutex{}

// Find a cached build for the given file path
func checkForCachedBuild(filePath string) (CachedBuild, bool) {
	cachedBuildsLock.RLock()
    defer cachedBuildsLock.RUnlock()
	cachedBuild, ok := cachedBuilds[getFullFilePath(filePath)]
	return cachedBuild, ok
}

// Add a build to the cache
func cacheBuild(filePath string, cachedBuild CachedBuild) {
	cachedBuildsLock.Lock()
    defer cachedBuildsLock.Unlock()
	cachedBuilds[getFullFilePath(filePath)] = cachedBuild
}

// Remove a build from the cache
func UpdateCacheOnFileChange(filePath string) {
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	delete(cachedBuilds, getFullFilePath(filePath))
}

func getFullFilePath(filePath string) string {
	fp, _ := filepath.Abs(filePath)
	return fp
}