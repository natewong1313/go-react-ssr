package update

import (
	"os"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
)

var VersionFilePath string

func getConfigDir() string {
	configDir, _ := os.UserConfigDir()
	configDirPath := filepath.Join(configDir, "gossr-cli")
	os.MkdirAll(configDirPath, os.ModePerm)
	return configDirPath
}

func createVersionFile() {
	file, err := os.Create(VersionFilePath)
	if err != nil {
		utils.HandleError(err)
	}
	defer file.Close()
	file.WriteString(getLatestVersion())
}

func updateVersionFile() {
	file, err := os.OpenFile(VersionFilePath, os.O_WRONLY, 0644)
	if err != nil {
		utils.HandleError(err)
	}
	defer file.Close()
	file.WriteString(getLatestVersion())
}
