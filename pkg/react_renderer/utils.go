package react_renderer

import "path/filepath"


func getFullFilePath(filePath string) string {
	fp, _ := filepath.Abs(filePath)
	return fp
}