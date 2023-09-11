package react_renderer

import (
	"fmt"
	"sync"
)

type CachedBuild struct {
	CompiledJS  string
	CompiledCSS string
}

var cachedBuilds = make(map[string]CachedBuild)
var cachedBuildsLock = sync.RWMutex{}

// Find a cached build for the given file path
func checkForCachedBuild(filePath string) (CachedBuild, bool) {
	cachedBuildsLock.RLock()
	defer cachedBuildsLock.RUnlock()
	cachedBuild, ok := cachedBuilds[filePath]
	return cachedBuild, ok
}

// Add a build to the cache
func cacheBuild(filePath string, cachedBuild CachedBuild) {
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	cachedBuilds[filePath] = cachedBuild
}

// Remove a build from the cache
func UpdateCacheOnFileChange(filePath string) string {
	filePath = getFullFilePath(filePath)
	fmt.Println("Looking for", filePath)
	filePathFoundInCache := deleteFromCache(filePath)
	if !filePathFoundInCache {
		fmt.Println("Not found in cache")
		filePath = getParentFilePathFromDependency(filePath)
		fmt.Println("Looking for", filePath)
		deleteFromCache(filePath)
	}
	return filePath
}

func deleteFromCache(filePath string) bool {
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	_, ok := cachedBuilds[filePath]
	if ok {
		delete(cachedBuilds, filePath)
	}
	return ok
}
