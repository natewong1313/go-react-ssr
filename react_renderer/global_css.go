package react_renderer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

var tempCssFolderPath string
var tempCssFilePath string

// Copies the contents of the global css file to a temporary css file
// and uses tailwind to compile it if tailwind is enabled
func BuildGlobalCSSFile() error {
	var err error
	// If the temporary folder has already been created, skip
	if tempCssFolderPath == "" {
		tempCssFolderPath, err = createTempCSSFolder()
		if err != nil {
			return err
		}
		tempCssFilePath, err = createTempCSSFile()
		if err != nil {
			return err
		}
	}
	// If tailwind is enabled, compile the global css file with tailwind and save the output to tempCssFilePath
	if config.C.TailwindConfigPath != "" {
		logger.L.Debug().Msg("Compiling tailwind css file")
		_, err = compileTailwindCssFile(tempCssFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a temporary folder in the local cache directory to store the temporary CSS file
func createTempCSSFolder() (string, error) {
	osCacheDir, _ := os.UserCacheDir()
	cacheFolderPath := filepath.Join(osCacheDir, "gossr-css")
	os.RemoveAll(cacheFolderPath)
	err := os.MkdirAll(cacheFolderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return cacheFolderPath, nil
}

// Create the temporary css file and copy the contents of the specified global css file to it
func createTempCSSFile() (string, error) {
	globalCssPath := utils.GetFullFilePath(config.C.GlobalCSSFilePath)
	tempFilePath := filepath.Join(tempCssFolderPath, "gossr-temporary.css")
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	// Copy contents of global css file to temp file
	globalCssFile, err := os.Open(globalCssPath)
	if err != nil {
		return "", err
	}
	defer globalCssFile.Close()
	_, err = io.Copy(tempFile, globalCssFile)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}
