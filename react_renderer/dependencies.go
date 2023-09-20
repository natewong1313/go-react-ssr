package react_renderer

import (
	"strings"
	"sync"

	"github.com/buger/jsonparser"
)

var routeIDToFileMap = map[string]string{}
var routeIDToFileMapLock = sync.RWMutex{}

// Updates the RouteToFileMap with the new file path
func updateRouteToFileMap(routeID, filePath string) {
	routeIDToFileMapLock.Lock()
	defer routeIDToFileMapLock.Unlock()
	routeIDToFileMap[routeID] = filePath
}

// Returns any routes that render a parent file
func GetRouteIDSForFile(filePath string) []string {
	routeIDToFileMapLock.RLock()
	defer routeIDToFileMapLock.RUnlock()
	var routes []string
	for route, file := range routeIDToFileMap {
		if file == filePath {
			routes = append(routes, route)
		}
	}
	return routes
}

var fileToDependenciesMap = map[string][]string{}
var fileToDependenciesMapLock = sync.RWMutex{}

// Parse dependencies from esbuild metafile
func getDependenciesFromMetafile(metafile string) []string {
	var dependencies []string
	jsonparser.ObjectEach([]byte(metafile), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if !strings.Contains(string(key), "/node_modules/") {
			dependencies = append(dependencies, getFullFilePath(string(key)))
		}
		return nil
	}, "inputs")
	return dependencies
}

// Updates the FileToDependenciesMap with the new file path
func updateFileToDependenciesMap(filePath string, dependencies []string) {
	fileToDependenciesMapLock.Lock()
	defer fileToDependenciesMapLock.Unlock()
	fileToDependenciesMap[filePath] = dependencies
}

// Returns the dependencies for the given file path
func GetDependencies(filePath string) []string {
	fileToDependenciesMapLock.RLock()
	defer fileToDependenciesMapLock.RUnlock()
	return fileToDependenciesMap[filePath]
}

func getParentFilePathFromDependency(filePath string) string {
	fileToDependenciesMapLock.RLock()
	defer fileToDependenciesMapLock.RUnlock()
	for parentFilePath, dependencies := range fileToDependenciesMap {
		for _, dependency := range dependencies {
			if dependency == filePath {
				return parentFilePath
			}
		}
	}
	return ""
}
