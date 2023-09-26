package create

import (
	"os"
	"path/filepath"
)

func createTempDir() string {
	osCacheDir, _ := os.UserCacheDir()
	tempDirPath := filepath.Join(osCacheDir, "gossr-cli")
	os.RemoveAll(tempDirPath)
	os.MkdirAll(tempDirPath, os.ModePerm)

	return tempDirPath
}
