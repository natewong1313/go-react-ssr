package react_renderer

import (
	"sync"

	"github.com/natewong1313/go-react-ssr/internal/logger"
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
func cacheBuild(routeID, filePath string, cachedBuild CachedBuild) {
	logger.L.Debug().Msgf("Caching build for %s", filePath)
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	cachedBuilds[routeID] = cachedBuild
}

// Remove a build from the cache
func UpdateCacheOnFileChange(filePath string) []string {
	filePath = getFullFilePath(filePath)
	filePathFoundInCache := deleteFromCache(filePath)
	if !filePathFoundInCache {
		filePaths := getParentFilePathsFromDependency(filePath)
		deleteFilesFromCache(filePaths)
		return filePaths
		// filePath = getParentFilePathFromDependency(filePath)
		// deleteFromCache(filePath)
	}
	return []string{filePath}
}

// Delete build from cache if exists
func deleteFromCache(filePath string) bool {
	cachedBuildsLock.Lock()
	defer cachedBuildsLock.Unlock()
	_, ok := cachedBuilds[filePath]
	if ok {
		delete(cachedBuilds, filePath)
	}
	return ok
}

func deleteFilesFromCache(filePaths []string) {
	for _, filePath := range filePaths {
		deleteFromCache(filePath)
	}
}
