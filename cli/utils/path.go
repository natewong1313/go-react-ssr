package utils

import (
	"io"
	"os"
)

func CheckPathExists(projectDir string) bool {
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		HandleError(err)
	}
	return true
}

func CheckPathEmpty(projectDir string) bool {
	f, err := os.Open(projectDir)
	if err != nil {
		HandleError(err)
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		return true
	}
	return false
}
