package utils

import "path/filepath"

// Returns the absolute path of the file without returning an error
func GetFullFilePath(filePath string) string {
	fp, _ := filepath.Abs(filePath)
	return filepath.ToSlash(fp)
}
