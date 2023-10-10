package update

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
)

func CheckNeedsUpdate() bool {
	configDirPath := getConfigDir()
	VersionFilePath = filepath.Join(configDirPath, "version")

	if _, err := os.Stat(VersionFilePath); errors.Is(err, os.ErrNotExist) {
		createVersionFile()
		return false
	}
	currentVersion, err := os.ReadFile(VersionFilePath)
	if err != nil {
		utils.HandleError(err)
	}
	return string(currentVersion) != getLatestVersion()
}
