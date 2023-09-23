package react_renderer

import (
	"sync"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Stores the route IDS that render a specific file
var routeIDToReactFileMap = map[string]string{}
var routeIDToReactFileMapLock = sync.RWMutex{}

// Updates the RouteToFileMap with the new file path
func updateRouteIDToReactFileMap(routeID, reactFilePath string) {
	routeIDToReactFileMapLock.Lock()
	defer routeIDToReactFileMapLock.Unlock()
	routeIDToReactFileMap[routeID] = reactFilePath
}

// Returns any routes that render a parent file
func GetRouteIDSForReactFile(reactFilePath string) []string {
	routeIDToReactFileMapLock.RLock()
	defer routeIDToReactFileMapLock.RUnlock()
	var routes []string
	for route, filePath := range routeIDToReactFileMap {
		if filePath == reactFilePath {
			routes = append(routes, route)
		}
	}
	return routes
}

// Returns all the route IDS that have been rendered
func GetAllRouteIDS() []string {
	routeIDToReactFileMapLock.RLock()
	defer routeIDToReactFileMapLock.RUnlock()
	var routes []string
	for route := range routeIDToReactFileMap {
		routes = append(routes, route)
	}
	return routes
}

// Store the react files and the depdenencies they import
var parentFileToDependenciesMap = map[string][]string{}
var parentFileToDependenciesMapLock = sync.RWMutex{}

// Updates the parentFileToDependenciesMap with the new dependencies
func updateParentFileDependencies(reactFilePath string, dependencies []string) {
	parentFileToDependenciesMapLock.Lock()
	defer parentFileToDependenciesMapLock.Unlock()
	parentFileToDependenciesMap[reactFilePath] = dependencies
}

// Returns any parent files that import the dependency
func getParentFilesFromDependency(dependencyPath string) []string {
	parentFileToDependenciesMapLock.RLock()
	defer parentFileToDependenciesMapLock.RUnlock()
	var parentFilePaths []string
	for parentFilePath, dependencies := range parentFileToDependenciesMap {
		for _, dependency := range dependencies {
			if dependency == dependencyPath {
				parentFilePaths = append(parentFilePaths, parentFilePath)
			}
		}
	}
	return parentFilePaths
}

// Takes in a file path and return any routeID's that either render the file
// or the file they render imports that file as a dependency
func GetRouteIDSWithFile(fileName string) []string {
	filePath := utils.GetFullFilePath(fileName)
	// Get the parent files that import the file as a dependency
	reactFilesWithDependency := getParentFilesFromDependency(filePath)
	if len(reactFilesWithDependency) == 0 {
		// If the file is not imported as a dependency, check if it is rendered directly
		reactFilesWithDependency = []string{filePath}
	}
	var routeIDS []string
	for _, reactFile := range reactFilesWithDependency {
		routeIDS = append(routeIDS, GetRouteIDSForReactFile(reactFile)...)
	}
	return routeIDS
}
