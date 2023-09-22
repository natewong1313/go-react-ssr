package utils

import "path/filepath"

func GetFullFilePath(filePath string) string {
	fp, _ := filepath.Abs(filePath)
	return fp
}
