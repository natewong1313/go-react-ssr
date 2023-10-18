package react

import "sync"

type ParentFileToDependenciesMap struct {
	Map  map[string][]string
	Lock sync.RWMutex
}
