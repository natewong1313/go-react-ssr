package react_renderer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/buger/jsonparser"
)

var routeToFileMap = map[string]string{}
var routeToFileMapLock = sync.RWMutex{}

// Updates the RouteToFileMap with the new file path
func updateRouteToFileMap(route, filePath string) {
	routeToFileMapLock.Lock()
	defer routeToFileMapLock.Unlock()
	routeToFileMap[route] = filePath
}

// Returns any routes that render a parent file
func GetRoutesForFile(filePath string) []string {
	routeToFileMapLock.RLock()
	defer routeToFileMapLock.RUnlock()
	var routes []string
	for route, file := range routeToFileMap {
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
			fmt.Println("adding dep", getFullFilePath(string(key)))
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
			fmt.Println("Found dep", dependency)
			if dependency == filePath {
				return parentFilePath
			}
		}
	}
	return ""
}
